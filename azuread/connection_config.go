package azuread

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type AzureADConfig struct {
	TenantID                *string `hcl:"tenant_id"`
	ClientID                *string `hcl:"client_id"`
	ClientSecret            *string `hcl:"client_secret"`
	EnableMsi               *bool   `hcl:"enable_msi"`
	MsiEndpoint             *string `hcl:"msi_endpoint"`
	Environment             *string `hcl:"environment"`
	InternalAPIRefreshToken *string `hcl:"internal_api_refresh_token"`
}

func ConfigInstance() any {
	return &AzureADConfig{}
}

// GetConfig retrieves and cast connection config from query data
func GetConfig(connection *plugin.Connection) AzureADConfig {
	if connection == nil || connection.Config == nil {
		return AzureADConfig{}
	}
	config, _ := connection.Config.(AzureADConfig)
	return config
}
