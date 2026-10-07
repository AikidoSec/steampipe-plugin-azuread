package azuread

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"uuid"
)

type graphObject = map[string]any

type graphEndpoint string

const (
	graphDefault graphEndpoint = ""
	graphBeta    graphEndpoint = "beta"
	graphPortal  graphEndpoint = "portal"
)

type tokenSource func(context.Context) (string, error)

type GraphClient struct {
	http        *http.Client
	graphURL    string
	portalURL   string
	version     string
	token       tokenSource
	portalToken tokenSource
}

func (c *GraphClient) endpointURL(endpoint graphEndpoint, path string) string {
	if endpoint == graphPortal {
		return strings.TrimRight(c.portalURL, "/") + "/" + strings.TrimLeft(path, "/")
	}

	version := c.version
	if endpoint == graphBeta {
		version = "beta"
	}

	if version == "beta" && (path == "groupSettings" || strings.HasPrefix(path, "groupSettings/")) {
		path = "settings" + strings.TrimPrefix(path, "groupSettings")
	}

	return strings.TrimRight(c.graphURL, "/") + "/" + version + "/" + strings.TrimLeft(path, "/")
}

func (c *GraphClient) get(ctx context.Context, endpoint graphEndpoint, path string, query url.Values) (graphObject, error) {
	u, err := url.Parse(c.endpointURL(endpoint, path))
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	return c.requestUserPage(ctx, endpoint, path, u.String())
}

// Non-premium tenants cannot return signInActivity.
func (c *GraphClient) requestUserPage(ctx context.Context, endpoint graphEndpoint, path, rawURL string) (graphObject, error) {
	row, err := c.request(ctx, endpoint, rawURL)

	var requestErr *RequestError
	if endpoint == graphPortal || (path != "users" && !strings.HasPrefix(path, "users/")) || !errors.As(err, &requestErr) || requestErr.Code != "Authentication_RequestFromNonPremiumTenantOrB2CTenant" {
		return row, err
	}

	u, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		return nil, err
	}

	query := u.Query()
	// Filters on signInActivity still require a premium license.
	if strings.Contains(query.Get("$filter"), "signInActivity") {
		return nil, err
	}

	var fields []string
	removed := false

	for _, field := range strings.Split(query.Get("$select"), ",") {
		if strings.TrimSpace(field) == "signInActivity" {
			removed = true
			continue
		}

		fields = append(fields, field)
	}

	if !removed {
		return nil, err
	}

	query.Set("$select", strings.Join(fields, ","))
	u.RawQuery = query.Encode()

	return c.request(ctx, endpoint, u.String())
}

// Returning false from visit stops pagination.
func (c *GraphClient) list(ctx context.Context, endpoint graphEndpoint, path string, query url.Values, visit func(graphObject) bool) error {
	u, err := url.Parse(c.endpointURL(endpoint, path))
	if err != nil {
		return err
	}

	u.RawQuery = query.Encode()

	next := u.String()
	seen := make(map[string]bool)

	for next != "" {
		if seen[next] {
			return fmt.Errorf("graph returned a repeated pagination URL")
		}

		seen[next] = true

		var page graphObject
		if len(seen) == 1 {
			page, err = c.requestUserPage(ctx, endpoint, path, next)
		} else {
			// Graph continuation URLs are opaque.
			page, err = c.request(ctx, endpoint, next)
		}
		if err != nil {
			return err
		}

		values, ok := page["value"].([]any)
		if !ok {
			return fmt.Errorf("graph collection response is missing its value array")
		}

		for _, value := range values {
			if err := ctx.Err(); err != nil {
				return err
			}

			row, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf("graph collection contains a non-object value")
			}

			if !visit(row) {
				return nil
			}
		}

		next = ""
		if link, exists := page["@odata.nextLink"]; exists && link != nil {
			var ok bool
			next, ok = link.(string)
			if !ok {
				return fmt.Errorf("graph collection contains a non-string pagination URL")
			}
		}
	}

	return nil
}

func (c *GraphClient) request(ctx context.Context, endpoint graphEndpoint, rawURL string) (graphObject, error) {
	base := c.graphURL
	source := c.token
	if endpoint == graphPortal {
		base, source = c.portalURL, c.portalToken
	}

	if source == nil {
		return nil, fmt.Errorf("internal portal authentication requires internal_api_refresh_token or AZURE_INTERNAL_API_REFRESH_TOKEN")
	}

	// Pagination URLs must belong to the token's API origin.
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	allowed, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	if u.Scheme != allowed.Scheme || !strings.EqualFold(u.Host, allowed.Host) || u.User != nil {
		return nil, fmt.Errorf("refusing Graph pagination URL outside the configured API origin")
	}

	for attempt := 0; ; attempt++ {
		token, err := source(ctx)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}

		requestID := uuid.New().String()
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("client-request-id", requestID)
		req.Header.Set("x-ms-client-request-id", requestID)
		req.Header.Set("x-ms-correlation-id", uuid.New().String())
		req.Header.Set("Accept", "application/json")

		if endpoint != graphPortal {
			req.Header.Set("ConsistencyLevel", "eventual")
		}

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}

		body, readErr := io.ReadAll(resp.Body)
		if err := resp.Body.Close(); err != nil {
			return nil, err
		}

		if readErr != nil {
			return nil, readErr
		}

		if retryableStatus(resp.StatusCode) && attempt < 3 {
			if err := waitForRetry(ctx, retryDelay(resp.Header.Get("Retry-After"), attempt)); err != nil {
				return nil, err
			}

			continue
		}

		var result graphObject
		decodeErr := json.Unmarshal(body, &result)
		if decodeErr != nil {
			return nil, fmt.Errorf("decoding API response: %w", decodeErr)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			requestErr := &RequestError{
				StatusCode:      resp.StatusCode,
				Code:            http.StatusText(resp.StatusCode),
				Message:         "API request failed",
				RequestID:       resp.Header.Get("request-id"),
				ClientRequestID: requestID,
			}

			if e, ok := result["error"].(map[string]any); ok {
				if code, ok := e["code"].(string); ok {
					requestErr.Code = code
				}

				if message, ok := e["message"].(string); ok {
					requestErr.Message = message
				}
			}

			return nil, requestErr
		}

		if e, ok := result["error"].(map[string]any); ok {
			code, _ := e["code"].(string)
			message, _ := e["message"].(string)

			return nil, &RequestError{
				StatusCode:      resp.StatusCode,
				Code:            code,
				Message:         message,
				RequestID:       resp.Header.Get("request-id"),
				ClientRequestID: requestID,
			}
		}

		if result == nil {
			return nil, errors.New("API returned a null object")
		}

		return result, nil
	}
}

func retryableStatus(status int) bool {
	return status == http.StatusTooManyRequests ||
		status == http.StatusInternalServerError ||
		status == http.StatusBadGateway ||
		status == http.StatusServiceUnavailable ||
		status == http.StatusGatewayTimeout
}

func retryDelay(header string, attempt int) time.Duration {
	if seconds, err := strconv.Atoi(header); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second
	}

	if deadline, err := http.ParseTime(header); err == nil {
		if delay := time.Until(deadline); delay > 0 {
			return delay
		}

		return 0
	}

	return time.Second * time.Duration(1<<attempt)
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
