# Table: azuread_device_registration_policy

Device registration policy collected from Microsoft Graph beta `/policies/deviceRegistrationPolicy`. Returns one row per tenant, including a `raw` JSON column with the full response.

Uses beta automatically, regardless of the connection's `graph_api_version`.

## Examples

```sql
select user_device_quota, azure_ad_join, raw
from azuread_device_registration_policy;
```

See [HTTP collection and authentication](../http-client.md) for configuration.
