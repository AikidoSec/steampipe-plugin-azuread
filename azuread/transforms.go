package azuread

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func graphField(path string) *transform.ColumnTransforms {
	return transform.FromP(graphFieldValue, path)
}

func graphValue(row any, path string) any {
	current := row

	for _, key := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil
		}

		current = object[key]
	}

	return current
}

func graphFieldValue(_ context.Context, d *transform.TransformData) (any, error) {
	return graphValue(d.HydrateItem, d.Param.(string)), nil
}

func graphTitle(fields ...string) *transform.ColumnTransforms {
	return transform.FromP(graphTitleValue, fields)
}

func graphTitleValue(_ context.Context, d *transform.TransformData) (any, error) {
	for _, field := range d.Param.([]string) {
		if value := graphValue(d.HydrateItem, field); value != nil {
			return value, nil
		}
	}

	return nil, nil
}

func graphTags(_ context.Context, d *transform.TransformData) (any, error) {
	row := d.HydrateItem.(graphObject)

	if labels, ok := row["assignedLabels"].([]any); ok {
		tags := map[string]any{}
		for _, value := range labels {
			label, ok := value.(map[string]any)
			if !ok {
				continue
			}

			if id, ok := label["labelId"].(string); ok {
				tags[id] = label["displayName"]
			}
		}

		if len(tags) == 0 {
			return nil, nil
		}

		return tags, nil
	}

	values, ok := row["tags"].([]any)
	if !ok {
		return nil, nil
	}

	tags := map[string]bool{}
	for _, value := range values {
		if tag, ok := value.(string); ok {
			tags[tag] = true
		}
	}

	return tags, nil
}

func graphSignInActivity(_ context.Context, d *transform.TransformData) (any, error) {
	value, ok := graphValue(d.HydrateItem, "signInActivity").(map[string]any)
	if !ok {
		return nil, nil
	}

	return graphObject{
		"LastSignInDateTime":                value["lastSignInDateTime"],
		"LastSignInRequestId":               value["lastSignInRequestId"],
		"LastNonInteractiveSignInDateTime":  value["lastNonInteractiveSignInDateTime"],
		"LastNonInteractiveSignInRequestId": value["lastNonInteractiveSignInRequestId"],
	}, nil
}

func graphMemberships(_ context.Context, d *transform.TransformData) (any, error) {
	values, ok := graphValue(d.HydrateItem, "memberOf").([]any)
	if !ok {
		return nil, nil
	}

	members := []graphObject{}
	for _, value := range values {
		if row, ok := value.(map[string]any); ok {
			members = append(members, graphObject{
				"id":          row["id"],
				"@odata.type": row["@odata.type"],
			})
		}
	}

	return members, nil
}

func graphRoleReference(_ context.Context, d *transform.TransformData) (any, error) {
	row, ok := graphValue(d.HydrateItem, d.Param.(string)).(map[string]any)
	if !ok {
		return nil, nil
	}

	result := graphObject{}
	for _, key := range []string{"id", "@odata.type", "description"} {
		if value, ok := row[key]; ok {
			result[key] = value
		}
	}

	if value, ok := row["displayName"]; ok {
		result["display_name"] = value
	}

	return result, nil
}

func graphLocationType(_ context.Context, d *transform.TransformData) (any, error) {
	row := d.HydrateItem.(graphObject)

	switch row["@odata.type"] {
	case "#microsoft.graph.ipNamedLocation":
		return "IP", nil
	case "#microsoft.graph.countryNamedLocation":
		return "Country", nil
	default:
		return "Unknown", nil
	}
}

func graphLocationInfo(_ context.Context, d *transform.TransformData) (any, error) {
	row := d.HydrateItem.(graphObject)

	switch row["@odata.type"] {
	case "#microsoft.graph.countryNamedLocation":
		return graphObject{
			"Countries_and_Regions":             row["countriesAndRegions"],
			"Get_Unknown_Countries_and_Regions": row["includeUnknownCountriesAndRegions"],
			"Lookup_Method":                     row["countryLookupMethod"],
		}, nil

	case "#microsoft.graph.ipNamedLocation":
		result := graphObject{"IsTrusted": row["isTrusted"]}
		for _, kind := range []string{"IPv4Cidr", "IPv6Cidr", "IPv4Range", "IPv6Range"} {
			result[kind] = []graphObject{}
		}

		ranges, _ := row["ipRanges"].([]any)
		for _, value := range ranges {
			r, ok := value.(map[string]any)
			if !ok {
				continue
			}

			var kind string
			switch r["@odata.type"] {
			case "#microsoft.graph.iPv4CidrRange":
				kind = "IPv4Cidr"
			case "#microsoft.graph.iPv6CidrRange":
				kind = "IPv6Cidr"
			case "#microsoft.graph.iPv4Range":
				kind = "IPv4Range"
			case "#microsoft.graph.iPv6Range":
				kind = "IPv6Range"
			default:
				continue
			}

			data := graphObject{"Address": r["cidrAddress"]}
			if strings.HasSuffix(kind, "Range") {
				data = graphObject{
					"Lower": r["lowerAddress"],
					"Upper": r["upperAddress"],
				}
			}

			result[kind] = append(result[kind].([]graphObject), data)
		}

		return result, nil
	}

	return nil, nil
}
