package azuread

import "context"

func policySettings(ctx context.Context, client *GraphClient, endpoint graphEndpoint, path, displayName string) ([]graphObject, error) {
	settings := []graphObject{}
	err := client.list(ctx, endpoint, path, nil, func(row graphObject) bool {
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
