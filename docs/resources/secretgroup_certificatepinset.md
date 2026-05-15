---
page_title: "anypoint_secretgroup_certificatepinset Resource - terraform-provider-anypoint"
subcategory: "Secrets Manager"
description: |-
  Manages a Certificate Pinset inside a Secret Group.
---

# anypoint_secretgroup_certificatepinset (Resource)

Manages a **Certificate Pinset** inside a Secret Group. A pinset is a bundle of trusted PEM certificates referenced by a TLS context for **certificate pinning** — Mule apps will only accept TLS connections whose server certificate matches one in the pinset.

The certificate content is uploaded as a base64-encoded PEM via `multipart/form-data`. The pinset is **immutable**: all fields are `ForceNew`; updating any value recreates the resource.

## Example Usage

```terraform
resource "anypoint_secretgroup_certificatepinset" "vendor_api_pins" {
  org_id          = var.org_id
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.app_secrets.id

  name = "vendor-api-pinset"

  # PEM file with one or more certificates concatenated
  certificate_pinset_base64 = filebase64("${path.module}/certs/vendor-pins.pem")
}

# Outputs returned by the API
output "pinset_algorithm" {
  value = anypoint_secretgroup_certificatepinset.vendor_api_pins.algorithm
}

output "pinset_expiration" {
  value = anypoint_secretgroup_certificatepinset.vendor_api_pins.expiration_date
}
```

### Use in a TLS context

```terraform
resource "anypoint_secretgroup_tlscontext_mule" "outbound" {
  org_id = var.org_id
  env_id = var.env_id
  sg_id  = anypoint_secretgroup.app_secrets.id
  name   = "vendor-outbound"

  # Reference the pinset by name
  certificate_pinset_path = anypoint_secretgroup_certificatepinset.vendor_api_pins.name

  # ... other TLS settings
}
```

## Schema

### Required

- `org_id` (String, ForceNew) Organization id.
- `env_id` (String, ForceNew) Environment id.
- `secret_group_id` (String, ForceNew) Parent Secret Group id.
- `name` (String, ForceNew) Pinset name (unique within the secret group).
- `certificate_pinset_base64` (String, ForceNew, Sensitive) Base64-encoded PEM content of the pinset. Use `filebase64()` if loading from disk.

### Read-Only

- `id` (String) Pinset id (UUID).
- `expiration_date` (String) ISO-8601 expiration date of the earliest certificate in the pinset.
- `algorithm` (String) Signature algorithm of the pinned certificate(s) (e.g. `SHA256withRSA`).

## Import

Import key is composite `<ORG_ID>/<ENV_ID>/<SG_ID>/<PINSET_ID>`:

```shell
terraform import anypoint_secretgroup_certificatepinset.vendor_api_pins \
  47ec5a2c.../a066601c.../39731075.../b2096f24...
```

> **Sensitive fields are not returned by the API.** After import, `certificate_pinset_base64` will be empty in state until you set it again. Since all fields are `ForceNew`, this will trigger a recreate on the next apply.
