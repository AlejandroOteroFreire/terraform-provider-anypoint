---
page_title: "anypoint_connected_app_scopes Resource - terraform-provider-anypoint"
subcategory: "Access Management"
description: |-
  Manages the scopes assigned to a Connected App.
---

# anypoint_connected_app_scopes (Resource)

Manages **scopes** assigned to a Connected App. Each scope can have **context parameters** (`org`, `env_id`) that limit where the scope applies.

A Connected App's scopes are not part of the `anypoint_connected_app` resource itself — they are a **separate** resource managed via `PATCH /connectedApplications/{id}/scopes`. This resource owns the full list: replacing this resource's `scopes` block replaces the entire scope set in Anypoint.

> ⚠️ **Requires the user-on-behalf provider** (`auth_type = "user"`). This Access Management operation is rejected by client-credentials tokens.
>
> See [`docs/index.md`](../index.md#user) for how to configure the `anypoint.admin` provider alias.

## Example Usage

### Org-level scopes (Datadog-style — read-only monitoring)

```terraform
resource "anypoint_connected_app_scopes" "datadog" {
  provider = anypoint.admin   # ← user-on-behalf provider required

  org_id           = var.org_id
  connected_app_id = anypoint_connected_app.datadog.id

  scopes {
    scope = "view:monitoring"
    context_params {
      org = var.org_id
    }
  }

  scopes {
    scope = "read:runtime_fabrics"
    context_params {
      org = var.org_id
    }
  }
}
```

### Mixed org + environment-scoped (CI/CD-style)

```terraform
resource "anypoint_connected_app_scopes" "cicd" {
  provider = anypoint.admin

  org_id           = var.org_id
  connected_app_id = anypoint_connected_app.cicd.id

  # Org-level: Exchange
  scopes {
    scope = "manage:exchange"
    context_params {
      org = var.org_id
    }
  }

  # Env-level: deploy on Sandbox
  scopes {
    scope = "create:applications"
    context_params {
      org    = var.org_id
      env_id = anypoint_env.sandbox.id
    }
  }

  scopes {
    scope = "delete:applications"
    context_params {
      org    = var.org_id
      env_id = anypoint_env.sandbox.id
    }
  }
}
```

## Schema

### Required

- `org_id` (String, ForceNew) The organization id where the connected app is defined.
- `connected_app_id` (String, ForceNew) The unique identifier of the connected app.
- `scopes` (Block List) List of scopes to assign. See [below](#nested-schema-for-scopes).

### Read-Only

- `id` (String) Same as `connected_app_id`.

### Nested Schema for `scopes`

Required:

- `scope` (String) The scope identifier (e.g. `manage:exchange`, `read:applications`). Use the canonical `verb:resource` form — **not** the display name. See the full catalog in your Anypoint UI under Access Management → Connected Apps → Add Scope.

Optional:

- `context_params` (Block List, Max 1) Context parameters scoping where the scope applies. Omit for scopes that don't take context (`profile`, `full`, `read:full`).

#### Nested Schema for `context_params`

Optional:

- `org` (String) The organization id where the scope applies (required for org- and env-level scopes).
- `env_id` (String) The environment id where the scope applies (only for env-level scopes like `create:applications`, `manage:apis`, etc.).

## Scope catalog

There are ~117 valid scopes split into three levels:

- **Org-scoped** (60): use `context_params { org = ... }`. Examples: `manage:exchange`, `admin:cloudhub`, `view:monitoring`.
- **Env-scoped** (55): use `context_params { org = ..., env_id = ... }`. Examples: `create:applications`, `manage:apis`, `read:secrets`.
- **No-context** (2): `profile`, `openid:google_wif`. Don't include `context_params`. Only valid for `auth_type = user` Connected Apps.

Refer to your Anypoint UI to discover the exact identifier for each scope.

## Import

```shell
terraform import anypoint_connected_app_scopes.example <CONNECTED_APP_ID>
```
