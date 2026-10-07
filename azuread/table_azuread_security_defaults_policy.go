package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdSecurityDefaultsPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_security_defaults_policy",
		Description: "Represents the Azure Active Directory security defaults policy",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Display name for this policy.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Identifier for this policy.",
				Transform:   graphField("id"),
			},
			{
				Name:        "is_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "If set to true, Azure Active Directory security defaults is enabled for the tenant.",
				Transform:   graphField("isEnabled"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Description for this policy.",
				Transform:   graphField("description"),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphField("displayName"),
			},
		}),
	}
}
