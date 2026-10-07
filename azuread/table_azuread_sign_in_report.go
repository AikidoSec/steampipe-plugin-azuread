package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdSignInReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_sign_in_report",
		Description: "Represents an Azure Active Directory (Azure AD) sign-in report.",
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
				Description: "Unique ID representing the sign-in activity.",
				Transform:   graphField("id"),
			},
			{
				Name:        "created_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Date and time (UTC) the sign-in was initiated.",
				Transform:   graphField("createdDateTime"),
			},
			{
				Name:        "user_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Display name of the user that initiated the sign-in.",
				Transform:   graphField("userDisplayName"),
			},
			{
				Name:        "user_principal_name",
				Type:        proto.ColumnType_STRING,
				Description: "User principal name of the user that initiated the sign-in.",
				Transform:   graphField("userPrincipalName"),
			},
			{
				Name:        "user_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the user that initiated the sign-in.",
				Transform:   graphField("userId"),
			},
			{
				Name:        "app_id",
				Type:        proto.ColumnType_STRING,
				Description: "Unique GUID representing the app ID in the Azure Active Directory.",
				Transform:   graphField("appId"),
			},
			{
				Name:        "app_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "App name displayed in the Azure Portal.",
				Transform:   graphField("appDisplayName"),
			},
			{
				Name:        "ip_address",
				Type:        proto.ColumnType_STRING,
				Description: "IP address of the client used to sign in.",
				Transform:   graphField("ipAddress"),
			},
			{
				Name:        "client_app_used",
				Type:        proto.ColumnType_STRING,
				Description: "Identifies the legacy client used for sign-in activity.",
				Transform:   graphField("clientAppUsed"),
			},
			{
				Name:        "correlation_id",
				Type:        proto.ColumnType_STRING,
				Description: "The request ID sent from the client when the sign-in is initiated; used to troubleshoot sign-in activity.",
				Transform:   graphField("correlationId"),
			},
			{
				Name:        "conditional_access_status",
				Type:        proto.ColumnType_STRING,
				Description: "Reports status of an activated conditional access policy. Possible values are: success, failure, notApplied, and unknownFutureValue.",
				Transform:   graphField("conditionalAccessStatus"),
			},
			{
				Name:        "is_interactive",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if a sign-in is interactive or not.",
				Transform:   graphField("isInteractive"),
			},
			{
				Name:        "risk_detail",
				Type:        proto.ColumnType_STRING,
				Description: "Provides the 'reason' behind a specific state of a risky user, sign-in or a risk event. The possible values are: none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, unknownFutureValue.",
				Transform:   graphField("riskDetail"),
			},
			{
				Name:        "risk_level_aggregated",
				Type:        proto.ColumnType_STRING,
				Description: "Aggregated risk level. The possible values are: none, low, medium, high, hidden, and unknownFutureValue.",
				Transform:   graphField("riskLevelAggregated"),
			},
			{
				Name:        "risk_level_during_sign_in",
				Type:        proto.ColumnType_STRING,
				Description: "Risk level during sign-in. The possible values are: none, low, medium, high, hidden, and unknownFutureValue.",
				Transform:   graphField("riskLevelDuringSignIn"),
			},
			{
				Name:        "risk_state",
				Type:        proto.ColumnType_STRING,
				Description: "Reports status of the risky user, sign-in, or a risk event. The possible values are: none, confirmedSafe, remediated, dismissed, atRisk, confirmedCompromised, unknownFutureValue.",
				Transform:   graphField("riskState"),
			},
			{
				Name:        "resource_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the resource the user signed into.",
				Transform:   graphField("resourceDisplayName"),
			},
			{
				Name:        "resource_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the resource that the user signed into.",
				Transform:   graphField("resourceId"),
			},

			// JSON fields
			{
				Name:        "risk_event_types",
				Type:        proto.ColumnType_JSON,
				Description: "Risk event types associated with the sign-in. The possible values are: unlikelyTravel, anonymizedIPAddress, maliciousIPAddress, unfamiliarFeatures, malwareInfectedIPAddress, suspiciousIPAddress, leakedCredentials, investigationsThreatIntelligence, generic, and unknownFutureValue.",
				Transform:   graphField("riskEventTypes"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_JSON,
				Description: "Sign-in status. Includes the error code and description of the error (in case of a sign-in failure).",
				Transform:   graphField("status"),
			},
			{
				Name:        "device_detail",
				Type:        proto.ColumnType_JSON,
				Description: "Device information from where the sign-in occurred; includes device ID, operating system, and browser.",
				Transform:   graphField("deviceDetail"),
			},
			{
				Name:        "location",
				Type:        proto.ColumnType_JSON,
				Description: "Provides the city, state, and country code where the sign-in originated.",
				Transform:   graphField("location"),
			},
			{
				Name:        "applied_conditional_access_policies",
				Type:        proto.ColumnType_JSON,
				Description: "Provides a list of conditional access policies that are triggered by the corresponding sign-in activity.",
				Transform:   graphField("appliedConditionalAccessPolicies"),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphField("id"),
			},
		}),
	}
}
