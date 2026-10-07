package azuread

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableAzureAdPasswordPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_password_policy",
		Description: "Password protection policy from the internal portal when configured, otherwise partial Microsoft Graph coverage.",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "graph_on_premises_password_check_mode",
				Description: "Graph password protection mode (Audit or Enforce); portal numeric mode is not inferred.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("graphOnPremisesPasswordCheckMode"),
			},
			{
				Name:        "lockout_threshold",
				Description: "Failed attempts before account lockout.",
				Type:        proto.ColumnType_INT,
				Transform:   graphField("lockoutThreshold"),
			},
			{
				Name:        "lockout_duration_in_seconds",
				Description: "Account lockout duration.",
				Type:        proto.ColumnType_INT,
				Transform:   graphField("lockoutDurationInSeconds"),
			},
			{
				Name:        "enforce_custom_banned_passwords",
				Description: "Whether the custom banned password list is enforced.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("enforceCustomBannedPasswords"),
			},
			{
				Name:        "custom_banned_passwords",
				Description: "Custom banned password list.",
				Type:        proto.ColumnType_JSON,
				Transform:   graphField("customBannedPasswords"),
			},
			{
				Name:        "enable_banned_password_check_on_premises",
				Description: "Whether on-premises password protection is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("enableBannedPasswordCheckOnPremises"),
			},
			{
				Name:        "banned_password_check_on_premises_mode",
				Description: "On-premises password protection mode.",
				Type:        proto.ColumnType_INT,
				Transform:   graphField("bannedPasswordCheckOnPremisesMode"),
			},
			{
				Name:        "data_source",
				Description: "Source: graph (partial coverage) or internal_portal.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("dataSource"),
			},
			{
				Name:        "raw",
				Type:        proto.ColumnType_JSON,
				Description: "Complete API response, including properties not exposed as individual columns.",
				Transform:   graphField("sourcePayload"),
			},
		}),
	}
}
