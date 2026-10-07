package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdAuthenticationMethodPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_authentication_method_policy",
		Description: "Authentication methods policy for the Microsoft Entra tenant.",
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
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The policy description.",
				Transform:   graphField("description"),
			},
			{
				Name:        "last_modified_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the policy was last updated.",
				Transform:   graphField("lastModifiedDateTime"),
			},
			{
				Name:        "policy_migration_state",
				Type:        proto.ColumnType_STRING,
				Description: "Migration state from legacy MFA and SSPR policies.",
				Transform:   graphField("policyMigrationState"),
			},
			{
				Name:        "policy_version",
				Type:        proto.ColumnType_STRING,
				Description: "The policy version.",
				Transform:   graphField("policyVersion"),
			},
			{
				Name:        "reconfirmation_in_days",
				Type:        proto.ColumnType_INT,
				Description: "The reconfirmation interval in days.",
				Transform:   graphField("reconfirmationInDays"),
			},
			{
				Name:        "registration_enforcement",
				Type:        proto.ColumnType_JSON,
				Description: "Authentication method registration requirements.",
				Transform:   graphField("registrationEnforcement"),
			},
			{
				Name:        "authentication_method_configurations",
				Type:        proto.ColumnType_JSON,
				Description: "Settings for each authentication method.",
				Transform:   graphField("authenticationMethodConfigurations"),
			},
			{
				Name:        "additional_data",
				Type:        proto.ColumnType_JSON,
				Description: "Additional policy properties returned by the API.",
				Transform:   transform.From(authenticationMethodAdditionalData),
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

func authenticationMethodAdditionalData(_ context.Context, d *transform.TransformData) (any, error) {
	result := graphObject{}

	for key, value := range d.HydrateItem.(graphObject) {
		switch key {
		case "id", "displayName", "description", "lastModifiedDateTime", "policyMigrationState", "policyVersion", "reconfirmationInDays", "registrationEnforcement", "authenticationMethodConfigurations", "@odata.type":
			continue
		default:
			result[key] = value
		}
	}

	return result, nil
}
