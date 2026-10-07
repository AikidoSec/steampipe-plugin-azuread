package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdCrossTenantAccessPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_cross_tenant_access_policy",
		Description: "Cross-tenant access policy and its default and partner configurations.",
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
				Name:        "allowed_cloud_endpoints",
				Type:        proto.ColumnType_JSON,
				Description: "Microsoft clouds with which the organization can collaborate.",
				Transform:   graphField("allowedCloudEndpoints"),
			},
			{
				Name:        "default_configuration",
				Type:        proto.ColumnType_JSON,
				Description: "Default configuration for interactions with external organizations.",
				Hydrate:     getCrossTenantDefault,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "partners",
				Type:        proto.ColumnType_JSON,
				Description: "Partner-specific configurations for external organizations.",
				Hydrate:     getCrossTenantPartners,
				Transform:   transform.FromValue(),
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

func getCrossTenantDefault(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	return client.get(ctx, graphDefault, "policies/crossTenantAccessPolicy/default", nil)
}

func getCrossTenantPartners(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	client, err := GetGraphClient(ctx, d)
	if err != nil {
		return nil, err
	}

	partners := []graphObject{}
	err = client.list(ctx, graphDefault, "policies/crossTenantAccessPolicy/partners", nil, func(row graphObject) bool {
		partners = append(partners, row)
		return true
	})
	if err != nil {
		return nil, err
	}

	return partners, nil
}
