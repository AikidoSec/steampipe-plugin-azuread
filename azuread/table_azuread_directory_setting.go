package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getAzureAdDirectorySetting(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	spec := graphTables[d.Table.Name]
	path, err := graphTablePath(spec, d, true)
	if err != nil {
		return nil, err
	}

	row, err := client.get(ctx, spec.endpoint, path, graphQuery(spec, d, true))
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	for _, setting := range flattenSettings(row) {
		if setting["name"] == name {
			return setting, nil
		}
	}

	return nil, nil
}

func listAzureAdDirectorySettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	spec := graphTables[d.Table.Name]
	err = client.list(ctx, spec.endpoint, spec.path, graphQuery(spec, d, false), func(row graphObject) bool {
		for _, setting := range flattenSettings(row) {
			d.StreamListItem(ctx, setting)
			if d.RowsRemaining(ctx) == 0 {
				return false
			}
		}

		return d.RowsRemaining(ctx) != 0
	})

	return nil, err
}

func flattenSettings(row graphObject) []graphObject {
	values, _ := row["values"].([]any)
	result := []graphObject{}

	for _, value := range values {
		if setting, ok := value.(map[string]any); ok {
			result = append(result, graphObject{
				"id":          row["id"],
				"displayName": row["displayName"],
				"templateId":  row["templateId"],
				"name":        setting["name"],
				"value":       setting["value"],
			})
		}
	}

	return result
}

func tableAzureAdDirectorySetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_setting",
		Description: "Represents the configurations that can be used to customize the tenant-wide and object-specific restrictions and allowed behavior",
		Get: &plugin.GetConfig{
			Hydrate:    getAzureAdDirectorySetting,
			KeyColumns: plugin.AllColumns([]string{"id", "name"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureAdDirectorySettings,
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
