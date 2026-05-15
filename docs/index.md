---
page_title: "Provider: Anypoint"
subcategory: "Mulesoft"
description: |-
  The Anypoint provider lets you manage MuleSoft Anypoint Platform resources using Terraform.
---

# Anypoint Provider

The `anypoint` provider lets you manage [MuleSoft Anypoint Platform](https://anypoint.mulesoft.com) resources using Terraform — Private Spaces, API instances and policies, environments and organizations, Connected Apps, Secrets Manager artifacts, Anypoint MQ destinations, Runtime Fabric clusters, and more.

This documentation covers the provider configuration. For per-resource reference, see the sidebar (or the [`docs/resources/`](https://github.com/mulesoft-anypoint/terraform-provider-anypoint/tree/master/docs/resources) and [`docs/data-sources/`](https://github.com/mulesoft-anypoint/terraform-provider-anypoint/tree/master/docs/data-sources) folders on GitHub).

## Example Usage

### Default — Connected App (client_credentials)

```terraform
terraform {
  required_providers {
    anypoint = {
      source  = "mulesoft-anypoint/anypoint"
      version = "~> 1.9"
    }
  }
}

provider "anypoint" {
  client_id     = var.anypoint_client_id      # or ANYPOINT_CLIENT_ID env var
  client_secret = var.anypoint_client_secret  # or ANYPOINT_CLIENT_SECRET env var
  cplane        = "us"                        # "us" (default), "eu", or "gov"
}
```

### Admin alias — user-on-behalf (Access Management)

Some resources (`anypoint_team_roles`, `anypoint_team_member`, `anypoint_connected_app_scopes`) require user-level permissions. Configure a second provider with `auth_type = "user"`:

```terraform
provider "anypoint" {
  alias         = "admin"
  auth_type     = "user"
  client_id     = var.anypoint_admin_client_id
  client_secret = var.anypoint_admin_client_secret
  username      = var.anypoint_admin_username
  password      = var.anypoint_admin_password
  cplane        = "us"
}

resource "anypoint_team_roles" "roles" {
  provider = anypoint.admin   # ← uses the admin alias
  org_id   = var.org_id
  team_id  = var.team_id
  # ...
}
```

## Authentication

The provider supports three authentication modes — pick one by setting `auth_type` (or by which fields you provide).

### `connected_app` (default)

OAuth2 `client_credentials` grant. The provider acts on its own behalf using a Connected App's credentials.

**Required**: `client_id`, `client_secret`.

```terraform
provider "anypoint" {
  client_id     = "..."
  client_secret = "..."
  # auth_type defaults to "connected_app"
}
```

**Connected App requirements**:
- Type: **App acts on its own behalf**
- Scopes assigned for the resources you intend to manage

### `user`

OAuth2 password grant. The provider authenticates on behalf of a user — gets the user's permissions on top of the Connected App's identity. Required for Access Management resources (teams, team roles, team members, connected app scopes management).

**Required**: `client_id`, `client_secret`, `username`, `password`.

```terraform
provider "anypoint" {
  auth_type     = "user"
  client_id     = "..."
  client_secret = "..."
  username      = "tf-bot@example.com"
  password      = "..."
}
```

**Requirements**:
- Connected App of type **App acts on behalf of a user** with the `Resource Owner Password Credentials` grant type enabled
- A service-account user with **MFA disabled** (password grant does not work with MFA)
- The user must have the roles required (e.g. `Organization Administrator`)

> Recommendation: use a Connected App + user pair dedicated **only** for Terraform, with a strong rotated password and minimal-needed roles.

### Pre-signed access token

Bypass authentication by providing a token directly (useful for short-lived runs or when getting a token via an external auth flow):

```terraform
provider "anypoint" {
  access_token = var.anypoint_access_token
}
```

## Control plane (`cplane`)

| `cplane` | Anypoint base URL |
|----------|-------------------|
| `us` (default) | `https://anypoint.mulesoft.com` |
| `eu`           | `https://eu1.anypoint.mulesoft.com` |
| `gov`          | `https://gov.anypoint.mulesoft.com` |

## Environment variables

Each field can be sourced from an environment variable as a fallback:

| Provider field | Env var |
|----------------|---------|
| `client_id`     | `ANYPOINT_CLIENT_ID` |
| `client_secret` | `ANYPOINT_CLIENT_SECRET` |
| `access_token`  | `ANYPOINT_ACCESS_TOKEN` |
| `username`      | `ANYPOINT_USERNAME` |
| `password`      | `ANYPOINT_PASSWORD` |
| `cplane`        | `ANYPOINT_CPLANE` |
| `auth_type`     | `ANYPOINT_AUTH_TYPE` |

## Multiple providers / multi-org setups

You can declare multiple `provider` blocks with different `alias` values to target different orgs, control planes, or auth modes:

```terraform
# Default — connected app for product resources
provider "anypoint" {
  client_id     = var.app_client_id
  client_secret = var.app_client_secret
  cplane        = "us"
}

# Admin — user-on-behalf for Access Management
provider "anypoint" {
  alias         = "admin"
  auth_type     = "user"
  client_id     = var.admin_client_id
  client_secret = var.admin_client_secret
  username      = var.admin_username
  password      = var.admin_password
  cplane        = "us"
}

# EU region — different cplane
provider "anypoint" {
  alias         = "eu"
  client_id     = var.eu_client_id
  client_secret = var.eu_client_secret
  cplane        = "eu"
}

resource "anypoint_team_roles" "roles" {
  provider = anypoint.admin
  # ...
}

resource "anypoint_private_space" "ps_eu" {
  provider = anypoint.eu
  # ...
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `auth_type` (String) Authentication type. One of `"connected_app"` (default) or `"user"`. `"user"` requires `client_id` + `client_secret` + `username` + `password` and uses OAuth2 password grant — needed for Access Management operations.
- `access_token` (String, Sensitive) Pre-signed access token. When set, no authentication call is performed and `client_id`/`client_secret`/`username`/`password` are ignored.
- `client_id` (String, Sensitive) Connected App client ID.
- `client_secret` (String, Sensitive) Connected App client secret.
- `username` (String, Sensitive) Anypoint user username. Only used when `auth_type = "user"`.
- `password` (String, Sensitive) Anypoint user password. Only used when `auth_type = "user"`.
- `cplane` (String) Anypoint control plane: `"us"` (default), `"eu"`, or `"gov"`.
