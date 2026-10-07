package azuread

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error { b.closed = true; return nil }
func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body))}
}
func testGraphClient(f roundTripFunc) *GraphClient {
	return &GraphClient{http: &http.Client{Transport: f}, graphURL: "https://graph.example", portalURL: "https://portal.example/api", version: "v1.0", token: func(context.Context) (string, error) { return "graph-token", nil }, portalToken: func(context.Context) (string, error) { return "portal-token", nil }}
}

func TestGraphPaginationPreservesOpaqueURLAndHeaders(t *testing.T) {
	next := "https://graph.example/v1.0/users?$skip=12&$select=id%2CdisplayName&$filter=accountEnabled%20eq%20true"
	calls := 0
	bodies := []*trackedBody{}
	client := testGraphClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Authorization") != "Bearer graph-token" || r.Header.Get("ConsistencyLevel") != "eventual" {
			t.Fatal("authentication/query headers missing")
		}
		var payload string
		if calls == 1 {
			if r.URL.Query().Get("$select") != "id,displayName" {
				t.Fatal("initial projection lost")
			}
			encoded, _ := json.Marshal(graphObject{"value": []graphObject{{"id": "one", "displayName": "A"}}, "@odata.nextLink": next})
			payload = string(encoded)
		} else {
			if r.URL.String() != next {
				t.Fatalf("nextLink changed: %s", r.URL)
			}
			payload = `{"value":[{"id":"two","displayName":"B"}]}`
		}
		body := &trackedBody{Reader: strings.NewReader(payload)}
		bodies = append(bodies, body)
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: body}, nil
	})
	var names []string
	err := client.list(context.Background(), graphDefault, "users", url.Values{"$select": {"id,displayName"}}, func(row graphObject) bool { names = append(names, row["displayName"].(string)); return true })
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || !reflect.DeepEqual(names, []string{"A", "B"}) {
		t.Fatalf("calls=%d rows=%v", calls, names)
	}
	for _, body := range bodies {
		if !body.closed {
			t.Fatal("response body leaked")
		}
	}
}

func TestGraphPaginationStopsAndRejectsUnsafeNextLinks(t *testing.T) {
	for _, next := range []string{"https://attacker.example/users", "http://graph.example/users", "https://user@graph.example/users"} {
		t.Run(next, func(t *testing.T) {
			calls := 0
			client := testGraphClient(func(r *http.Request) (*http.Response, error) {
				calls++
				return response(200, fmt.Sprintf(`{"value":[{"id":"one"}],"@odata.nextLink":%q}`, next)), nil
			})
			if err := client.list(context.Background(), graphDefault, "users", nil, func(graphObject) bool { return true }); err == nil {
				t.Fatal("unsafe nextLink accepted")
			}
			if calls != 1 {
				t.Fatal("token sent to untrusted origin")
			}
			calls = 0
			if err := client.list(context.Background(), graphDefault, "users", nil, func(graphObject) bool { return false }); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatal("paging continued after limit")
			}
		})
	}
}

func TestGraphRetriesErrorsAndCancellation(t *testing.T) {
	for _, status := range []int{429, 500, 502, 503, 504} {
		t.Run(strconvStatus(status), func(t *testing.T) {
			calls := 0
			var bodies []*trackedBody
			client := testGraphClient(func(r *http.Request) (*http.Response, error) {
				calls++
				code := status
				payload := `{"error":{"code":"Busy","message":"retry"}}`
				if calls == 2 {
					code = 200
					payload = `{"id":"ok"}`
				}
				body := &trackedBody{Reader: strings.NewReader(payload)}
				bodies = append(bodies, body)
				return &http.Response{StatusCode: code, Header: http.Header{"Retry-After": []string{"0"}}, Body: body}, nil
			})
			row, err := client.get(context.Background(), graphDefault, "users/a", nil)
			if err != nil || row["id"] != "ok" || calls != 2 {
				t.Fatalf("row=%v err=%v calls=%d", row, err, calls)
			}
			for _, body := range bodies {
				if !body.closed {
					t.Fatal("retry body leaked")
				}
			}
		})
	}
	calls := 0
	client := testGraphClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return response(403, `{"error":{"code":"Authorization_RequestDenied","message":"denied"}}`), nil
	})
	_, err := client.get(context.Background(), graphDefault, "users", nil)
	var apiErr *RequestError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 403 || apiErr.Code != "Authorization_RequestDenied" || calls != 1 {
		t.Fatalf("unexpected error/fallback: %v calls=%d", err, calls)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := waitForRetry(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Fatalf("retry ignored cancellation: %v", err)
	}
	if retryDelay("garbage", 2) != 4*time.Second {
		t.Fatal("missing retry header should use bounded backoff")
	}
	if retryDelay(time.Now().Add(time.Hour).UTC().Format(http.TimeFormat), 0) < 59*time.Minute {
		t.Fatal("HTTP-date retry header ignored")
	}
}
func strconvStatus(status int) string { return fmt.Sprint(status) }

func TestGraphRetryExhaustionAndMalformedResponses(t *testing.T) {
	calls := 0
	client := testGraphClient(func(r *http.Request) (*http.Response, error) {
		calls++
		resp := response(429, `{"error":{"code":"TooManyRequests","message":"slow down"}}`)
		resp.Header.Set("Retry-After", "0")
		return resp, nil
	})
	_, err := client.get(context.Background(), graphDefault, "users", nil)
	if err == nil || calls != 4 {
		t.Fatalf("retries not bounded: calls=%d err=%v", calls, err)
	}
	for _, payload := range []string{`not json`, `null`, `{"error":{"code":"Bad","message":"failed"}}`, `{"value":null}`} {
		client = testGraphClient(func(r *http.Request) (*http.Response, error) { return response(200, payload), nil })
		if err := client.list(context.Background(), graphDefault, "users", nil, func(graphObject) bool { return true }); err == nil {
			t.Fatalf("accepted %s", payload)
		}
	}
}

func TestGraphEndpointRoutingAndSeparatePortalToken(t *testing.T) {
	client := testGraphClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://portal.example/api/Directories/Properties" || r.Header.Get("Authorization") != "Bearer portal-token" {
			t.Fatalf("wrong portal request: %v", r)
		}
		return response(200, `{"objectId":"tenant"}`), nil
	})
	for _, tt := range []struct {
		version    string
		endpoint   graphEndpoint
		path, want string
	}{
		{"v1.0", graphDefault, "users", "https://graph.example/v1.0/users"},
		{"v1.0", graphBeta, "policies/deviceRegistrationPolicy", "https://graph.example/beta/policies/deviceRegistrationPolicy"},
		{"beta", graphDefault, "identity/conditionalAccess/policies", "https://graph.example/beta/identity/conditionalAccess/policies"},
		{"v1.0", graphDefault, "groupSettings/a", "https://graph.example/v1.0/groupSettings/a"},
		{"beta", graphDefault, "groupSettings/a", "https://graph.example/beta/settings/a"},
	} {
		client.version = tt.version
		if got := client.endpointURL(tt.endpoint, tt.path); got != tt.want {
			t.Fatalf("got %s want %s", got, tt.want)
		}
	}
	if _, err := client.get(context.Background(), graphPortal, "Directories/Properties", nil); err != nil {
		t.Fatal(err)
	}
	client.portalToken = nil
	if _, err := client.get(context.Background(), graphPortal, "Directories/Properties", nil); err == nil {
		t.Fatal("missing portal credentials should fail clearly")
	}
}

func TestPortalTokenCachingRotationAndErrors(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		expected := "initial"
		if calls > 1 {
			expected = "rotated"
		}
		if r.Method != "POST" || r.Form.Get("refresh_token") != expected || r.Form.Get("grant_type") != "refresh_token" || r.Form.Get("client_id") != "1950a258-227b-4e31-a9cf-717495945fc2" {
			t.Fatalf("incorrect portal grant: %v", r.Form)
		}
		if calls == 1 {
			return response(200, `{"access_token":"short","refresh_token":"rotated","expires_in":1}`), nil
		}
		return response(200, `{"access_token":"long","expires_in":"3600"}`), nil
	})}
	source := newPortalTokenSource(client, "https://login.example/tenant/oauth2/token", "initial")
	for _, want := range []string{"short", "long", "long"} {
		token, err := source(context.Background())
		if err != nil || token != want {
			t.Fatalf("token=%s err=%v", token, err)
		}
	}
	if calls != 2 {
		t.Fatal("valid portal token was not cached")
	}
	client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return response(400, `{"error":"invalid_grant","error_description":"secret must not leak"}`), nil
	})
	source = newPortalTokenSource(client, "https://login.example/tenant/oauth2/token", "initial")
	_, err := source(context.Background())
	if err == nil || strings.Contains(err.Error(), "secret must not leak") {
		t.Fatalf("unsafe token error: %v", err)
	}
}
