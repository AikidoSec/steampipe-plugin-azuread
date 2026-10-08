package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdConditionalAccessNamedLocation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_conditional_access_named_location",
		Description: "Represents an Azure Active Directory (Azure AD) Conditional Access Named Location.",
		Get: &plugin.GetConfig{
			Hydrate: getGraphTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "location_type",
					Require: plugin.Optional,
				},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "raw",
				Type:        proto.ColumnType_JSON,
				Description: "The complete named location response, including location and IP range types.",
				Transform:   transform.FromValue(),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the identifier of a Named Location object.",
				Transform:   graphField("id"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies a display name for the Named Location object.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "location_type",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the type of the Named Location object: IP or Country.",
				Transform:   transform.From(graphLocationType),
			},
			{
				Name:        "created_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The create date of the Named Location object.",
				Transform:   graphField("createdDateTime"),
			},
			{
				Name:        "modified_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The modification date of Named Location object.",
				Transform:   graphField("modifiedDateTime"),
			},
			{
				Name:        "location_info",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies some location information for the Named Location object. Now supported: IP (v4/6 and CIDR/Range), odata_type, IsTrusted (for IP named locations only). Country (and regions, if exist), lookup method, UnkownCountriesAndRegions (for country named locations only).",
				Transform:   transform.From(graphLocationInfo),
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
