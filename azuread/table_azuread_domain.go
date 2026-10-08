package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_domain",
		Description: "Represents an Azure Active Directory (Azure AD) domain.",
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
				Description: "The fully qualified name of the domain.",
				Transform:   graphField("id"),
			},
			{
				Name:        "authentication_type",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the configured authentication type for the domain. The value is either Managed or Federated. Managed indicates a cloud managed domain where Azure AD performs user authentication. Federated indicates authentication is federated with an identity provider such as the tenant's on-premises Active Directory via Active Directory Federation Services.",
				Transform:   graphField("authenticationType"),
			},

			// Other fields
			{
				Name:        "is_default",
				Type:        proto.ColumnType_BOOL,
				Description: "true if this is the default domain that is used for user creation. There is only one default domain per company.",
				Transform:   graphField("isDefault"),
			},
			{
				Name:        "is_admin_managed",
				Type:        proto.ColumnType_BOOL,
				Description: "The value of the property is false if the DNS record management of the domain has been delegated to Microsoft 365. Otherwise, the value is true.",
				Transform:   graphField("isAdminManaged"),
			},
			{
				Name:        "is_initial",
				Type:        proto.ColumnType_BOOL,
				Description: "true if this is the initial domain created by Microsoft Online Services (companyname.onmicrosoft.com). There is only one initial domain per company.",
				Transform:   graphField("isInitial"),
			},
			{
				Name:        "is_root",
				Type:        proto.ColumnType_BOOL,
				Description: "true if the domain is a verified root domain. Otherwise, false if the domain is a subdomain or unverified.",
				Transform:   graphField("isRoot"),
			},
			{
				Name:        "is_verified",
				Type:        proto.ColumnType_BOOL,
				Description: "true if the domain has completed domain ownership verification.",
				Transform:   graphField("isVerified"),
			},

			// Json fields
			{
				Name:        "supported_services",
				Type:        proto.ColumnType_JSON,
				Description: "The capabilities assigned to the domain. Can include 0, 1 or more of following values: Email, Sharepoint, EmailInternalRelayOnly, OfficeCommunicationsOnline, SharePointDefaultDomain, FullRedelegation, SharePointPublic, OrgIdAuthentication, Yammer, Intune. The values which you can add/remove using Graph API include: Email, OfficeCommunicationsOnline, Yammer.",
				Transform:   graphField("supportedServices"),
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
