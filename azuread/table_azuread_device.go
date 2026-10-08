package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdDevice(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_device",
		Description: "Represents an Azure AD device.",
		Get: &plugin.GetConfig{
			Hydrate:    getGraphTable,
			KeyColumns: plugin.SingleColumn("id"),
		},
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
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
					Name:    "operating_system",
					Require: plugin.Optional,
				},
				{
					Name:    "operating_system_version",
					Require: plugin.Optional,
				},
				{
					Name:    "profile_type",
					Require: plugin.Optional,
				},
				{
					Name:    "trust_type",
					Require: plugin.Optional,
				},
				{
					Name:    "filter",
					Require: plugin.Optional,
				},
			},
		},

		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "compliance_expiration_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the device's compliance expires.",
				Transform:   graphField("complianceExpirationDateTime"),
			},
			{
				Name:        "device_category",
				Type:        proto.ColumnType_STRING,
				Description: "The device category.",
				Transform:   graphField("deviceCategory"),
			},
			{
				Name:        "device_ownership",
				Type:        proto.ColumnType_STRING,
				Description: "The device ownership type.",
				Transform:   graphField("deviceOwnership"),
			},
			{
				Name:        "enrollment_profile_name",
				Type:        proto.ColumnType_STRING,
				Description: "The enrollment profile name.",
				Transform:   graphField("enrollmentProfileName"),
			},
			{
				Name:        "on_premises_last_sync_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The last synchronization from the on-premises directory.",
				Transform:   graphField("onPremisesLastSyncDateTime"),
			},
			{
				Name:        "on_premises_sync_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the device is synchronized from an on-premises directory.",
				Transform:   graphField("onPremisesSyncEnabled"),
			},
			{
				Name:        "registration_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "When the device was registered.",
				Transform:   graphField("registrationDateTime"),
			},
			{
				Name:        "system_labels",
				Type:        proto.ColumnType_JSON,
				Description: "Labels applied to the device by the system.",
				Transform:   graphField("systemLabels"),
			},

			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the device. Inherited from directoryObject.",
				Transform:   graphField("id"),
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name displayed for the device.",
				Transform:   graphField("displayName"),
			},
			{
				Name:        "account_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "True if the account is enabled; otherwise, false.",
				Transform:   graphField("accountEnabled"),
			},
			{
				Name:        "device_id",
				Type:        proto.ColumnType_STRING,
				Description: "Unique identifier set by Azure Device Registration Service at the time of registration.",
				Transform:   graphField("deviceId"),
			},
			{
				Name:        "approximate_last_sign_in_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The timestamp type represents date and time information using ISO 8601 format and is always in UTC time.",
				Transform:   graphField("approximateLastSignInDateTime"),
			},

			// Other fields
			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
				Description: "Odata query to search for resources.",
			},
			{
				Name:        "is_compliant",
				Type:        proto.ColumnType_BOOL,
				Description: "True if the device is compliant; otherwise, false.",
				Transform:   graphField("isCompliant"),
			},
			{
				Name:        "is_managed",
				Type:        proto.ColumnType_BOOL,
				Description: "True if the device is managed; otherwise, false.",
				Transform:   graphField("isManaged"),
			},
			{
				Name:        "mdm_app_id",
				Type:        proto.ColumnType_STRING,
				Description: "Application identifier used to register device into MDM.",
				Transform:   graphField("mdmAppId"),
			},
			{
				Name:        "operating_system",
				Type:        proto.ColumnType_STRING,
				Description: "The type of operating system on the device.",
				Transform:   graphField("operatingSystem"),
			},
			{
				Name:        "operating_system_version",
				Type:        proto.ColumnType_STRING,
				Description: "The version of the operating system on the device.",
				Transform:   graphField("operatingSystemVersion"),
			},
			{
				Name:        "profile_type",
				Type:        proto.ColumnType_STRING,
				Description: "A string value that can be used to classify device types.",
				Transform:   graphField("profileType"),
			},
			{
				Name:        "trust_type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of trust for the joined device. Possible values: Workplace (indicates bring your own personal devices), AzureAd (Cloud only joined devices), ServerAd (on-premises domain joined devices joined to Azure AD).",
				Transform:   graphField("trustType"),
			},

			// JSON fields
			{
				Name:        "extension_attributes",
				Type:        proto.ColumnType_JSON,
				Description: "Contains extension attributes 1-15 for the device. The individual extension attributes are not selectable. These properties are mastered in cloud and can be set during creation or update of a device object in Azure AD.",
				Transform:   graphField("extensionAttributes"),
			},
			{
				Name:        "member_of",
				Type:        proto.ColumnType_JSON,
				Description: "A list the groups and directory roles that the device is a direct member of.",
				Hydrate:     getGraphMemberOf,
				Transform:   transform.FromValue(),
			},

			{
				Name:        "registered_user_ids",
				Type:        proto.ColumnType_JSON,
				Description: "IDs of users registered to the device.",
				Hydrate:     getGraphRegisteredUsers,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphTitle("displayName", "deviceId"),
			},
		}),
	}
}

func getGraphRegisteredUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	return graphRelationship(ctx, d, h, "registeredUsers", true)
}
