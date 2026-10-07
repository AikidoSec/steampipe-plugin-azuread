package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdAuthorizationPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_authorization_policy",
		Description: "Represents a policy that can control Azure Active Directory authorization settings.",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "allow_user_consent_for_risky_apps",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether users can consent to risky applications.",
				Transform:   graphField("allowUserConsentForRiskyApps"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Display name for this policy.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the authorization policy.",
				Transform:   graphField("id"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Description of this policy.",
				Transform:   graphField("description"),
			},

			// Other fields
			{
				Name:        "allowed_to_sign_up_email_based_subscriptions",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether users can sign up for email based subscriptions.",
				Transform:   graphField("allowedToSignUpEmailBasedSubscriptions"),
			},
			{
				Name:        "allowed_to_use_sspr",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the Self-Serve Password Reset feature can be used by users on the tenant.",
				Transform:   graphField("allowedToUseSSPR"),
			},
			{
				Name:        "allowed_email_verified_users_to_join_organization",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether a user can join the tenant by email validation.",
				Transform:   graphField("allowEmailVerifiedUsersToJoinOrganization"),
			},
			{
				Name:        "allow_invites_from",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates who can invite external users to the organization. Possible values are: none, adminsAndGuestInviters, adminsGuestInvitersAndAllMembers, everyone.",
				Transform:   graphField("allowInvitesFrom"),
			},
			{
				Name:        "block_msol_powershell",
				Type:        proto.ColumnType_BOOL,
				Description: "To disable the use of MSOL PowerShell set this property to true. This will also disable user-based access to the legacy service endpoint used by MSOL PowerShell. This does not affect Azure AD Connect or Microsoft Graph.",
				Transform:   graphField("blockMsolPowerShell"),
			},
			{
				Name:        "guest_user_role_id",
				Type:        proto.ColumnType_STRING,
				Description: "Represents role templateId for the role that should be granted to guest user.",
				Transform:   graphField("guestUserRoleId"),
			},

			// JSON fields
			{
				Name:        "default_user_role_permissions",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies certain customizable permissions for default user role.",
				Transform:   graphField("defaultUserRolePermissions"),
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
