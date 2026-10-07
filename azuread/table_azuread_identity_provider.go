package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdIdentityProvider(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_identity_provider",
		Description: "Represents an Azure Active Directory (Azure AD) identity provider.",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_UnsupportedQuery", "Invalid filter clause"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "filter",
					Require: plugin.Optional,
				},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the identity provider.",
				Transform:   graphField("id"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name of the identity provider.",
				Transform:   graphField("displayName"),
			},

			// Other fields
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The identity provider type is a required field. For B2B scenario: Google, Facebook. For B2C scenario: Microsoft, Google, Amazon, LinkedIn, Facebook, GitHub, Twitter, Weibo, QQ, WeChat, OpenIDConnect.",
				Transform:   graphField("identityProviderType"),
			},
			{
				Name:        "client_id",
				Type:        proto.ColumnType_STRING,
				Description: "The client ID for the application. This is the client ID obtained when registering the application with the identity provider.",
				Transform:   graphField("clientId"),
			},
			{
				Name:        "client_secret",
				Type:        proto.ColumnType_STRING,
				Description: "The client secret for the application. This is the client secret obtained when registering the application with the identity provider. This is write-only. A read operation will return ****.",
				Transform:   graphField("clientSecret"),
			},
			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
				Description: "Odata query to search for resources.",
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
