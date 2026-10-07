package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_user",
		Description: "Represents an Azure AD user account.",
		Get: &plugin.GetConfig{
			Hydrate: getGraphTable,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Request_ResourceNotFound", "Invalid object identifier"}),
			},
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{
					Name:    "user_principal_name",
					Require: plugin.Optional,
				},
				{
					Name:    "filter",
					Require: plugin.Optional,
				},

				// Other fields for filtering OData
				{
					Name:    "user_type",
					Require: plugin.Optional,
				},
				{
					Name:      "account_enabled",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "surname",
					Require: plugin.Optional,
				},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "custom_security_attributes",
				Type:        proto.ColumnType_JSON,
				Description: "Custom security attribute assignments, subject to attribute permissions.",
				Transform:   graphField("customSecurityAttributes"),
			},
			{
				Name:        "identities",
				Type:        proto.ColumnType_JSON,
				Description: "Identities used to sign in to this account.",
				Transform:   graphField("identities"),
			},
			{
				Name:        "last_password_change_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the user's password was last changed.",
				Transform:   graphField("lastPasswordChangeDateTime"),
			},
			{
				Name:        "security_identifier",
				Type:        proto.ColumnType_STRING,
				Description: "The user's security identifier.",
				Transform:   graphField("securityIdentifier"),
			},
			{
				Name:        "deleted_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the user was deleted, if applicable.",
				Transform:   graphField("deletedDateTime"),
			},
			{
				Name:        "registered_device_ids",
				Type:        proto.ColumnType_JSON,
				Description: "Identifiers of devices registered by the user.",
				Hydrate:     getGraphRegisteredDevices,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the user. Should be treated as an opaque identifier.",
				Transform:   graphField("id"),
			},
			{
				Name:        "user_principal_name",
				Type:        proto.ColumnType_STRING,
				Description: "Principal email of the active directory user.",
				Transform:   graphField("userPrincipalName"),
			},
			{
				Name:        "account_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "True if the account is enabled; otherwise, false.",
				Transform:   graphField("accountEnabled"),
			},
			{
				Name:        "user_type",
				Type:        proto.ColumnType_STRING,
				Description: "A string value that can be used to classify user types in your directory.",
				Transform:   graphField("userType"),
			},
			{
				Name:        "given_name",
				Type:        proto.ColumnType_STRING,
				Description: "The given name (first name) of the user.",
				Transform:   graphField("givenName"),
			},
			{
				Name:        "surname",
				Type:        proto.ColumnType_STRING,
				Description: "Family name or last name of the active directory user.",
				Transform:   graphField("surname"),
			},

			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
				Description: "Odata query to search for resources.",
			},

			// Other fields
			{
				Name:        "created_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time at which the user was created.",
				Transform:   graphField("createdDateTime"),
			},
			{
				Name:        "mail",
				Type:        proto.ColumnType_STRING,
				Description: "The SMTP address for the user, for example, jeff@contoso.onmicrosoft.com.",
				Transform:   graphField("mail"),
			},
			{
				Name:        "mail_nickname",
				Type:        proto.ColumnType_STRING,
				Description: "The mail alias for the user.",
				Transform:   graphField("mailNickname"),
			},
			{
				Name:        "password_policies",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two may be specified together; for example: DisablePasswordExpiration, DisableStrongPassword.",
				Transform:   graphField("passwordPolicies"),
			},
			{
				Name:        "sign_in_sessions_valid_from_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).",
				Transform:   graphField("signInSessionsValidFromDateTime"),
			},
			{
				Name:        "usage_location",
				Type:        proto.ColumnType_STRING,
				Description: "A two letter country code (ISO standard 3166), required for users that will be assigned licenses due to legal requirement to check for availability of services in countries.",
				Transform:   graphField("usageLocation"),
			},
			{
				Name:        "external_user_state",
				Type:        proto.ColumnType_STRING,
				Description: "For an external user invited to the tenant using the invitation API, this property represents the invited user's invitation status",
				Transform:   graphField("externalUserState"),
			},

			// Job Information
			{
				Name:        "employee_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier assigned to the employee.",
				Transform:   graphField("employeeId"),
			},
			{
				Name:        "employee_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of employment (e.g., full-time, part-time, contractor).",
				Transform:   graphField("employeeType"),
			},
			{
				Name:        "company_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the company the user is associated with.",
				Transform:   graphField("companyName"),
			},
			{
				Name:        "job_title",
				Type:        proto.ColumnType_STRING,
				Description: "The job title of the user.",
				Transform:   graphField("jobTitle"),
			},
			{
				Name:        "department",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the department in which the user works.",
				Transform:   graphField("department"),
			},
			{
				Name:        "office_location",
				Type:        proto.ColumnType_STRING,
				Description: "The physical location of the user's office.",
				Transform:   graphField("officeLocation"),
			},
			{
				Name:        "manager",
				Type:        proto.ColumnType_STRING,
				Description: "The manager of the user.",
				Hydrate:     getGraphManager,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "employee_hire_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date when the user was hired.",
				Transform:   graphField("employeeHireDate"),
			},

			// On-premises
			{
				Name:        "on_premises_distinguished_name",
				Type:        proto.ColumnType_STRING,
				Description: "The distinguished name of the user in the on-premises Active Directory.",
				Transform:   graphField("onPremisesDistinguishedName"),
			},
			{
				Name:        "on_premises_domain_name",
				Type:        proto.ColumnType_STRING,
				Description: "The domain name of the user in the on-premises Active Directory.",
				Transform:   graphField("onPremisesDomainName"),
			},
			{
				Name:        "on_premises_immutable_id",
				Type:        proto.ColumnType_STRING,
				Description: "Used to associate an on-premises Active Directory user account with their Azure AD user object.",
				Transform:   graphField("onPremisesImmutableId"),
			},
			{
				Name:        "on_premises_sam_account_name",
				Type:        proto.ColumnType_STRING,
				Description: "The Security Account Manager (SAM) account name of the user in the on-premises Active Directory.",
				Transform:   graphField("onPremisesSamAccountName"),
			},
			{
				Name:        "on_premises_security_identifier",
				Type:        proto.ColumnType_STRING,
				Description: "The security identifier (SID) of the user in the on-premises Active Directory.",
				Transform:   graphField("onPremisesSecurityIdentifier"),
			},
			{
				Name:        "on_premises_user_principal_name",
				Type:        proto.ColumnType_STRING,
				Description: "The User Principal Name (UPN) of the user in the on-premises Active Directory.",
				Transform:   graphField("onPremisesUserPrincipalName"),
			},
			{
				Name:        "on_premises_sync_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the user is synchronized with on-premises Active Directory.",
				Transform:   graphField("onPremisesSyncEnabled"),
			},
			{
				Name:        "on_premises_last_sync_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time when the user's information was last synchronized with the on-premises Active Directory.",
				Transform:   graphField("onPremisesLastSyncDateTime"),
			},

			// Json fields
			{
				Name:        "member_of",
				Type:        proto.ColumnType_JSON,
				Description: "A list the groups and directory roles that the user is a direct member of.",
				Hydrate:     getGraphMemberOf,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "im_addresses",
				Type:        proto.ColumnType_JSON,
				Description: "The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user.",
				Transform:   graphField("imAddresses"),
			},
			{
				Name:        "other_mails",
				Type:        proto.ColumnType_JSON,
				Description: "A list of additional email addresses for the user.",
				Transform:   graphField("otherMails"),
			},
			{
				Name:        "password_profile",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies the password profile for the user. The profile contains the user’s password. This property is required when a user is created.",
				Transform:   graphField("passwordProfile"),
			},
			{
				Name:        "sign_in_activity",
				Type:        proto.ColumnType_JSON,
				Description: "Get the last signed-in date and request ID of the sign-in for a given user",
				Transform:   transform.From(graphSignInActivity),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphTitle("displayName", "userPrincipalName"),
			},
		}),
	}
}
