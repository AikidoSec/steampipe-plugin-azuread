# Table: azuread_password_reset_policy

Password reset policy collected from internal portal `/PasswordReset/PasswordResetPolicies`. Returns one row per tenant, including a `raw` JSON column with the full response.

With `internal_api_refresh_token` (or `AZURE_INTERNAL_API_REFRESH_TOKEN`) configured, uses the internal portal and requires a tenant ID on Azure Public Cloud. Otherwise, queries Microsoft Graph with the normal connection credentials. Graph coverage is partial; unavailable portal fields are NULL. The `data_source` column identifies the selected source, and `raw` preserves its source payload. See the fallback field mapping in [HTTP collection and authentication](../http-client.md#graph-sources-when-no-portal-token-is-configured).

## Examples

```sql
select enablement_type, registration_required_on_sign_in, raw
from azuread_password_reset_policy;
```

See [HTTP collection and authentication](../http-client.md) for configuration.
