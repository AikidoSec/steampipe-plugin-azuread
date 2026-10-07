package azuread

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdDeviceRegistrationPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_device_registration_policy",
		Description: "Tenant device registration policy from Microsoft Graph beta.",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Policy identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("id"),
			},
			{
				Name:        "display_name",
				Description: "Policy display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("displayName"),
			},
			{
				Name:        "description",
				Description: "Policy description.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("description"),
			},
			{
				Name:        "user_device_quota",
				Description: "Maximum number of devices per user.",
				Type:        proto.ColumnType_INT,
				Transform:   graphField("userDeviceQuota"),
			},
			{
				Name:        "multi_factor_auth_configuration",
				Description: "MFA configuration for device registration.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("multiFactorAuthConfiguration"),
			},
			{
				Name:        "azure_ad_registration",
				Description: "Allowed device registration membership.",
				Type:        proto.ColumnType_JSON,
				Transform:   graphField("azureADRegistration"),
			},
			{
				Name:        "azure_ad_join",
				Description: "Allowed device join membership and local administrator configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   graphField("azureADJoin"),
			},
			{
				Name:        "local_admin_password",
				Description: "Local administrator password configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   graphField("localAdminPassword"),
			},
			{
				Name:        "raw",
				Type:        proto.ColumnType_JSON,
				Description: "Complete API response, including properties not exposed as individual columns.",
				Transform:   transform.FromValue(),
			},
		}),
	}
}
