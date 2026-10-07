package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdServicePrincipal(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_service_principal",
		Description: "Represents an Azure Active Directory (Azure AD) service principal.",
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
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{
					Name:    "app_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:      "account_enabled",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
				{
					Name:    "service_principal_type",
					Require: plugin.Optional,
				},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "deleted_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the service principal was deleted, if applicable.",
				Transform:   graphField("deletedDateTime"),
			},
			{
				Name:        "disabled_by_microsoft_status",
				Type:        proto.ColumnType_STRING,
				Description: "The service principal's Microsoft disablement status.",
				Transform:   graphField("disabledByMicrosoftStatus"),
			},
			{
				Name:        "preferred_single_sign_on_mode",
				Type:        proto.ColumnType_STRING,
				Description: "The preferred single sign-on mode.",
				Transform:   graphField("preferredSingleSignOnMode"),
			},
			{
				Name:        "preferred_token_signing_key_thumbprint",
				Type:        proto.ColumnType_STRING,
				Description: "The preferred token signing key thumbprint.",
				Transform:   graphField("preferredTokenSigningKeyThumbprint"),
			},
			{
				Name:        "token_encryption_key_id",
				Type:        proto.ColumnType_STRING,
				Description: "The public key used to encrypt tokens.",
				Transform:   graphField("tokenEncryptionKeyId"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the service principal.",
				Transform:   graphField("id"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name for the service principal.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "app_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the associated application (its appId property).",
				Transform:   graphField("appId"),
			},

			// Other fields
			{
				Name:        "account_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "true if the service principal account is enabled; otherwise, false.",
				Transform:   graphField("accountEnabled"),
			},
			{
				Name:        "app_display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name exposed by the associated application.",
				Transform:   graphField("appDisplayName"),
			},
			{
				Name:        "app_owner_organization_id",
				Type:        proto.ColumnType_STRING,
				Description: "Contains the tenant id where the application is registered. This is applicable only to service principals backed by applications.",
				Transform:   graphField("appOwnerOrganizationId"),
			},
			{
				Name:        "app_role_assignment_required",
				Type:        proto.ColumnType_BOOL,
				Description: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false.",
				Transform:   graphField("appRoleAssignmentRequired"),
			},
			{
				Name:        "service_principal_type",
				Type:        proto.ColumnType_STRING,
				Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Azure AD internally.",
				Transform:   graphField("servicePrincipalType"),
			},
			{
				Name:        "sign_in_audience",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the Microsoft accounts that are supported for the current application. Supported values are: AzureADMyOrg, AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, PersonalMicrosoftAccount.",
				Transform:   graphField("signInAudience"),
			},
			{
				Name:        "app_description",
				Type:        proto.ColumnType_STRING,
				Description: "The description exposed by the associated application.",
				Transform:   graphField("appDescription"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Free text field to provide an internal end-user facing description of the service principal.",
				Transform:   graphField("description"),
			},
			{
				Name:        "login_url",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the URL where the service provider redirects the user to Azure AD to authenticate. Azure AD uses the URL to launch the application from Microsoft 365 or the Azure AD My Apps. When blank, Azure AD performs IdP-initiated sign-on for applications configured with SAML-based single sign-on.",
				Transform:   graphField("loginUrl"),
			},
			{
				Name:        "logout_url",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the URL that will be used by Microsoft's authorization service to logout an user using OpenId Connect front-channel, back-channel or SAML logout protocols.",
				Transform:   graphField("logoutUrl"),
			},

			// JSON fields
			{
				Name:        "add_ins",
				Type:        proto.ColumnType_JSON,
				Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts.",
				Transform:   graphField("addIns"),
			},
			{
				Name:        "alternative_names",
				Type:        proto.ColumnType_JSON,
				Description: "Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities.",
				Transform:   graphField("alternativeNames"),
			},
			{
				Name:        "app_roles",
				Type:        proto.ColumnType_JSON,
				Description: "The roles exposed by the application which this service principal represents.",
				Transform:   graphField("appRoles"),
			},
			{
				Name:        "info",
				Type:        proto.ColumnType_JSON,
				Description: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs.",
				Transform:   graphField("info"),
			},
			{
				Name:        "key_credentials",
				Type:        proto.ColumnType_JSON,
				Description: "The collection of key credentials associated with the service principal.",
				Transform:   graphField("keyCredentials"),
			},
			{
				Name:        "notification_email_addresses",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications.",
				Transform:   graphField("notificationEmailAddresses"),
			},
			{
				Name:        "owner_ids",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGraphOwners,
				Transform:   transform.FromValue(),
				Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object.",
			},
			{
				Name:        "password_credentials",
				Type:        proto.ColumnType_JSON,
				Description: "Represents a password credential associated with a service principal.",
				Transform:   graphField("passwordCredentials"),
			},
			{
				Name:        "oauth2_permission_scopes",
				Type:        proto.ColumnType_JSON,
				Description: "The published permission scopes.",
				Transform:   graphField("oauth2PermissionScopes"),
			},
			{
				Name:        "reply_urls",
				Type:        proto.ColumnType_JSON,
				Description: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application.",
				Transform:   graphField("replyUrls"),
			},
			{
				Name:        "service_principal_names",
				Type:        proto.ColumnType_JSON,
				Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD.",
				Transform:   graphField("servicePrincipalNames"),
			},
			{
				Name:        "tags_src",
				Type:        proto.ColumnType_JSON,
				Description: "Custom strings that can be used to categorize and identify the service principal.",
				Transform:   graphField("tags"),
			},

			// Standard columns
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(graphTags),
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
