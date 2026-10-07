package azuread

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func (c *GraphClient) graphPolicyFallback(ctx context.Context, table string) (graphObject, error) {
	row := graphObject{"dataSource": "graph"}

	if table == "azuread_password_policy" {
		settings, err := c.policySettings(ctx, graphBeta, "settings", "Password Rule Settings")
		if err != nil {
			return nil, err
		}

		row["sourcePayload"] = graphObject{"settings": settings}

		for _, setting := range settings {
			values := policySettingValues(setting)

			for from, to := range map[string]string{
				"LockoutThreshold":         "lockoutThreshold",
				"LockoutDurationInSeconds": "lockoutDurationInSeconds",
			} {
				if v, ok := values[from]; ok {
					n, err := strconv.Atoi(v)
					if err != nil {
						return nil, fmt.Errorf("invalid %s: %w", from, err)
					}

					row[to] = n
				}
			}

			for from, to := range map[string]string{
				"EnableBannedPasswordCheck":           "enforceCustomBannedPasswords",
				"EnableBannedPasswordCheckOnPremises": "enableBannedPasswordCheckOnPremises",
			} {
				if v, ok := values[from]; ok {
					b, err := strconv.ParseBool(v)
					if err != nil {
						return nil, fmt.Errorf("invalid %s: %w", from, err)
					}

					row[to] = b
				}
			}

			if v, ok := values["BannedPasswordList"]; ok {
				list := []string{}
				if v != "" {
					list = strings.Split(v, "\t")
				}

				row["customBannedPasswords"] = list
			}

			// Graph uses Audit/Enforce strings; the portal uses a numeric enum.
			if v, ok := values["BannedPasswordCheckOnPremisesMode"]; ok {
				row["graphOnPremisesPasswordCheckMode"] = v
			}
		}

		return row, nil
	}

	policy, err := c.get(ctx, graphDefault, "policies/authorizationPolicy", nil)
	if err != nil {
		return nil, err
	}

	payload := graphObject{"authorizationPolicy": policy}
	row["sourcePayload"] = payload
	permissions, _ := policy["defaultUserRolePermissions"].(map[string]any)

	switch table {
	case "azuread_directory_properties":
		row["usersCanRegisterApps"] = permissions["allowedToCreateApps"]
		row["allowInvitesFrom"] = policy["allowInvitesFrom"]

	case "azuread_password_reset_policy":
		// This applies to administrators, not the user SSPR enablement scope.
		row["administratorsAllowedToUseSSPR"] = policy["allowedToUseSSPR"]

	case "azuread_self_service_group_management":
		row["usersCanCreateSecurityGroups"] = permissions["allowedToCreateSecurityGroups"]

		settings, err := c.policySettings(ctx, graphDefault, "groupSettings", "Group.Unified")
		if err != nil {
			return nil, err
		}

		payload["groupSettings"] = settings

		for _, setting := range settings {
			values := policySettingValues(setting)

			if v, ok := values["EnableGroupCreation"]; ok {
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("invalid EnableGroupCreation: %w", err)
				}

				row["usersCanCreateMicrosoft365Groups"] = b
			}

			if v, ok := values["GroupCreationAllowedGroupId"]; ok {
				row["groupCreationAllowedGroupId"] = v
			}
		}

	default:
		return nil, fmt.Errorf("no Graph policy mapping for %s", table)
	}

	return row, nil
}

func (c *GraphClient) policySettings(ctx context.Context, endpoint graphEndpoint, path, displayName string) ([]graphObject, error) {
	settings := []graphObject{}
	err := c.list(ctx, endpoint, path, nil, func(row graphObject) bool {
		if row["displayName"] == displayName {
			settings = append(settings, row)
		}

		return true
	})

	return settings, err
}

func policySettingValues(setting graphObject) map[string]string {
	result := map[string]string{}
	values, _ := setting["values"].([]any)

	for _, item := range values {
		if value, ok := item.(map[string]any); ok {
			name, nameOK := value["name"].(string)
			text, textOK := value["value"].(string)

			if nameOK && textOK {
				result[name] = text
			}
		}
	}

	return result
}
