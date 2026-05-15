---
page_title: "anypoint_private_space_advanced_config Resource - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Manages advanced configuration of a Private Space — ingress settings and IAM role.
---

# anypoint_private_space_advanced_config (Resource)

Manages **advanced configuration** of a Private Space:
- **Ingress configuration**: read response timeout, protocol, port log level, per-IP log filters
- **IAM role**: enable/disable AWS IAM role for the Private Space (required for some integrations like AWS S3 access from inside a Mule app)

This resource is typically used **after** `anypoint_private_space` is created and the space is `ACTIVE`. It updates settings via a `PATCH` to the same endpoint.

> **One per Private Space.** Use a single `anypoint_private_space_advanced_config` per Private Space.

## Example Usage

### Minimal — enable IAM role, defaults for ingress

```terraform
resource "anypoint_private_space_advanced_config" "main" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id

  enable_iam_role = true
}
```

### Full ingress configuration with log filters

```terraform
resource "anypoint_private_space_advanced_config" "main" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id

  enable_iam_role = true

  ingress_configuration {
    read_response_timeout = "600"            # seconds
    protocol              = "https-redirect" # or "https"
    port_log_level        = "INFO"           # "ERROR" | "INFO" | "DEBUG"

    log_filters {
      ip    = "203.0.113.42"
      level = "DEBUG"
    }

    log_filters {
      ip    = "203.0.113.43"
      level = "INFO"
    }
  }
}
```

## Schema

### Required

- `org_id` (String, ForceNew) The organization id where the private space is defined.
- `private_space_id` (String, ForceNew) The unique identifier of the private space.

### Optional

- `enable_iam_role` (Boolean, Default `false`) Whether to enable IAM role for the private space.
- `ingress_configuration` (Block List, Max 1) Ingress configuration. See [below](#nested-schema-for-ingress_configuration).

### Read-Only

- `id` (String) Same as `private_space_id`.

### Nested Schema for `ingress_configuration`

Optional:

- `read_response_timeout` (String, Default `"300"`) Read response timeout in seconds.
- `protocol` (String, Default `"https-redirect"`) Ingress protocol. Valid values: `"https-redirect"`, `"https"`.
- `port_log_level` (String, Default `"ERROR"`) Default port log level. Valid values: `"ERROR"`, `"INFO"`, `"DEBUG"`.
- `log_filters` (Block List) Per-IP log level overrides. See [below](#nested-schema-for-log_filters).

#### Nested Schema for `log_filters`

Required:

- `ip` (String) IP address for the log filter.
- `level` (String) Log level for this IP filter. Valid values: `"ERROR"`, `"INFO"`, `"DEBUG"`.

## Import

```shell
terraform import anypoint_private_space_advanced_config.main <private_space_id>
```

> ⚠️ **On destroy** the resource resets the advanced config to defaults (`enable_iam_role = false`, empty ingress filters). It does NOT delete the Private Space itself.
