package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_group",
		Description: "Represents an Azure AD group.",
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
				ShouldIgnoreErrorFunc: isIgnorableErrorPredicate([]string{"Invalid filter clause"}),
			},
			KeyColumns: plugin.KeyColumnSlice{
				// Key fields
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "filter",
					Require: plugin.Optional,
				},
				{
					Name:    "mail",
					Require: plugin.Optional,
				},
				{
					Name:      "mail_enabled",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
				{
					Name:      "on_premises_sync_enabled",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
				{
					Name:      "security_enabled",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
			},
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "preferred_data_location",
				Type:        proto.ColumnType_STRING,
				Description: "The preferred location for the group's data.",
				Transform:   graphField("preferredDataLocation"),
			},
			{
				Name:        "deleted_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time the group was deleted, if applicable.",
				Transform:   graphField("deletedDateTime"),
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
				Description: "The unique identifier for the group.",
				Transform:   graphField("id"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "An optional description for the group.",
				Transform:   graphField("description"),
			},
			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
				Description: "Odata query to search for groups.",
			},

			// Other fields
			{
				Name:        "classification",
				Type:        proto.ColumnType_STRING,
				Description: "Describes a classification for the group (such as low, medium or high business impact).",
				Transform:   graphField("classification"),
			},
			{
				Name:        "created_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time at which the group was created.",
				Transform:   graphField("createdDateTime"),
			},
			{
				Name:        "expiration_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Timestamp of when the group is set to expire.",
				Transform:   graphField("expirationDateTime"),
			},
			{
				Name:        "is_assignable_to_role",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether this group can be assigned to an Azure Active Directory role or not.",
				Transform:   graphField("isAssignableToRole"),
			},
			{
				Name:        "is_subscribed_by_mail",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the signed-in user is subscribed to receive email conversations. Default value is true.",
				Hydrate:     getGraphGroupSubscription,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "mail",
				Type:        proto.ColumnType_STRING,
				Description: "The SMTP address for the group, for example, \"serviceadmins@contoso.onmicrosoft.com\".",
				Transform:   graphField("mail"),
			},
			{
				Name:        "mail_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Specifies whether the group is mail-enabled.",
				Transform:   graphField("mailEnabled"),
			},
			{
				Name:        "mail_nickname",
				Type:        proto.ColumnType_STRING,
				Description: "The mail alias for the user.",
				Transform:   graphField("mailNickname"),
			},
			{
				Name:        "membership_rule",
				Type:        proto.ColumnType_STRING,
				Description: "The mail alias for the group, unique in the organization.",
				Transform:   graphField("membershipRule"),
			},
			{
				Name:        "membership_rule_processing_state",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates whether the dynamic membership processing is on or paused. Possible values are On or Paused.",
				Transform:   graphField("membershipRuleProcessingState"),
			},
			{
				Name:        "on_premises_domain_name",
				Type:        proto.ColumnType_STRING,
				Description: "Contains the on-premises Domain name synchronized from the on-premises directory.",
				Transform:   graphField("onPremisesDomainName"),
			},
			{
				Name:        "on_premises_last_sync_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Indicates the last time at which the group was synced with the on-premises directory.",
				Transform:   graphField("onPremisesLastSyncDateTime"),
			},
			{
				Name:        "on_premises_net_bios_name",
				Type:        proto.ColumnType_STRING,
				Description: "Contains the on-premises NetBiosName synchronized from the on-premises directory.",
				Transform:   graphField("onPremisesNetBiosName"),
			},
			{
				Name:        "on_premises_sam_account_name",
				Type:        proto.ColumnType_STRING,
				Description: "Contains the on-premises SAM account name synchronized from the on-premises directory.",
				Transform:   graphField("onPremisesSamAccountName"),
			},
			{
				Name:        "on_premises_security_identifier",
				Type:        proto.ColumnType_STRING,
				Description: "Contains the on-premises security identifier (SID) for the group that was synchronized from on-premises to the cloud.",
				Transform:   graphField("onPremisesSecurityIdentifier"),
			},
			{
				Name:        "on_premises_sync_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "True if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default).",
				Transform:   graphField("onPremisesSyncEnabled"),
			},
			{
				Name:        "renewed_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Timestamp of when the group was last renewed. This cannot be modified directly and is only updated via the renew service action.",
				Transform:   graphField("renewedDateTime"),
			},
			{
				Name:        "security_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Specifies whether the group is a security group.",
				Transform:   graphField("securityEnabled"),
			},
			{
				Name:        "security_identifier",
				Type:        proto.ColumnType_STRING,
				Description: "Security identifier of the group, used in Windows scenarios.",
				Transform:   graphField("securityIdentifier"),
			},
			{
				Name:        "visibility",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or Hiddenmembership.",
				Transform:   graphField("visibility"),
			},

			// JSON fields
			{
				Name:        "assigned_labels",
				Type:        proto.ColumnType_JSON,
				Description: "The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group.",
				Transform:   graphField("assignedLabels"),
			},
			{
				Name:        "group_types",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies the group type and its membership. If the collection contains Unified, the group is a Microsoft 365 group; otherwise, it's either a security group or distribution group. For details, see [groups overview](https://docs.microsoft.com/en-us/graph/api/resources/groups-overview?view=graph-rest-1.0).",
				Transform:   graphField("groupTypes"),
			},
			{
				Name:        "member_ids",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGraphMembers,
				Transform:   transform.FromValue(),
				Description: "Id of Users and groups that are members of this group.",
			},
			{
				Name:        "owner_ids",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getGraphOwners,
				Transform:   transform.FromValue(),
				Description: "Id od the owners of the group. The owners are a set of non-admin users who are allowed to modify this object.",
			},
			{
				Name:        "proxy_addresses",
				Type:        proto.ColumnType_JSON,
				Description: "Email addresses for the group that direct to the same group mailbox. For example: [\"SMTP: bob@contoso.com\", \"smtp: bob@sales.contoso.com\"]. The any operator is required to filter expressions on multi-valued properties.",
				Transform:   graphField("proxyAddresses"),
			},
			{
				Name:        "resource_behavior_options",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies the group behaviors that can be set for a Microsoft 365 group during creation. Possible values are AllowOnlyMembersToPost, HideGroupInOutlook, SubscribeNewGroupMembers, WelcomeEmailDisabled.",
				Transform:   graphField("resourceBehaviorOptions"),
			},
			{
				Name:        "resource_provisioning_options",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies the group resources that are provisioned as part of Microsoft 365 group creation, that are not normally part of default group creation. Possible value is Team.",
				Transform:   graphField("resourceProvisioningOptions"),
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

func getGraphMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	return graphRelationship(ctx, d, h, "members", true)
}

func getGraphOwners(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	return graphRelationship(ctx, d, h, "owners", true)
}
