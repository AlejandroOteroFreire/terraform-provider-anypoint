---
page_title: "anypoint_secretgroup_sharedsecret Resource - terraform-provider-anypoint"
subcategory: "Secrets Manager"
description: |-
  Manages a Shared Secret inside a Secret Group.
---

# anypoint_secretgroup_sharedsecret (Resource)

Manages a **Shared Secret** inside a Secret Group. Shared Secrets are credentials consumable by Mule applications (and TLS contexts, keystores, etc.) — they live alongside keystores, truststores and certificates in the same Secret Group.

Four types are supported:

| `type` | Use case | Required fields |
|--------|----------|-----------------|
| `UsernamePassword` | DB/HTTP basic-auth credentials | `username`, `password` |
| `S3Credential`     | AWS S3 access keys | `access_key_id`, `secret_access_key` |
| `SymmetricKey`     | AES-style symmetric encryption key | `key` (base64-encoded) |
| `Blob`             | Arbitrary opaque secret | `content` |

The `type` is immutable (`ForceNew`) — to change it, destroy and recreate.

## Example Usage

### Username/Password

```terraform
resource "anypoint_secretgroup_sharedsecret" "db_creds" {
  org_id          = var.org_id
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.app_secrets.id

  name = "db-creds"
  type = "UsernamePassword"

  username = "appuser"
  password = var.db_password   # mark sensitive in your variable
}
```

### AWS S3 credentials

```terraform
resource "anypoint_secretgroup_sharedsecret" "s3" {
  org_id          = var.org_id
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.app_secrets.id

  name              = "s3-uploader"
  type              = "S3Credential"
  access_key_id     = var.aws_access_key_id
  secret_access_key = var.aws_secret_access_key
  expiration_date   = "2026-12-31T23:59:59Z"
}
```

### Symmetric key (base64-encoded)

```terraform
resource "anypoint_secretgroup_sharedsecret" "aes" {
  org_id          = var.org_id
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.app_secrets.id

  name = "payload-aes-key"
  type = "SymmetricKey"
  key  = filebase64("${path.module}/keys/aes.key")
}
```

### Opaque blob

```terraform
resource "anypoint_secretgroup_sharedsecret" "license" {
  org_id          = var.org_id
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.app_secrets.id

  name    = "vendor-license"
  type    = "Blob"
  content = file("${path.module}/license.txt")
}
```

## Schema

### Required

- `org_id` (String, ForceNew) Organization id.
- `env_id` (String, ForceNew) Environment id.
- `secret_group_id` (String, ForceNew) Parent Secret Group id.
- `name` (String) Shared secret name (unique within the secret group).
- `type` (String, ForceNew) One of `UsernamePassword`, `S3Credential`, `SymmetricKey`, `Blob`.

### Optional

- `expiration_date` (String, Computed) ISO-8601 expiration date. Computed from the API if not set.
- `username` (String) Only when `type = "UsernamePassword"`.
- `password` (String, Sensitive) Only when `type = "UsernamePassword"`.
- `access_key_id` (String) Only when `type = "S3Credential"`.
- `secret_access_key` (String, Sensitive) Only when `type = "S3Credential"`.
- `key` (String, Sensitive) Base64-encoded symmetric key. Only when `type = "SymmetricKey"`.
- `content` (String, Sensitive) Opaque content. Only when `type = "Blob"`.

### Read-Only

- `id` (String) Shared secret id (UUID).

## Import

Import key is composite `<ORG_ID>/<ENV_ID>/<SG_ID>/<SECRET_ID>`:

```shell
terraform import anypoint_secretgroup_sharedsecret.db_creds \
  47ec5a2c.../a066601c.../39731075.../b2096f24...
```

> **Sensitive fields are not returned by the API.** After import, `password`, `secret_access_key`, `key`, and `content` will be empty in state until you set them via Terraform (which will trigger an update).
