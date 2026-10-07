package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdDirectorySetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_setting",
		Description: "Represents the configurations that can be used to customize the tenant-wide and object-specific restrictions and allowed behavior",
		Get: &plugin.GetConfig{
			Hydrate:    getGraphTable,
			KeyColumns: plugin.AllColumns([]string{"id", "name"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Display name of this group of settings, which comes from the associated template.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Unique identifier for these settings.",
				Transform:   graphField("id"),
			},
			{
				Name:        "template_id",
				Type:        proto.ColumnType_STRING,
				Description: "Unique identifier for the template used to create this group of settings.",
				Transform:   graphField("templateId"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the setting.",
				Transform:   graphField("name"),
			},
			{
				Name:        "value",
				Type:        proto.ColumnType_STRING,
				Description: "Value of the setting.",
				Transform:   graphField("value"),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphField("name"),
			},
		}),
	}
}
