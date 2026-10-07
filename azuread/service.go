package azuread

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func GetGraphClient(_ context.Context, d *plugin.QueryData) (*GraphClient, error) {
	const cacheKey = "graphHTTPClient"
	if d.ConnectionManager != nil {
		if cached, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
			return cached.(*GraphClient), nil
		}
	}

	client, err := newGraphClient(GetConfig(d.Connection))
	if err != nil {
		return nil, err
	}

	if d.ConnectionManager != nil {
		d.ConnectionManager.Cache.Set(cacheKey, client)
	}

	return client, nil
}

func configString(value *string, env string) string {
	if value != nil {
		return *value
	}

	return os.Getenv(env)
}

func newGraphClient(config AzureADConfig) (*GraphClient, error) {
	tenant := configString(config.TenantID, "AZURE_TENANT_ID")
	clientID := configString(config.ClientID, "AZURE_CLIENT_ID")
	secret := configString(config.ClientSecret, "AZURE_CLIENT_SECRET")
	environment := configString(config.Environment, "AZURE_ENVIRONMENT")
	version := configString(config.GraphAPIVersion, "AZURE_GRAPH_API_VERSION")

	if version == "" {
		version = "v1.0"
	}

	if version != "v1.0" && version != "beta" {
		return nil, fmt.Errorf("graph_api_version must be v1.0 or beta")
	}

	graphURL, loginURL, cloudConfig, err := graphCloud(environment)
	if err != nil {
		return nil, err
	}

	options := azcore.ClientOptions{Cloud: cloudConfig}
	var credential azcore.TokenCredential

	switch {
	case tenant != "" && clientID != "" && secret != "":
		credential, err = azidentity.NewClientSecretCredential(tenant, clientID, secret, &azidentity.ClientSecretCredentialOptions{
			ClientOptions: options,
		})

	case config.EnableMsi != nil && *config.EnableMsi:
		msiOptions := &azidentity.ManagedIdentityCredentialOptions{ClientOptions: options}
		if clientID != "" {
			msiOptions.ID = azidentity.ClientID(clientID)
		}

		credential, err = azidentity.NewManagedIdentityCredential(msiOptions)

	case config.ClientID != nil || config.ClientSecret != nil || clientID != "" || secret != "":
		return nil, fmt.Errorf("client-secret authentication requires non-empty tenant_id, client_id, and client_secret")

	default:
		credential, err = azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{TenantID: tenant})
	}

	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		// API and token endpoints must not redirect credentials to another service.
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	client := &GraphClient{
		http:      httpClient,
		graphURL:  graphURL,
		version:   version,
		portalURL: "https://main.iam.ad.ext.azure.com/api",
		token: func(ctx context.Context) (string, error) {
			token, err := credential.GetToken(ctx, policy.TokenRequestOptions{
				Scopes: []string{graphURL + "/.default"},
			})

			return token.Token, err
		},
	}

	refresh := configString(config.InternalAPIRefreshToken, "AZURE_INTERNAL_API_REFRESH_TOKEN")
	if refresh != "" {
		if graphURL != "https://graph.microsoft.com" {
			return nil, fmt.Errorf("internal portal collection is only available in AZUREPUBLICCLOUD")
		}

		if tenant == "" {
			return nil, fmt.Errorf("tenant_id is required with internal_api_refresh_token")
		}

		client.portalToken = newPortalTokenSource(httpClient, loginURL+"/"+url.PathEscape(tenant)+"/oauth2/token", refresh)
	}

	return client, nil
}

func graphCloud(environment string) (string, string, cloud.Configuration, error) {
	switch environment {
	case "", "AZUREPUBLICCLOUD":
		return "https://graph.microsoft.com", "https://login.microsoftonline.com", cloud.AzurePublic, nil
	case "AZURECHINACLOUD":
		return "https://microsoftgraph.chinacloudapi.cn", "https://login.chinacloudapi.cn", cloud.AzureChina, nil
	case "AZUREUSGOVERNMENTCLOUD":
		return "https://graph.microsoft.us", "https://login.microsoftonline.us", cloud.AzureGovernment, nil
	default:
		return "", "", cloud.Configuration{}, fmt.Errorf("unsupported Azure environment %q", environment)
	}
}

// The portal requires its own delegated token, not a Microsoft Graph token.
func newPortalTokenSource(client *http.Client, tokenURL, refreshToken string) tokenSource {
	var mu sync.Mutex
	var accessToken string
	var expires time.Time

	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()

		if accessToken != "" && time.Now().Add(time.Minute).Before(expires) {
			return accessToken, nil
		}

		form := url.Values{
			"client_id":     {"1950a258-227b-4e31-a9cf-717495945fc2"},
			"grant_type":    {"refresh_token"},
			"refresh_token": {refreshToken},
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}()

		var result struct {
			AccessToken  string          `json:"access_token"`
			RefreshToken string          `json:"refresh_token"`
			ExpiresIn    json.RawMessage `json:"expires_in"`
			Error        string          `json:"error"`
		}

		if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
			return "", fmt.Errorf("decoding internal portal token response: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 || result.Error != "" || result.AccessToken == "" {
			return "", fmt.Errorf("internal portal token acquisition failed (HTTP %d, code %s)", resp.StatusCode, result.Error)
		}

		seconds, err := strconv.Atoi(strings.Trim(string(result.ExpiresIn), `"`))
		if err != nil || seconds <= 0 {
			seconds = 300
		}

		accessToken, expires = result.AccessToken, time.Now().Add(time.Duration(seconds)*time.Second)
		if result.RefreshToken != "" {
			refreshToken = result.RefreshToken
		}

		return accessToken, nil
	}
}

// https://github.com/Azure/go-autorest/blob/3fb5326fea196cd5af02cf105ca246a0fba59021/autorest/azure/cli/token.go#L126
// NewAuthorizerFromCLIWithResource creates an Authorizer configured from Azure CLI 2.0 for local development scenarios.
func getTenantFromCLI() (string, error) {
	// This is the path that a developer can set to tell this class what the install path for Azure CLI is.
	const azureCLIPath = "AzureCLIPath"

	// The default install paths are used to find Azure CLI. This is for security, so that any path in the calling program's Path environment is not used to execute Azure CLI.
	azureCLIDefaultPathWindows := fmt.Sprintf("%s\\Microsoft SDKs\\Azure\\CLI2\\wbin; %s\\Microsoft SDKs\\Azure\\CLI2\\wbin", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

	// Default path for non-Windows.
	const azureCLIDefaultPath = "/bin:/sbin:/usr/bin:/usr/local/bin"

	// Execute Azure CLI to get token
	var cliCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cliCmd = exec.Command(fmt.Sprintf("%s\\system32\\cmd.exe", os.Getenv("windir")))
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s;%s", os.Getenv(azureCLIPath), azureCLIDefaultPathWindows))
		cliCmd.Args = append(cliCmd.Args, "/c", "az")
	} else {
		cliCmd = exec.Command("az")
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s:%s", os.Getenv(azureCLIPath), azureCLIDefaultPath))
	}
	cliCmd.Args = append(cliCmd.Args, "account", "get-access-token", "--resource-type=ms-graph", "-o", "json")

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		return "", fmt.Errorf("Invoking Azure CLI failed with the following error: %v", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"accessToken"`
		ExpiresOn   string `json:"expiresOn"`
		Tenant      string `json:"tenant"`
		TokenType   string `json:"tokenType"`
	}
	err = json.Unmarshal(output, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.Tenant, nil
}
