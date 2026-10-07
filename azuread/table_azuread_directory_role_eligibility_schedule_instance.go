package azuread

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAzureAdDirectoryRoleEligibilityScheduleInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuread_directory_role_eligibility_schedule_instance",
		Description: "Represents the schedule instances for role eligibility operations on Azure AD resources.",
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
				Description: "The unique identifier for the role eligibility schedule instance.",
				Transform:   graphField("id"),
			},
			{
				Name:        "principal_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the principal that's in the scope of the role eligibility.",
				Transform:   graphField("principalId"),
			},
			{
				Name:        "role_definition_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the role definition that's in the scope of the role eligibility.",
				Transform:   graphField("roleDefinitionId"),
			},
			{
				Name:        "directory_scope_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the directory scope that's in the scope of the role eligibility.",
				Transform:   graphField("directoryScopeId"),
			},
			{
				Name:        "app_scope_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the app scope that's in the scope of the role eligibility.",
				Transform:   graphField("appScopeId"),
			},
			{
				Name:        "start_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The start date and time of the role eligibility schedule instance.",
				Transform:   graphField("startDateTime"),
			},
			{
				Name:        "end_date_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The end date and time of the role eligibility schedule instance.",
				Transform:   graphField("endDateTime"),
			},
			{
				Name:        "member_type",
				Type:        proto.ColumnType_STRING,
				Description: "How the role eligibility is inherited. It can be Inherited, Direct, or Group.",
				Transform:   graphField("memberType"),
			},
			{
				Name:        "role_eligibility_schedule_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier of the role eligibility schedule.",
				Transform:   graphField("roleEligibilityScheduleId"),
			},

			// JSON columns for complex objects
			{
				Name:        "app_scope",
				Type:        proto.ColumnType_JSON,
				Description: "The app scope that's in the scope of the role eligibility.",
				Transform:   transform.FromP(graphRoleReference, "appScope"),
			},
			{
				Name:        "directory_scope",
				Type:        proto.ColumnType_JSON,
				Description: "The directory scope that's in the scope of the role eligibility.",
				Transform:   transform.FromP(graphRoleReference, "directoryScope"),
			},
			{
				Name:        "principal",
				Type:        proto.ColumnType_JSON,
				Description: "The principal that's in the scope of the role eligibility.",
				Transform:   transform.FromP(graphRoleReference, "principal"),
			},
			{
				Name:        "role_definition",
				Type:        proto.ColumnType_JSON,
				Description: "The role definition that's in the scope of the role eligibility.",
				Transform:   transform.FromP(graphRoleReference, "roleDefinition"),
			},

			// Standard columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   graphTitle("id", "principalId"),
			},
		}),
	}
}
