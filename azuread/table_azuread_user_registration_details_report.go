package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdUserRegistrationDetailsReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user_registration_details_report",
		Description: "Represents an Azure Active Directory (Azure AD) user-registration-details report.",
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
				Name:        "user_principal_name",
				Type:        proto.ColumnType_STRING,
				Description: "User principal name of the user that initiated the sign-in.",
				Transform:   graphField("userPrincipalName"),
			},
			{
				Name:        "user_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Display name of the user that initiated the sign-in.",
				Transform:   graphField("userDisplayName"),
			},
			{
				Name:        "user_type",
				Type:        proto.ColumnType_STRING,
				Description: "A string value that can be used to classify user types in your directory.",
				Transform:   graphField("userType"),
			},
			{
				Name:        "is_admin",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user has an admin role in the tenant",
				Transform:   graphField("isAdmin"),
			},
			{
				Name:        "is_mfa_registered",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user has registered a strong authentication method for multifactor authentication",
				Transform:   graphField("isMfaRegistered"),
			},
			{
				Name:        "is_mfa_capable",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user has registered a strong authentication method for multifactor authentication",
				Transform:   graphField("isMfaCapable"),
			},
			{
				Name:        "is_sspr_registered",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user has registered the required number of authentication methods for self-service password reset",
				Transform:   graphField("isSsprRegistered"),
			},
			{
				Name:        "is_sspr_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user is allowed to perform self-service password reset by policy",
				Transform:   graphField("isSsprEnabled"),
			},
			{
				Name:        "is_sspr_capable",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user has registered the required number of authentication methods for self-service password reset and the user is allowed to perform self-service password reset by policy",
				Transform:   graphField("isSsprCapable"),
			},
			{
				Name:        "is_passwordless_capable",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user has registered a passwordless strong authentication method",
				Transform:   graphField("isPasswordlessCapable"),
			},
			{
				Name:        "is_system_preferred_authentication_method_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether system preferred authentication method is enabled",
				Transform:   graphField("isSystemPreferredAuthenticationMethodEnabled"),
			},
			{
				Name:        "user_preferred_method_for_secondary_authentication",
				Type:        proto.ColumnType_STRING,
				Description: "The method the user selected as the default second-factor for performing multifactor authentication",
				Transform:   graphField("userPreferredMethodForSecondaryAuthentication"),
			},
			{
				Name:        "last_updated_datetime",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time (UTC) when the report was last updated",
				Transform:   graphField("lastUpdatedDateTime"),
			},

			// JSON fields
			{
				Name:        "methods_registered",
				Type:        proto.ColumnType_JSON,
				Description: "Collection of authentication methods registered",
				Transform:   graphField("methodsRegistered"),
			},
			{
				Name:        "system_preferred_authentication_methods",
				Type:        proto.ColumnType_JSON,
				Description: "Collection of authentication methods that the system determined to be the most secure authentication methods among the registered methods for second factor authentication.",
				Transform:   graphField("systemPreferredAuthenticationMethods"),
			},
		}),
	}
}
