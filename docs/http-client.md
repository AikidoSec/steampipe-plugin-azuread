# Graph HTTP collection

The plugin queries Microsoft Graph using client-secret, certificate, managed-identity, or Azure CLI credentials. Internal portal queries require a separate delegated refresh token.

## API routing

Routing is selected per resource, without requiring a version setting:

| Resources | API |
| --- | --- |
| Users, groups, applications, service principals, devices, roles, assignments, registration reports, security defaults, authorization policy | Microsoft Graph v1.0 |
| Conditional access policies | `/beta/identity/conditionalAccess/policies` |
| Named locations | `/beta/identity/conditionalAccess/namedLocations` |
| Directory settings | `/beta/settings` |
| Device registration policy | `/beta/policies/deviceRegistrationPolicy` |
| Group members | `/beta/groups/{id}/members` (documented workaround for omitted service principals in v1.0) |
| Internal portal resources | Separate portal API (see below) |

Audit logs, domains, and PIM eligibility use v1.0. Directory settings return one SQL row per setting value.

The optional `graph_api_version = "beta"` connection setting (or `AZURE_GRAPH_API_VERSION=beta`) switches otherwise-v1.0 resources to beta. Resources explicitly routed to beta above always use beta, even when the connection setting is `v1.0`. There is no automatic fallback on a permission error.

Azure China uses `https://microsoftgraph.chinacloudapi.cn`. Azure US Government uses `https://graph.microsoft.us`, with the corresponding authentication authority and token scope.

## Internal portal resources

These tables call `https://main.iam.ad.ext.azure.com/api` only when queried:

| Table | Path |
| --- | --- |
| `azuread_password_reset_policy` | `/PasswordReset/PasswordResetPolicies` |
| `azuread_password_policy` | `/AuthenticationMethods/PasswordPolicy` |
| `azuread_directory_properties` | `/Directories/Properties` |
| `azuread_self_service_group_management` | `/Directories/SsgmProperties` |

To select the internal portal source, configure `tenant_id` / `AZURE_TENANT_ID` and a delegated refresh token supplied as `internal_api_refresh_token` / `AZURE_INTERNAL_API_REFRESH_TOKEN`. This uses the Azure PowerShell client ID and OAuth refresh-token grant. The refresh token must be valid for that grant and the portal audience; an ordinary Microsoft Graph access token or client secret alone is insufficient. Access tokens are cached until shortly before expiration; rotated refresh tokens are retained in memory for the lifetime of the cached client, not written to configuration.

Internal portal collection is optional and only configured for Azure Public Cloud. Without a configured refresh token, these tables query Microsoft Graph using the normal connection credentials. Graph coverage is partial: unsupported fields remain SQL NULL. With a token configured, HTTP/authentication errors are surfaced rather than represented as empty or disabled policies. These are internal, unversioned APIs: endpoint availability and delegated permissions must be verified against the tenant. Each policy table includes a `raw` JSON column. For portal collection it contains the original response; for Graph collection it contains the source objects grouped by endpoint. The four dual-source tables expose `data_source` (`graph` or `internal_portal`).

```hcl
connection "azuread" {
  plugin = "azuread"
  tenant_id = "YOUR_TENANT_ID"
  # Existing Graph credentials may be supplied through AZURE_CLIENT_ID and AZURE_CLIENT_SECRET.
  # Supply AZURE_INTERNAL_API_REFRESH_TOKEN separately for portal tables.
}
```

```sql
select user_device_quota, azure_ad_join from azuread_device_registration_policy;
select lockout_threshold, custom_banned_passwords from azuread_password_policy;
select enablement_type, raw from azuread_password_reset_policy;
select restrict_non_admin_users, raw from azuread_directory_properties;
select self_service_group_management_enabled, raw from azuread_self_service_group_management;
```

## Query behavior

Memberships, owners, registered users, and registered devices require separate API calls when their columns are selected. User `manager` returns the manager's object ID, or NULL when no manager exists.

Nested JSON columns can contain additional properties returned by the API. Enum values are represented as API strings.

The client streams collections page by page, follows complete `@odata.nextLink` URLs, retains request headers across pages, and stops when the query has enough rows. Malformed or repeated pagination links return errors. Every request includes Graph's documented `client-request-id` plus the `x-ms-client-request-id` and `x-ms-correlation-id` headers. API errors preserve the client request ID and the server's `request-id` header for diagnosis. User queries selecting or filtering `signInActivity` use a maximum page size of 500. On the specific `Authentication_RequestFromNonPremiumTenantOrB2CTenant` licensing error, the initial user request is retried without selecting `signInActivity`, which remains NULL. This fallback is not applied to continuation URLs or queries filtering on sign-in activity. Other licensing and permission errors remain visible. There is no legacy Azure AD Graph fallback.

Users, groups, devices, applications, and service principals explicitly select requested fields. Application and service-principal key credentials therefore include public key material when permitted. User registered devices are collected through paginated relationship requests. The legacy application column `is_authorization_service_enabled` remains available but returns NULL because Microsoft Graph does not expose that property. The `oauth2_require_post_response` column reads Graph's `oauth2RequiredPostResponse` property, and device `extension_attributes` reads `extensionAttributes`.

Requests have a 60-second HTTP timeout and retry HTTP 429/500/502/503/504 up to three times. Retry delays follow `Retry-After` (seconds or HTTP date), otherwise use exponential backoff. Cancellation stops requests and retry waits. Cross-origin next links and redirects are rejected.

### API limits and permissions

Service-principal pages contain up to 100 objects. User pages contain up to 999 objects, or 500 when selecting or filtering `signInActivity`. Collections follow server-provided pagination links. [Directory-role membership](https://learn.microsoft.com/en-us/graph/api/directoryrole-list-members?view=graph-rest-1.0) returns a default of 1,000 objects and does not support `$top` pagination.

Group membership uses beta to include service principals omitted by the [v1.0 endpoint](https://learn.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-1.0). User app-role queries include indirect assignments through direct group memberships.

US Government L4 is supported; DoD L5 is not. Internal portal endpoints are private, unversioned APIs outside the public Graph API contract.

Device-registration policy requires `Policy.Read.DeviceConfiguration` for application access. See [device registration permissions](https://learn.microsoft.com/en-us/graph/api/deviceregistrationpolicy-get?view=graph-rest-beta).

### Graph sources when no portal token is configured

| Table | Graph endpoint | Exposed fallback fields |
| --- | --- | --- |
| `azuread_password_policy` | `/beta/settings`, `Password Rule Settings` | Lockout threshold/duration, custom banned passwords and enforcement, on-premises check enabled; string mode in `graph_on_premises_password_check_mode`. |
| `azuread_password_reset_policy` | `/v1.0/policies/authorizationPolicy` | `administrators_allowed_to_use_sspr` only; user SSPR scope, groups, registration and notification fields remain NULL. |
| `azuread_directory_properties` | `/v1.0/policies/authorizationPolicy` | `users_can_register_apps` and role-scoped `allow_invites_from`; portal IDs, display name, invitation boolean and access restrictions remain NULL. |
| `azuread_self_service_group_management` | `/v1.0/policies/authorizationPolicy` and `/v1.0/groupSettings`, `Group.Unified` | `users_can_create_security_groups`, `users_can_create_microsoft365_groups`, `group_creation_allowed_group_id`. Portal management fields remain NULL because creation and management differ. |

The connection's beta override also applies to the otherwise-v1.0 endpoints. Collections are paginated. Missing settings objects or values are left unknown, not filled with assumed defaults; the tenant row still exposes the source payload and `data_source`. No portal numeric enum is inferred from Graph's Audit/Enforce strings. Source selection happens before requests: a configured but invalid portal token produces an error, not a silent switch to Graph.

Graph application credentials require consented permissions for the selected endpoints (for example `Policy.Read.All` for authorization policy and permissions to read directory/group settings). See [authorization policy](https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-get?view=graph-rest-1.0), [group settings](https://learn.microsoft.com/en-us/graph/group-directory-settings), and [Microsoft365DSC's password settings implementation](https://www.powershellgallery.com/packages/Microsoft365DSC/1.25.723.2/Content/DSCResources%5CMSFT_AADPasswordRuleSettings%5CMSFT_AADPasswordRuleSettings.psm1).
