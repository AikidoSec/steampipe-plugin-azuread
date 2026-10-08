package azuread

import (
	"context"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getSelfServiceGroupManagementFromGraph(ctx context.Context, client *GraphClient) (graphObject, error) {
	policy, err := client.get(ctx, graphDefault, "policies/authorizationPolicy", nil)
	if err != nil {
		return nil, err
	}

	permissions, _ := policy["defaultUserRolePermissions"].(map[string]any)
	settings, err := policySettings(ctx, client, graphBeta, "settings", "Group.Unified")
	if err != nil {
		return nil, err
	}

	row := graphObject{
		"dataSource": "graph",
		"sourcePayload": graphObject{
			"authorizationPolicy": policy,
			"groupSettings":       settings,
		},
		"usersCanCreateSecurityGroups": permissions["allowedToCreateSecurityGroups"],
	}

	for _, setting := range settings {
		values := policySettingValues(setting)

		if v, ok := values["EnableGroupCreation"]; ok {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, fmt.Errorf("invalid EnableGroupCreation: %w", err)
			}

			row["usersCanCreateMicrosoft365Groups"] = b
		}

		if v, ok := values["GroupCreationAllowedGroupId"]; ok {
			row["groupCreationAllowedGroupId"] = v
		}
	}

	return row, nil
}

func tableAzureAdSelfServiceGroupManagement() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_self_service_group_management",
		Description: "Self-service group management settings from the internal portal when configured, otherwise partial Microsoft Graph coverage.",
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
				Name:        "users_can_create_security_groups",
				Description: "Graph default-user security group creation permission.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("usersCanCreateSecurityGroups"),
			},
			{
				Name:        "users_can_create_microsoft365_groups",
				Description: "Graph Microsoft 365 group creation setting.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("usersCanCreateMicrosoft365Groups"),
			},
			{
				Name:        "group_creation_allowed_group_id",
				Description: "Graph group allowed to create Microsoft 365 groups.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("groupCreationAllowedGroupId"),
			},
			{
				Name:        "self_service_group_management_enabled",
				Description: "Whether self-service group management is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("selfServiceGroupManagementEnabled"),
			},
			{
				Name:        "security_groups_enabled",
				Description: "Whether security groups are enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("securityGroupsEnabled"),
			},
			{
				Name:        "users_can_manage_security_groups",
				Description: "Who can manage security groups.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("usersCanManageSecurityGroups"),
			},
			{
				Name:        "office365_groups_enabled",
				Description: "Whether Microsoft 365 groups are enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("office365GroupsEnabled"),
			},
			{
				Name:        "users_can_manage_office_groups",
				Description: "Who can manage Microsoft 365 groups.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("usersCanManageOfficeGroups"),
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
