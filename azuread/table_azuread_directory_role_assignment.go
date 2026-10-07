package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdDirectoryRoleAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_assignment",
		Description: "Represents the role assignments for Azure AD resources.",
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
				Description: "The unique identifier for the role assignment.",
				Transform:   graphField("id"),
			},
			{
				Name:        "principal_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the principal that's in the scope of the role assignment.",
				Transform:   graphField("principalId"),
			},
			{
				Name:        "role_definition_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the role definition that's in the scope of the role assignment.",
				Transform:   graphField("roleDefinitionId"),
			},
			{
				Name:        "directory_scope_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the directory scope that's in the scope of the role assignment.",
				Transform:   graphField("directoryScopeId"),
			},
			{
				Name:        "app_scope_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the app scope that's in the scope of the role assignment.",
				Transform:   graphField("appScopeId"),
			},
			{
				Name:        "condition",
				Type:        proto.ColumnType_STRING,
				Description: "The condition which describes the circumstances under which the role assignment is valid.",
				Transform:   graphField("condition"),
			},

			// JSON fields
			{
				Name:        "principal",
				Type:        proto.ColumnType_JSON,
				Description: "The principal (user, group, or service principal) that the role is assigned to.",
				Transform:   transform.FromP(graphRoleReference, "principal"),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphTitle("id", "principalId"),
			},
		}),
	}
}
