package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdDirectoryRoleDefinition(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_definition",
		Description: "Represents the role definitions for Azure AD directory resources.",
		Get: &plugin.GetConfig{
			Hydrate: getGraphTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the role definition.",
				Transform:   graphField("id"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name for the role definition.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description for the role definition.",
				Transform:   graphField("description"),
			},
			{
				Name:        "template_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the role template.",
				Transform:   graphField("templateId"),
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Description: "The version of the role definition.",
				Transform:   graphField("version"),
			},
			{
				Name:        "is_built_in",
				Type:        proto.ColumnType_BOOL,
				Description: "Flag indicating if the role definition is built in.",
				Transform:   graphField("isBuiltIn"),
			},
			{
				Name:        "is_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Flag indicating whether the role is enabled for assignment.",
				Transform:   graphField("isEnabled"),
			},

			// JSON fields
			{
				Name:        "resource_scopes",
				Type:        proto.ColumnType_JSON,
				Description: "List of scopes that the role definition applies to.",
				Transform:   graphField("resourceScopes"),
			},
			{
				Name:        "role_permissions",
				Type:        proto.ColumnType_JSON,
				Description: "List of permissions included in this role.",
				Transform:   graphField("rolePermissions"),
			},
			{
				Name:        "inherits_permissions_from",
				Type:        proto.ColumnType_JSON,
				Description: "Read-only collection of role definitions that the given role definition inherits from.",
				Transform:   graphField("inheritsPermissionsFrom"),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphTitle("displayName", "id"),
			},
		}),
	}
}
