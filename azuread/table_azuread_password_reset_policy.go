package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getPasswordResetPolicyFromGraph(ctx context.Context, client *GraphClient) (graphObject, error) {
	policy, err := client.get(ctx, graphDefault, "policies/authorizationPolicy", nil)
	if err != nil {
		return nil, err
	}

	return graphObject{
		"dataSource":    "graph",
		"sourcePayload": graphObject{"authorizationPolicy": policy},
		// This applies to administrators, not the user SSPR enablement scope.
		"administratorsAllowedToUseSSPR": policy["allowedToUseSSPR"],
	}, nil
}

func tableAzureAdPasswordResetPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_password_reset_policy",
		Description: "Self-service password reset policy from the internal portal when configured, otherwise partial Microsoft Graph coverage.",
		List: &plugin.ListConfig{
			Hydrate: listGraphTable,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "Policy object identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   graphField("objectId"),
			},
			{
				Name:        "administrators_allowed_to_use_sspr",
				Description: "Graph administrator SSPR permission; not the user SSPR scope.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("administratorsAllowedToUseSSPR"),
			},
			{
				Name:        "enablement_type",
				Description: "Password reset enablement mode.",
				Type:        proto.ColumnType_INT,
				Transform:   graphField("enablementType"),
			},
			{
				Name:        "number_of_authentication_methods_required",
				Description: "Required authentication method count.",
				Type:        proto.ColumnType_INT,
				Transform:   graphField("numberOfAuthenticationMethodsRequired"),
			},
			{
				Name:        "registration_required_on_sign_in",
				Description: "Whether registration is required at sign-in.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("registrationRequiredOnSignIn"),
			},
			{
				Name:        "notify_users_on_password_reset",
				Description: "Whether users are notified after a password reset.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("notifyUsersOnPasswordReset"),
			},
			{
				Name:        "notify_on_admin_password_reset",
				Description: "Whether administrator password resets generate notifications.",
				Type:        proto.ColumnType_BOOL,
				Transform:   graphField("notifyOnAdminPasswordReset"),
			},
			{
				Name:        "password_reset_enabled_group_ids",
				Description: "Groups enabled for password reset.",
				Type:        proto.ColumnType_JSON,
				Transform:   graphField("passwordResetEnabledGroupIds"),
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
