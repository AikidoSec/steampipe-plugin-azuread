package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdServicePrincipalAppRoleAssignment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_service_principal_app_role_assignment",
		Description: "Represents an application role assigned to a service principal.",
		Get: &plugin.GetConfig{
			Hydrate: getGraphTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "service_principal_id",
					Require: plugin.Required,
				},
				{
					Name:    "id",
					Require: plugin.Required,
				},
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{
					Name:    "service_principal_id",
					Require: plugin.Required,
				},

				// Other fields for filtering OData
				{
					Name:    "resource_id",
					Require: plugin.Optional,
				},
				{
					Name:    "principal_display_name",
					Require: plugin.Optional,
				},
			},
		},

		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "A unique identifier for the appRoleAssignment key.",
				Transform:   graphField("id"),
			},
			{
				Name:        "app_role_id",
				Type:        proto.ColumnType_STRING,
				Description: "The identifier (id) for the app role which is assigned to the principal. This app role must be exposed in the appRoles property on the resource application's service principal (resourceId). If the resource application has not declared any app roles, a default app role ID of 00000000-0000-0000-0000-000000000000 can be specified to signal that the principal is assigned to the resource app without any specific app roles.",
				Transform:   graphField("appRoleId"),
			},
			{
				Name:        "resource_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier (id) for the resource service principal for which the assignment is made.",
				Transform:   graphField("resourceId"),
			},
			{
				Name:        "resource_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name of the resource app's service principal to which the assignment is made.",
				Transform:   graphField("resourceDisplayName"),
			},

			// Other fields
			{
				Name:        "created_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
				Transform:   graphField("createdDateTime"),
			},
			{
				Name:        "deleted_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time when the app role assignment was deleted. Always null for an appRoleAssignment object that hasn't been deleted.",
				Transform:   graphField("deletedDateTime"),
			},

			{
				Name:        "principal_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier (id) for the user, security group, or service principal being granted the app role.",
				Transform:   graphField("principalId"),
			},
			{
				Name:        "principal_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name of the user, group, or service principal that was granted the app role assignment.",
				Transform:   graphField("principalDisplayName"),
			},
			{
				Name:        "principal_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the assigned principal. This can either be User, Group, or ServicePrincipal.",
				Transform:   graphField("principalType"),
			},

			// Standard columns
			{
				Name:        "service_principal_id",
				Type:        proto.ColumnType_STRING,
				Description: "The identifier (id) of the service principal.",
				Transform:   graphField("resourceId"),
			},
		},
	}
}
