package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdDirectoryRoleTemplate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_template",
		Description: "Represents a directory role template in Azure Active Directory (Azure AD). A directory role template specifies the property values of a directory role.",
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
				Description: "The unique identifier for the directory role template.",
				Transform:   graphField("id"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description to set for the directory role.",
				Transform:   graphField("description"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name to set for the directory role.",
				Transform:   graphField("displayName"),
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
