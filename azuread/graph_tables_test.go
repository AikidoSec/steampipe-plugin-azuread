package azuread

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/quals"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func stringQual(value string) *proto.QualValue {
	return &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: value}}
}
func queryData(table *plugin.Table, columns ...string) *plugin.QueryData {
	return &plugin.QueryData{Table: table, QueryContext: &plugin.QueryContext{Columns: columns}, EqualsQuals: plugin.KeyColumnEqualsQualMap{}, Quals: plugin.KeyColumnQualMap{}}
}
func addQual(d *plugin.QueryData, name, operator string, value *proto.QualValue) {
	d.Quals[name] = &plugin.KeyColumnQuals{Name: name, Quals: quals.QualSlice{&quals.Qual{Column: name, Operator: operator, Value: value}}}
	if operator == "=" {
		d.EqualsQuals[name] = value
	}
}

func TestOriginalTableSchemasRemainAvailable(t *testing.T) {
	data, err := os.ReadFile("testdata/original_schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var original map[string][][2]string
	if err := json.Unmarshal(data, &original); err != nil {
		t.Fatal(err)
	}
	tables := Plugin(context.Background()).TableMap
	if len(original) != 26 || len(tables) != 31 {
		t.Fatalf("unexpected table counts: original=%d current=%d", len(original), len(tables))
	}
	for name, columns := range original {
		table, ok := tables[name]
		if !ok {
			t.Fatalf("missing table %s", name)
		}
		for _, column := range columns {
			current := graphColumn(table, column[0])
			if current == nil || current.Type.String() != column[1] {
				t.Errorf("column contract changed: %s.%s (%v)", name, column[0], current)
			}
		}
	}
	for name, table := range tables {
		if _, ok := graphTables[name]; !ok {
			t.Errorf("missing endpoint for %s", name)
		}
		for _, column := range table.Columns {
			if column.Transform == nil {
				t.Errorf("missing JSON transform: %s.%s", name, column.Name)
				continue
			}
			if column.Hydrate != nil || column.Name == "filter" {
				continue
			}
			if _, err := column.Transform.Execute(context.Background(), &transform.TransformData{HydrateItem: graphObject{}, ColumnName: column.Name}); err != nil {
				t.Errorf("sparse %s.%s: %v", name, column.Name, err)
			}
		}
	}
}

func TestGraphProjectionFiltersAndNavigation(t *testing.T) {
	table := tableAzureAdUser(context.Background())
	d := queryData(table, "id", "title", "account_enabled", "sign_in_activity", "member_of", "manager", "tenant_id")
	addQual(d, "display_name", "=", stringQual("O'Brien"))
	addQual(d, "account_enabled", "<>", &proto.QualValue{Value: &proto.QualValue_BoolValue{BoolValue: true}})
	query := graphQuery(graphTables[table.Name], d, false)
	if query.Get("$select") != "accountEnabled,displayName,id,signInActivity,userPrincipalName" {
		t.Fatalf("incorrect projection: %s", query.Encode())
	}
	if query.Get("$top") != "500" || query.Get("$count") != "true" {
		t.Fatalf("incorrect paging: %s", query.Encode())
	}
	filter := query.Get("$filter")
	if !strings.Contains(filter, "displayName eq 'O''Brien'") || !strings.Contains(filter, "accountEnabled eq false") {
		t.Fatalf("bad filter %s", filter)
	}
	if query.Get("$expand") != "" {
		t.Fatal("membership must be independently paginated")
	}
	limit := int64(3)
	d.QueryContext.Limit = &limit
	if got := graphQuery(graphTables[table.Name], d, false).Get("$top"); got != "3" {
		t.Fatalf("LIMIT not pushed down: %s", got)
	}
	addQual(d, "filter", "=", stringQual("startswith(displayName,'A')"))
	if got := graphQuery(graphTables[table.Name], d, false).Get("$filter"); got != "startswith(displayName,'A')" {
		t.Fatalf("raw filter changed: %s", got)
	}
	if query := graphQuery(graphTables[table.Name], d, true); query.Get("$top") != "" || query.Get("$filter") != "" {
		t.Fatal("collection options on GET")
	}
	table = tableAzureAdDirectoryRoleEligibilityScheduleInstance(context.Background())
	d = queryData(table, "principal", "role_definition", "app_scope")
	if got := graphQuery(graphTables[table.Name], d, false).Get("$expand"); got != "principal,roleDefinition,appScope" {
		t.Fatalf("missing role expansions: %s", got)
	}
	table = tableAzureAdIdentityProvider(context.Background())
	d = queryData(table, "name")
	if got := graphQuery(graphTables[table.Name], d, false).Get("$top"); got != "" {
		t.Fatal("unsupported top on identity providers")
	}
}

func TestGraphParentPathsAndFlattenedSettings(t *testing.T) {
	d := queryData(tableAzureAdApplicationAppRoleAssignment(context.Background()))
	d.EqualsQuals["app_id"] = stringQual("app-guid")
	d.EqualsQuals["id"] = stringQual("assignment/id")
	spec := graphTables[d.Table.Name]
	path, err := graphTablePath(spec, d, true)
	if err != nil || path != "servicePrincipals(appId='app-guid')/appRoleAssignedTo/assignment%2Fid" {
		t.Fatalf("invalid alternate key: %s %v", path, err)
	}
	row := graphObject{"id": "assignment"}
	addParentColumn(row, spec, d)
	if row["appId"] != "app-guid" {
		t.Fatal("parent ID lost")
	}
	delete(d.EqualsQuals, "app_id")
	if _, err := graphTablePath(spec, d, false); err == nil {
		t.Fatal("missing required parent accepted")
	}
	rows := flattenSettings(graphObject{"id": "setting", "displayName": "Group.Unified", "templateId": "template", "values": []any{graphObject{"name": "EnableGroupCreation", "value": "false"}, graphObject{"name": "AllowGuests", "value": "true"}}})
	if len(rows) != 2 || rows[0]["value"] != "false" || rows[1]["id"] != "setting" || rows[1]["name"] != "AllowGuests" {
		t.Fatalf("bad flattened settings: %v", rows)
	}
}

func TestGraphTransformsPreserveLegacyShapes(t *testing.T) {
	tables := Plugin(context.Background()).TableMap
	cases := []struct {
		table, column string
		row           graphObject
		want          any
	}{
		{"azuread_user", "account_enabled", graphObject{"accountEnabled": false}, false},
		{"azuread_user", "account_enabled", graphObject{}, nil},
		{"azuread_user", "title", graphObject{"userPrincipalName": "someone@example.test"}, "someone@example.test"},
		{"azuread_application", "tags", graphObject{"tags": []any{"a", "b"}}, map[string]bool{"a": true, "b": true}},
		{"azuread_group", "tags", graphObject{"assignedLabels": []any{graphObject{"labelId": "a", "displayName": "Private"}}}, map[string]any{"a": "Private"}},
		{"azuread_conditional_access_policy", "authentication_strength", graphObject{"grantControls": graphObject{"authenticationStrength": graphObject{"allowedCombinations": []any{"fido2"}}}}, []any{"fido2"}},
		{"azuread_conditional_access_policy", "disable_resilience_defaults", graphObject{"sessionControls": graphObject{"disableResilienceDefaults": false}}, false},
		{"azuread_directory_role_eligibility_schedule_instance", "role_definition", graphObject{"roleDefinition": graphObject{"id": "a", "displayName": "Admin", "@odata.type": "#role"}}, graphObject{"id": "a", "display_name": "Admin", "@odata.type": "#role"}},
		{"azuread_conditional_access_named_location", "location_type", graphObject{"@odata.type": "#microsoft.graph.countryNamedLocation"}, "Country"},
		{"azuread_conditional_access_named_location", "location_info", graphObject{"@odata.type": "#microsoft.graph.countryNamedLocation", "countriesAndRegions": []any{"RO"}, "includeUnknownCountriesAndRegions": false, "countryLookupMethod": "clientIpAddress"}, graphObject{"Countries_and_Regions": []any{"RO"}, "Get_Unknown_Countries_and_Regions": false, "Lookup_Method": "clientIpAddress"}},
		{"azuread_user", "sign_in_activity", graphObject{"signInActivity": graphObject{"lastSignInRequestId": "a"}}, graphObject{"LastSignInDateTime": nil, "LastSignInRequestId": "a", "LastNonInteractiveSignInDateTime": nil, "LastNonInteractiveSignInRequestId": nil}},
	}
	for _, tt := range cases {
		t.Run(tt.table+"/"+tt.column, func(t *testing.T) {
			col := graphColumn(tables[tt.table], tt.column)
			got, err := col.Transform.Execute(context.Background(), &transform.TransformData{HydrateItem: tt.row, ColumnName: tt.column})
			if err != nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v err=%v", got, tt.want, err)
			}
		})
	}
}

func TestGetHydratorRoutesAllTables(t *testing.T) {
	ctx := context.Background()
	cache, err := connection.NewConnectionCache("http-test", 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	for name, table := range Plugin(ctx).TableMap {
		t.Run(name, func(t *testing.T) {
			d := queryData(table, "id")
			d.ConnectionManager = connection.NewManager(cache)
			d.EqualsQuals["id"] = stringQual("object")
			d.EqualsQuals["name"] = stringQual("option")
			spec := graphTables[name]
			if spec.parentColumn != "" {
				d.EqualsQuals[spec.parentColumn] = stringQual("parent")
			}
			client := testGraphClient(func(r *http.Request) (*http.Response, error) {
				expectedPath, _ := graphTablePath(spec, d, true)
				expected := (&GraphClient{graphURL: "https://graph.example", portalURL: "https://portal.example/api"}).endpointURL(spec.endpoint, expectedPath)
				if q := graphQuery(spec, d, true).Encode(); q != "" {
					expected += "?" + q
				}
				if r.URL.String() != expected {
					t.Fatalf("unexpected request: %s want %s", r.URL, expected)
				}
				if name == "azuread_directory_setting" {
					return response(200, `{"id":"object","values":[{"name":"option","value":"setting-value"}]}`), nil
				}
				return response(200, `{"id":"object","accountEnabled":false}`), nil
			})
			d.ConnectionManager.Cache.Set("graphHTTPClient", client)
			var row any
			var err error
			if name == "azuread_directory_setting" {
				row, err = getAzureAdDirectorySetting(ctx, d, nil)
			} else {
				row, err = getGraphTable(ctx, d, nil)
			}
			if err != nil {
				t.Fatal(err)
			}
			if row.(graphObject)["id"] != "object" {
				t.Fatalf("unexpected row: %v", row)
			}
		})
	}
}

func TestAssetEndpointDefaults(t *testing.T) {
	client := &GraphClient{graphURL: "https://graph.microsoft.com"}
	expected := map[string]string{
		"azuread_user":                              "https://graph.microsoft.com/v1.0/users",
		"azuread_group":                             "https://graph.microsoft.com/v1.0/groups",
		"azuread_application":                       "https://graph.microsoft.com/v1.0/applications",
		"azuread_service_principal":                 "https://graph.microsoft.com/v1.0/servicePrincipals",
		"azuread_device":                            "https://graph.microsoft.com/v1.0/devices",
		"azuread_conditional_access_policy":         "https://graph.microsoft.com/beta/identity/conditionalAccess/policies",
		"azuread_conditional_access_named_location": "https://graph.microsoft.com/beta/identity/conditionalAccess/namedLocations",
		"azuread_directory_setting":                 "https://graph.microsoft.com/beta/settings",
		"azuread_device_registration_policy":        "https://graph.microsoft.com/beta/policies/deviceRegistrationPolicy",
	}
	for table, want := range expected {
		spec := graphTables[table]
		if got := client.endpointURL(spec.endpoint, spec.path); got != want {
			t.Errorf("%s: got %s, want %s", table, got, want)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single lowercase word",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "single uppercase word",
			input:    "HELLO",
			expected: "hello",
		},
		{
			name:     "space separated",
			input:    "hello world",
			expected: "helloWorld",
		},
		{
			name:     "underscore separated",
			input:    "hello_world",
			expected: "helloWorld",
		},
		{
			name:     "hyphen separated",
			input:    "hello-world",
			expected: "helloWorld",
		},
		{
			name:     "dot separated",
			input:    "hello.world",
			expected: "helloWorld",
		},
		{
			name:     "mixed separators",
			input:    "hello_world-test.example",
			expected: "helloWorldTestExample",
		},
		{
			name:     "uppercase words",
			input:    "HELLO_WORLD",
			expected: "helloWorld",
		},
		{
			name:     "mixed casing",
			input:    "HeLLo_WoRLD",
			expected: "helloWorld",
		},
		{
			name:     "multiple consecutive separators",
			input:    "hello---___...world",
			expected: "helloWorld",
		},
		{
			name:     "leading separators",
			input:    "__hello_world",
			expected: "helloWorld",
		},
		{
			name:     "trailing separators",
			input:    "hello_world___",
			expected: "helloWorld",
		},
		{
			name:     "leading and trailing separators",
			input:    "___hello_world---",
			expected: "helloWorld",
		},
		{
			name:     "only separators",
			input:    "_- ._",
			expected: "",
		},
		{
			name:     "unicode",
			input:    "über_café",
			expected: "überCafé",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toCamelCase(tt.input)

			if got != tt.expected {
				t.Errorf(
					"toCamelCase(%q) = %q; want %q",
					tt.input,
					got,
					tt.expected,
				)
			}
		})
	}
}
