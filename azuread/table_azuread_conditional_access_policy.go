package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdConditionalAccessPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_conditional_access_policy",
		Description: "Represents an Azure Active Directory (Azure AD) Conditional Access Policy.",
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
					Name:    "state",
					Require: plugin.Optional,
				},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "All conditions that determine when the policy applies.",
				Transform:   graphField("conditions"),
			},
			{
				Name:        "grant_controls",
				Type:        proto.ColumnType_JSON,
				Description: "All access grant controls for the policy.",
				Transform:   graphField("grantControls"),
			},
			{
				Name:        "session_controls",
				Type:        proto.ColumnType_JSON,
				Description: "All session controls for the policy.",
				Transform:   graphField("sessionControls"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the identifier of a conditionalAccessPolicy object.",
				Transform:   graphField("id"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies a display name for the conditionalAccessPolicy object.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "state",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the state of the conditionalAccessPolicy object. Possible values are: enabled, disabled, enabledForReportingButNotEnforced.",
				Transform:   graphField("state"),
			},
			{
				Name:        "created_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The create date of the conditional access policy.",
				Transform:   graphField("createdDateTime"),
			},
			{
				Name:        "modified_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The modification date of the conditional access policy.",
				Transform:   graphField("modifiedDateTime"),
			},
			{
				Name:        "operator",
				Type:        proto.ColumnType_STRING,
				Description: "Defines the relationship of the grant controls. Possible values: AND, OR.",
				Transform:   graphField("grantControls.operator"),
			},

			// Json fields
			{
				Name:        "applications",
				Type:        proto.ColumnType_JSON,
				Description: "Applications and user actions included in and excluded from the policy.",
				Transform:   graphField("conditions.applications"),
			},
			{
				Name:        "application_enforced_restrictions",
				Type:        proto.ColumnType_JSON,
				Description: "Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control.",
				Transform:   graphField("sessionControls.applicationEnforcedRestrictions"),
			},
			{
				Name:        "built_in_controls",
				Type:        proto.ColumnType_JSON,
				Description: "List of values of built-in controls required by the policy. Possible values: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.",
				Transform:   graphField("grantControls.builtInControls"),
			},
			{
				Name:        "authentication_strength",
				Type:        proto.ColumnType_JSON,
				Description: "List combinations of authentication methods allowed by the policy. For example: password, Federated Multi-Factor, FIDO2 security key",
				Transform:   graphField("grantControls.authenticationStrength.allowedCombinations"),
			},
			{
				Name:        "client_app_types",
				Type:        proto.ColumnType_JSON,
				Description: "Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, easSupported, other.",
				Transform:   graphField("conditions.clientAppTypes"),
			},
			{
				Name:        "custom_authentication_factors",
				Type:        proto.ColumnType_JSON,
				Description: "List of custom controls IDs required by the policy.",
				Transform:   graphField("grantControls.customAuthenticationFactors"),
			},
			{
				Name:        "cloud_app_security",
				Type:        proto.ColumnType_JSON,
				Description: "Session control to apply cloud app security.",
				Transform:   graphField("sessionControls.cloudAppSecurity"),
			},
			{
				Name:        "locations",
				Type:        proto.ColumnType_JSON,
				Description: "Locations included in and excluded from the policy.",
				Transform:   graphField("conditions.locations"),
			},
			{
				Name:        "persistent_browser",
				Type:        proto.ColumnType_JSON,
				Description: "Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly.",
				Transform:   graphField("sessionControls.persistentBrowser"),
			},
			{
				Name:        "platforms",
				Type:        proto.ColumnType_JSON,
				Description: "Platforms included in and excluded from the policy.",
				Transform:   graphField("conditions.platforms"),
			},
			{
				Name:        "sign_in_frequency",
				Type:        proto.ColumnType_JSON,
				Description: "Session control to enforce signin frequency.",
				Transform:   graphField("sessionControls.signInFrequency"),
			},
			{
				Name:        "sign_in_risk_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.",
				Transform:   graphField("conditions.signInRiskLevels"),
			},
			{
				Name:        "terms_of_use",
				Type:        proto.ColumnType_JSON,
				Description: "List of terms of use IDs required by the policy.",
				Transform:   graphField("grantControls.termsOfUse"),
			},
			{
				Name:        "users",
				Type:        proto.ColumnType_JSON,
				Description: "Users, groups, and roles included in and excluded from the policy.",
				Transform:   graphField("conditions.users"),
			},
			{
				Name:        "user_risk_levels",
				Type:        proto.ColumnType_JSON,
				Description: "User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.",
				Transform:   graphField("conditions.userRiskLevels"),
			},
			{
				Name:        "disable_resilience_defaults",
				Type:        proto.ColumnType_BOOL,
				Description: "Session control that determines whether it is acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage or not.",
				Transform:   graphField("sessionControls.disableResilienceDefaults"),
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
