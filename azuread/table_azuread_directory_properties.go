package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getDirectoryPropertiesFromGraph(ctx context.Context, client *GraphClient) (graphObject, error) {
	policy, err := client.get(ctx, graphDefault, "policies/authorizationPolicy", nil)
	if err != nil {
		return nil, err
	}

	permissions, _ := policy["defaultUserRolePermissions"].(map[string]any)

	return graphObject{
		"dataSource":           "graph",
		"sourcePayload":        graphObject{"authorizationPolicy": policy},
		"usersCanRegisterApps": permissions["allowedToCreateApps"],
		"allowInvitesFrom":     policy["allowInvitesFrom"],
	}, nil
}

func tableAzureAdDirectoryProperties() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_properties",
		Description: "Directory properties from the internal portal when configured, otherwise partial Microsoft Graph coverage.",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Directory object identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("objectId"),
			},
			{
				Name:        "display_name",
				Description: "Directory display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("displayName"),
			},
			{
				Name:        "allow_invites_from",
				Description: "Graph guest invitation policy, retaining its role-based scope.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("allowInvitesFrom"),
			},
			{
				Name:        "users_can_register_apps",
				Description: "Whether users can register applications.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("usersCanRegisterApps"),
			},
			{
				Name:        "users_can_add_external_users",
				Description: "Whether users can add external users.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("usersCanAddExternalUsers"),
			},
			{
				Name:        "restrict_directory_access",
				Description: "Whether directory access is restricted.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("restrictDirectoryAccess"),
			},
			{
				Name:        "restrict_non_admin_users",
				Description: "Whether non-administrator access is restricted.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("restrictNonAdminUsers"),
			},
			{
				Name:        "data_source",
				Description: "Source: graph (partial coverage) or internal_portal.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("dataSource"),
			},
			{
				Name:        "raw",
				Type:        proto.ColumnType_JSON,
				Description: "Complete API response, including properties not exposed as individual columns.",
				Transform:   graphField("sourcePayload"),
			},
		}),
	}
}
