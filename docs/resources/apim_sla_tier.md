---
page_title: "anypoint_apim_sla_tier Resource - terraform-provider-anypoint"
subcategory: "API Manager"
description: |-
  Manages an SLA tier for an API Manager instance.
---

# anypoint_apim_sla_tier (Resource)

Manages a **SLA tier** for an API Manager instance. An SLA tier defines rate limits (max requests per time window) that get applied to client applications subscribed to that tier.

Typical use: combine with the `anypoint_apim_policy_rate_limiting` policy to enforce tier-based throttling.

## Example Usage

### Basic — Bronze tier with 100 req/min

```terraform
resource "anypoint_apim_sla_tier" "bronze" {
  org_id          = var.org_id
  env_id          = var.env_id
  api_instance_id = anypoint_apim_mule4.my_api.id

  name         = "Bronze"
  description  = "Basic tier — 100 requests per minute"
  auto_approve = true
  status       = "ACTIVE"

  limits {
    time_period_in_milliseconds = 60000   # 1 minute
    maximum_requests            = 100
    visible                     = true
  }
}
```

### Multi-window — Silver tier with two limits

```terraform
resource "anypoint_apim_sla_tier" "silver" {
  org_id          = var.org_id
  env_id          = var.env_id
  api_instance_id = anypoint_apim_mule4.my_api.id

  name         = "Silver"
  description  = "1000 req/min + 50000 req/day"
  auto_approve = false   # requires manual approval
  status       = "ACTIVE"

  limits {
    time_period_in_milliseconds = 60000      # 1 minute
    maximum_requests            = 1000
  }

  limits {
    time_period_in_milliseconds = 86400000   # 1 day
    maximum_requests            = 50000
  }
}
```

## Schema

### Required

- `org_id` (String, ForceNew) Organization id.
- `env_id` (String, ForceNew) Environment id where the API instance lives.
- `api_instance_id` (String, ForceNew) Numeric API instance id (the one from API Manager, NOT the asset/exchange id).
- `name` (String) Tier name (e.g. `"Bronze"`, `"Silver"`, `"Gold"`).
- `limits` (Block List) Rate-limit windows. See [below](#nested-schema-for-limits).

### Optional

- `description` (String) Free-form description shown in the dev portal.
- `auto_approve` (Boolean, Default `false`) If `true`, client subscriptions to this tier are auto-approved without manual review.
- `status` (String, Default `"ACTIVE"`) Tier status. One of `"ACTIVE"`, `"INACTIVE"`.

### Read-Only

- `id` (String) Numeric tier id.

### Nested Schema for `limits`

Required:

- `time_period_in_milliseconds` (Int) Time window. Common values:
  - `1000` — 1 second
  - `60000` — 1 minute
  - `3600000` — 1 hour
  - `86400000` — 1 day
- `maximum_requests` (Int) Maximum requests allowed within the time window.

Optional:

- `visible` (Boolean, Default `true`) Whether this limit is shown publicly in the dev portal.

## Import

Import key uses the composite `<ORG_ID>/<ENV_ID>/<API_INSTANCE_ID>/<TIER_ID>`:

```shell
terraform import anypoint_apim_sla_tier.bronze \
  47ec5a2c-c3ce-4994-af2e-11ed525c5b78/a066601c-e7b9-4b8b-95bc-d8a2ae5a91c2/12345/678
```
