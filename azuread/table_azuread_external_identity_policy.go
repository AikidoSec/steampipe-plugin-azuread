package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdExternalIdentityPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_external_identity_policy",
		Description: "Tenant-wide external identity policy.",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The policy identifier.",
				Transform:   graphField("id"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The policy display name.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "allow_external_identities_to_leave",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether external users can leave the tenant through self-service controls.",
				Transform:   graphField("allowExternalIdentitiesToLeave"),
			},
			{
				Name:        "allow_deleted_identities_data_removal",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether deleted identities' data can be removed.",
				Transform:   graphField("allowDeletedIdentitiesDataRemoval"),
			},
			{
				Name:        "deleted_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the policy was deleted.",
				Transform:   graphField("deletedDateTime"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphTitle("displayName", "id"),
			},
		}),
	}
}
