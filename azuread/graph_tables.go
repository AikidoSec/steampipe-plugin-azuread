package azuread

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type graphFilter struct {
	// Graph property name; empty disables filter pushdown for this column
	field string
	// Emit the value as an OData literal rather than a quoted string
	unquoted bool
}

type graphTableSpec struct {
	// API-relative resource path
	path string
	// API target: Graph v1.0 by default, Graph beta, or the internal portal
	endpoint graphEndpoint
	// The endpoint returns one object rather than a paginated collection
	singleton bool
	// Required query column whose value replaces {parent} in the path
	parentColumn string
	// Default $top value, zero omits $top
	pageSize int
	// Build $select from the requested table columns
	project bool
	// Column-to-$select overrides, an empty field list omits the column
	selectFields map[string][]string
	// Per-column filter overrides for Graph property names and value quoting
	filters map[string]graphFilter
	// Include $count for every list request.
	count bool
	// Maximum page size when a property appears in $select or $filter
	pageSizeLimits map[string]int
	// Fetch partial Graph data when an internal portal token is not configured
	fallback func(context.Context, *GraphClient) (graphObject, error)
}

var graphTables = map[string]graphTableSpec{
	"azuread_user": {
		selectFields: map[string][]string{
			"title": {"displayName", "userPrincipalName"},
		},
		count: true,
		pageSizeLimits: map[string]int{
			"signInActivity": 500,
		},
		path:     "users",
		pageSize: 999,
		project:  true,
	},
	"azuread_group": {
		selectFields: map[string][]string{
			"title": {"displayName"},
			"tags":  {"assignedLabels"},
		},
		count:    true,
		path:     "groups",
		pageSize: 999,
		project:  true,
	},
	"azuread_application": {
		selectFields: map[string][]string{
			"title":                            {"displayName"},
			"is_authorization_service_enabled": nil,
		},
		count:    true,
		path:     "applications",
		pageSize: 999,
		project:  true,
	},
	"azuread_service_principal": {
		selectFields: map[string][]string{
			"title": {"displayName"},
		},
		count:    true,
		path:     "servicePrincipals",
		pageSize: 100,
		project:  true,
	},
	"azuread_device": {
		selectFields: map[string][]string{
			"title": {"displayName", "deviceId"},
		},
		count:    true,
		path:     "devices",
		pageSize: 999,
		project:  true,
	},
	"azuread_domain": {
		// Microsoft documents known issues with $top on this endpoint.
		path: "domains",
	},
	"azuread_directory_role": {
		path: "directoryRoles",
	},
	"azuread_directory_role_template": {
		path: "directoryRoleTemplates",
	},
	"azuread_directory_role_definition": {
		path: "roleManagement/directory/roleDefinitions",
	},
	"azuread_directory_role_assignment": {
		path: "roleManagement/directory/roleAssignments",
	},
	"azuread_directory_role_eligibility_schedule_instance": {
		path: "roleManagement/directory/roleEligibilityScheduleInstances",
	},
	"azuread_directory_setting": {
		path:     "settings",
		endpoint: graphBeta,
	},
	"azuread_identity_provider": {
		filters: map[string]graphFilter{
			"name": {
				field: "displayName",
			},
		},
		path: "identity/identityProviders",
	},
	"azuread_conditional_access_policy": {
		path:     "identity/conditionalAccess/policies",
		endpoint: graphBeta,
		pageSize: 1000,
	},
	"azuread_conditional_access_named_location": {
		filters: map[string]graphFilter{
			"location_type": {},
		},
		path:     "identity/conditionalAccess/namedLocations",
		endpoint: graphBeta,
		pageSize: 1000,
	},
	"azuread_sign_in_report": {
		path:     "auditLogs/signIns",
		pageSize: 999,
	},
	"azuread_directory_audit_report": {
		path:     "auditLogs/directoryAudits",
		pageSize: 1000,
	},
	"azuread_user_registration_details_report": {
		path: "reports/authenticationMethods/userRegistrationDetails",
	},
	"azuread_admin_consent_request_policy": {
		path:      "policies/adminConsentRequestPolicy",
		singleton: true,
	},
	"azuread_security_defaults_policy": {
		path:      "policies/identitySecurityDefaultsEnforcementPolicy",
		singleton: true,
	},
	"azuread_authorization_policy": {
		path:      "policies/authorizationPolicy",
		singleton: true,
	},
	"azuread_user_app_role_assignment": {
		count: true,
		filters: map[string]graphFilter{
			"resource_id": {
				field:    "resourceId",
				unquoted: true,
			},
		},
		path:         "users/{parent}/appRoleAssignments",
		parentColumn: "user_id",
		pageSize:     999,
	},
	"azuread_group_app_role_assignment": {
		filters: map[string]graphFilter{
			"resource_id": {
				field:    "resourceId",
				unquoted: true,
			},
		},
		path:         "groups/{parent}/appRoleAssignments",
		parentColumn: "group_id",
		pageSize:     999,
	},
	"azuread_service_principal_app_role_assignment": {
		filters: map[string]graphFilter{
			"resource_id": {
				field:    "resourceId",
				unquoted: true,
			},
		},
		path:         "servicePrincipals/{parent}/appRoleAssignments",
		parentColumn: "service_principal_id",
		pageSize:     999,
	},
	"azuread_service_principal_app_role_assigned_to": {
		filters: map[string]graphFilter{
			"resource_id": {
				field:    "resourceId",
				unquoted: true,
			},
		},
		path:         "servicePrincipals/{parent}/appRoleAssignedTo",
		parentColumn: "service_principal_id",
	},
	"azuread_application_app_role_assigned_to": {
		filters: map[string]graphFilter{
			"resource_id": {
				field:    "resourceId",
				unquoted: true,
			},
		},
		path:         "servicePrincipals(appId='{parent}')/appRoleAssignedTo",
		parentColumn: "app_id",
	},
	"azuread_device_registration_policy": {
		path:      "policies/deviceRegistrationPolicy",
		endpoint:  graphBeta,
		singleton: true,
	},
	"azuread_password_reset_policy": {
		fallback:  getPasswordResetPolicyFromGraph,
		path:      "PasswordReset/PasswordResetPolicies",
		endpoint:  graphPortal,
		singleton: true,
	},
	"azuread_password_policy": {
		fallback:  getPasswordPolicyFromGraph,
		path:      "AuthenticationMethods/PasswordPolicy",
		endpoint:  graphPortal,
		singleton: true,
	},
	"azuread_directory_properties": {
		fallback:  getDirectoryPropertiesFromGraph,
		path:      "Directories/Properties",
		endpoint:  graphPortal,
		singleton: true,
	},
	"azuread_self_service_group_management": {
		fallback:  getSelfServiceGroupManagementFromGraph,
		path:      "Directories/SsgmProperties",
		endpoint:  graphPortal,
		singleton: true,
	},
}

func graphTablePath(spec graphTableSpec, d *plugin.QueryData, get bool) (string, error) {
	path := spec.path

	if spec.parentColumn != "" {
		parent := d.EqualsQuals[spec.parentColumn].GetStringValue()
		if parent == "" {
			return "", fmt.Errorf("%s is required", spec.parentColumn)
		}

		if spec.parentColumn == "app_id" {
			parent = strings.ReplaceAll(parent, "'", "''")
		}

		path = strings.ReplaceAll(path, "{parent}", url.PathEscape(parent))
	}

	if get && !spec.singleton {
		id := d.EqualsQuals["id"].GetStringValue()
		if id == "" {
			return "", fmt.Errorf("id is required")
		}

		path += "/" + url.PathEscape(id)
	}

	return path, nil
}

func listGraphTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	spec, ok := graphTables[d.Table.Name]
	if !ok {
		return nil, fmt.Errorf("no API endpoint for %s", d.Table.Name)
	}

	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	if spec.endpoint == graphPortal && client.portalToken == nil {
		if spec.fallback == nil {
			return nil, fmt.Errorf("no Graph policy mapping for %s", d.Table.Name)
		}

		row, err := spec.fallback(ctx, client)
		if err != nil {
			return nil, err
		}

		d.StreamListItem(ctx, row)

		return nil, nil
	}

	path, err := graphTablePath(spec, d, false)
	if err != nil {
		return nil, err
	}

	query := graphQuery(spec, d, false)
	if spec.singleton {
		row, err := client.get(ctx, spec.endpoint, path, query)
		if err != nil {
			return nil, err
		}

		if spec.endpoint == graphPortal {
			payload := row
			row = graphObject{}

			for key, value := range payload {
				row[key] = value
			}

			row["dataSource"] = "internal_portal"
			row["sourcePayload"] = payload
		}

		d.StreamListItem(ctx, row)

		return nil, nil
	}

	err = client.list(ctx, spec.endpoint, path, query, func(row graphObject) bool {
		addParentColumn(row, spec, d)
		d.StreamListItem(ctx, row)

		return d.RowsRemaining(ctx) != 0
	})

	return nil, err
}

func getGraphTable(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	spec, ok := graphTables[d.Table.Name]
	if !ok {
		return nil, fmt.Errorf("no API endpoint for %s", d.Table.Name)
	}

	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	path, err := graphTablePath(spec, d, true)
	if err != nil {
		return nil, err
	}

	row, err := client.get(ctx, spec.endpoint, path, graphQuery(spec, d, true))
	if err != nil {
		return nil, err
	}

	addParentColumn(row, spec, d)

	return row, nil
}

func addParentColumn(row graphObject, spec graphTableSpec, d *plugin.QueryData) {
	if spec.parentColumn != "" {
		row[toCamelCase(spec.parentColumn)] = d.EqualsQuals[spec.parentColumn].GetStringValue()
	}
}

func graphQuery(spec graphTableSpec, d *plugin.QueryData, get bool) url.Values {
	query := url.Values{}
	columns := d.QueryContext.Columns

	if spec.project {
		selected := map[string]bool{"id": true}

		for _, name := range columns {
			column := graphColumn(d.Table, name)
			if column == nil || column.Hydrate != nil || name == "filter" || name == "tenant_id" {
				continue
			}

			if fields, ok := spec.selectFields[name]; ok {
				for _, field := range fields {
					selected[field] = true
				}

				continue
			}

			path := toCamelCase(name)
			if column.Transform != nil && len(column.Transform.Transforms) > 0 {
				if field, ok := column.Transform.Transforms[0].Param.(string); ok {
					path = field
				}
			}

			selected[strings.Split(path, ".")[0]] = true
		}

		keys := make([]string, 0, len(selected))
		for key := range selected {
			keys = append(keys, key)
		}

		sort.Strings(keys)
		query.Set("$select", strings.Join(keys, ","))
	}

	if get || spec.singleton {
		return query
	}

	if raw := d.EqualsQuals["filter"].GetStringValue(); raw != "" {
		query.Set("$filter", raw)
	} else {
		filters := graphFilters(spec, d)
		if len(filters) > 0 {
			query.Set("$filter", strings.Join(filters, " and "))
		}
	}

	if spec.count {
		query.Set("$count", "true")
	}

	pageSize := spec.pageSize
	for field, maximum := range spec.pageSizeLimits {
		if strings.Contains(query.Get("$select"), field) || strings.Contains(query.Get("$filter"), field) {
			if pageSize > maximum {
				pageSize = maximum
			}
		}
	}

	if pageSize > 0 {
		if limit := d.QueryContext.Limit; limit != nil && *limit > 0 && *limit < int64(pageSize) {
			pageSize = int(*limit)
		}

		query.Set("$top", strconv.Itoa(pageSize))
	}

	return query
}

func graphFilters(spec graphTableSpec, d *plugin.QueryData) []string {
	if d.Table.List == nil {
		return nil
	}

	var filters []string
	for _, key := range d.Table.List.KeyColumns {
		name := key.Name
		if name == spec.parentColumn || name == "filter" {
			continue
		}

		col := graphColumn(d.Table, name)
		if col == nil {
			continue
		}

		filter, ok := spec.filters[name]
		if !ok {
			filter.field = toCamelCase(name)
		}
		if filter.field == "" {
			continue
		}

		quals := d.Quals[name]
		if quals == nil {
			continue
		}

		for _, q := range quals.Quals {
			op := map[string]string{
				"=":  "eq",
				"<>": "ne",
				">":  "gt",
				">=": "ge",
				"<":  "lt",
				"<=": "le",
			}[q.Operator]
			if op == "" {
				continue
			}

			var value string
			switch col.Type {
			case proto.ColumnType_BOOL:
				boolValue := q.Value.GetBoolValue()
				if q.Operator == "<>" {
					op = "eq"
					boolValue = !boolValue
				}

				value = strconv.FormatBool(boolValue)
			case proto.ColumnType_TIMESTAMP:
				value = q.Value.GetTimestampValue().AsTime().Format(time.RFC3339Nano)
			default:
				value = q.Value.GetStringValue()
				if !filter.unquoted {
					value = "'" + strings.ReplaceAll(value, "'", "''") + "'"
				}
			}

			filters = append(filters, filter.field+" "+op+" "+value)
		}
	}

	return filters
}

func graphColumn(table *plugin.Table, name string) *plugin.Column {
	for _, col := range table.Columns {
		if col.Name == name {
			return col
		}
	}

	return nil
}

// graphRelationship returns connected resources for a given resource. For example the members of a group.
func graphRelationship(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, relationship string, idsOnly bool) (any, error) {
	spec := graphTables[d.Table.Name]

	row := h.Item.(graphObject)
	id, _ := row["id"].(string)
	if id == "" {
		return nil, fmt.Errorf("relationship resource has no id")
	}

	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var ids = []string{}
	var members = []graphObject{}

	err = client.list(ctx, spec.endpoint, spec.path+"/"+url.PathEscape(id)+"/"+relationship, url.Values{"$select": {"id"}}, func(row graphObject) bool {
		if idsOnly {
			if id, ok := row["id"].(string); ok {
				ids = append(ids, id)
			}
		} else {
			members = append(members, graphObject{
				"id":          row["id"],
				"@odata.type": row["@odata.type"],
			})
		}

		return true
	})
	if err != nil {
		return nil, err
	}

	if idsOnly {
		return ids, nil
	}

	return members, nil
}

func getGraphGroupSubscription(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	row := h.Item.(graphObject)
	id, _ := row["id"].(string)
	result, err := client.get(ctx, graphDefault, "groups/"+url.PathEscape(id), url.Values{"$select": {"isSubscribedByMail"}})
	if err != nil {
		return nil, err
	}

	return result["isSubscribedByMail"], nil
}

func getGraphManager(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	row := h.Item.(graphObject)
	id, _ := row["id"].(string)
	result, err := client.get(ctx, graphDefault, "users/"+url.PathEscape(id)+"/manager", url.Values{"$select": {"id"}})
	if err != nil {
		if getErrorObject(err).StatusCode == 404 {
			return nil, nil
		}

		return nil, err
	}

	return result["id"], nil
}

func toCamelCase(s string) string {
	var n strings.Builder
	n.Grow(len(s))

	lower := cases.Lower(language.English)
	title := cases.Title(language.English)

	capNext := false

	for _, v := range s {
		if v == '_' || v == ' ' || v == '-' || v == '.' {
			if n.Len() > 0 {
				capNext = true
			}
			continue
		}

		if n.Len() == 0 {
			n.WriteString(lower.String(string(v)))
			continue
		}

		if capNext {
			n.WriteString(title.String(string(v)))
			capNext = false
			continue
		}

		n.WriteString(lower.String(string(v)))
	}

	return n.String()
}
