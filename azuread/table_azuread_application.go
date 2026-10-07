package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_application",
		Description: "Represents an Azure Active Directory (Azure AD) application.",
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
					Name:    "app_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "publisher_domain",
					Require: plugin.Optional,
				},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "application_template_id",
				Type:        proto.ColumnType_STRING,
				Description: "The application template identifier.",
				Transform:   graphField("applicationTemplateId"),
			},
			{
				Name:        "app_roles",
				Type:        proto.ColumnType_JSON,
				Description: "Application roles exposed by this application.",
				Transform:   graphField("appRoles"),
			},
			{
				Name:        "certification",
				Type:        proto.ColumnType_JSON,
				Description: "Application certification details.",
				Transform:   graphField("certification"),
			},
			{
				Name:        "deleted_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the application was deleted, if applicable.",
				Transform:   graphField("deletedDateTime"),
			},
			{
				Name:        "disabled_by_microsoft_status",
				Type:        proto.ColumnType_STRING,
				Description: "The application's Microsoft disablement status.",
				Transform:   graphField("disabledByMicrosoftStatus"),
			},
			{
				Name:        "group_membership_claims",
				Type:        proto.ColumnType_STRING,
				Description: "Group membership claims requested in tokens.",
				Transform:   graphField("groupMembershipClaims"),
			},
			{
				Name:        "is_device_only_auth_supported",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether device-only authentication is supported.",
				Transform:   graphField("isDeviceOnlyAuthSupported"),
			},
			{
				Name:        "is_fallback_public_client",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the fallback application type is a public client.",
				Transform:   graphField("isFallbackPublicClient"),
			},
			{
				Name:        "request_signature_verification",
				Type:        proto.ColumnType_JSON,
				Description: "Authentication request signature verification settings.",
				Transform:   graphField("requestSignatureVerification"),
			},
			{
				Name:        "required_resource_access",
				Type:        proto.ColumnType_JSON,
				Description: "Required resource APIs and their delegated permissions and application roles.",
				Transform:   graphField("requiredResourceAccess"),
			},
			{
				Name:        "service_management_reference",
				Type:        proto.ColumnType_STRING,
				Description: "Reference to application management information.",
				Transform:   graphField("serviceManagementReference"),
			},
			{
				Name:        "service_principal_lock_configuration",
				Type:        proto.ColumnType_JSON,
				Description: "Restrictions on editing sensitive service principal properties.",
				Transform:   graphField("servicePrincipalLockConfiguration"),
			},
			{
				Name:        "token_encryption_key_id",
				Type:        proto.ColumnType_STRING,
				Description: "The public key used to encrypt tokens.",
				Transform:   graphField("tokenEncryptionKeyId"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The display name for the application.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the application.",
				Transform:   graphField("id"),
			},
			{
				Name:        "app_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the application that is assigned to an application by Azure AD.",
				Transform:   graphField("appId"),
			},

			// Other fields
			{
				Name:        "created_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time.",
				Transform:   graphField("createdDateTime"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Free text field to provide a description of the application object to end users.",
				Transform:   graphField("description"),
			},
			{
				Name:        "is_authorization_service_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Legacy field unavailable in Microsoft Graph; returned as null.",
				Transform:   graphField("isAuthorizationServiceEnabled"),
			},
			{
				Name:        "oauth2_require_post_response",
				Type:        proto.ColumnType_BOOL,
				Description: "Specifies whether, as part of OAuth 2.0 token requests, Azure AD allows POST requests, as opposed to GET requests. The default is false, which specifies that only GET requests are allowed.",
				Transform:   graphField("oauth2RequiredPostResponse"),
				Default:     false,
			},
			{
				Name:        "publisher_domain",
				Type:        proto.ColumnType_STRING,
				Description: "The verified publisher domain for the application.",
				Transform:   graphField("publisherDomain"),
			},
			{
				Name:        "sign_in_audience",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the Microsoft accounts that are supported for the current application.",
				Transform:   graphField("signInAudience"),
			},

			// JSON fields
			{
				Name:        "api",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies settings for an application that implements a web API.",
				Transform:   graphField("api"),
			},
			{
				Name:        "identifier_uris",
				Type:        proto.ColumnType_JSON,
				Description: "The URIs that identify the application within its Azure AD tenant, or within a verified custom domain if the application is multi-tenant.",
				Transform:   graphField("identifierUris"),
			},
			{
				Name:        "info",
				Type:        proto.ColumnType_JSON,
				Description: "Basic profile information of the application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience.",
				Transform:   graphField("info"),
			},
			{
				Name:        "key_credentials",
				Type:        proto.ColumnType_JSON,
				Description: "The collection of key credentials associated with the application.",
				Transform:   graphField("keyCredentials"),
			},
			{
				Name:        "owner_ids",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGraphOwners,
				Transform:   transform.FromValue(),
				Description: "Id of the owners of the application. The owners are a set of non-admin users who are allowed to modify this object.",
			},
			{
				Name:        "parental_control_settings",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies parental control settings for an application.",
				Transform:   graphField("parentalControlSettings"),
			},
			{
				Name:        "password_credentials",
				Type:        proto.ColumnType_JSON,
				Description: "The collection of password credentials associated with the application.",
				Transform:   graphField("passwordCredentials"),
			},
			{
				Name:        "spa",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.",
				Transform:   graphField("spa"),
			},
			{
				Name:        "tags_src",
				Type:        proto.ColumnType_JSON,
				Description: "Custom strings that can be used to categorize and identify the application.",
				Transform:   graphField("tags"),
			},
			{
				Name:        "web",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies settings for a web application.",
				Transform:   graphField("web"),
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
