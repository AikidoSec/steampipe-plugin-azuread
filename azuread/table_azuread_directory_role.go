package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdDirectoryRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role",
		Description: "Represents an Azure Active Directory (Azure AD) directory role.",
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
				Name:        "deleted_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the directory role was deleted, if applicable.",
				Transform:   graphField("deletedDateTime"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the directory role.",
				Transform:   graphField("id"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description for the directory role.",
				Transform:   graphField("description"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name for the directory role.",
				Transform:   graphField("displayName"),
			},

			// Other fields
			{
				Name:        "role_template_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the directoryRoleTemplate that this role is based on. The property must be specified when activating a directory role in a tenant with a POST operation. After the directory role has been activated, the property is read only.",
				Transform:   graphField("roleTemplateId"),
			},

			// Json fields
			{
				Name:        "member_ids",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGraphMembers,
				Transform:   transform.FromValue(),
				Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object.",
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
