# Terraform Provider Anypoint

![Terraform](https://img.shields.io/badge/terraform->=1.0.x-623CE4?logo=terraform)
![Go](https://img.shields.io/badge/go-1.20%2B-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MPL%202.0-blue)
![Discord](https://img.shields.io/badge/Discord-join-7289DA?logo=discord&logoColor=white)

A community-driven Terraform provider for managing **MuleSoft Anypoint Platform** resources as code.

Whether you need to provision Private Spaces, configure API Manager policies, manage Access Management (teams, users, roles), or wire up entire Connected Apps with their scopes — this provider lets you describe it all declaratively.

> **Note:** this is a **community fork** (mulesoft-anypoint org) with extensions on top of the upstream. It is not an official MuleSoft product. See [Disclaimer](#disclaimer).

---

## Table of Contents

- [Why this provider?](#why-this-provider)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Architecture / How it works](#architecture--how-it-works)
- [Resources by Category](#resources-by-category)
- [Building from source](#building-from-source)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [Disclaimer](#disclaimer)

---

## Why this provider?

- **Full lifecycle**: declarative management of Private Spaces, Cloudhub 1.0/2.0, RTF, API Manager, Anypoint MQ, Secrets Manager, Access Management, Identity Providers, and more
- **Two auth modes**: standard `client_credentials` for automation, plus `auth_type = "user"` (OAuth2 password grant) for Access Management operations that need elevated permissions
- **Extra resources** beyond the upstream: full Private Space VPN Connection management (BGP / static, dual-tunnel), AWS Transit Gateway attachments, SLA Tiers, Shared Secrets, Certificate Pinsets, and more

## Installation

This fork is **not published** on the public Terraform Registry. Install it from the GitHub Releases page.

> 📘 **Full installation guide** with CI/CD examples, migration from `mulesoft-anypoint/anypoint` 1.8.x, troubleshooting, and the `dev_overrides` workflow: **[INSTALL.md](INSTALL.md)**.

### Quick install — from GitHub Releases

1. Download the release archive matching your OS/arch from the [Releases page](https://github.com/AlejandroOteroFreire/terraform-provider-anypoint/releases) — for example:
   - macOS Apple Silicon: `terraform-provider-anypoint_2.0.0_darwin_arm64.zip`
   - Linux amd64: `terraform-provider-anypoint_2.0.0_linux_amd64.zip`
2. (Optional) Verify the checksum:
   ```bash
   curl -sL <SHA256SUMS-url> | shasum -c -a 256 --ignore-missing
   ```
3. Unzip into Terraform's plugin directory using the expected layout:
   ```bash
   VERSION=2.0.0
   OS_ARCH=darwin_arm64   # adjust for your platform
   NAMESPACE=AlejandroOteroFreire
   TARGET="$HOME/.terraform.d/plugins/registry.terraform.io/${NAMESPACE}/anypoint/${VERSION}/${OS_ARCH}"

   mkdir -p "$TARGET"
   unzip -o terraform-provider-anypoint_${VERSION}_${OS_ARCH}.zip -d "$TARGET"
   chmod +x "$TARGET/terraform-provider-anypoint_v${VERSION}"
   ```
4. Use it in your config:
   ```hcl
   terraform {
     required_providers {
       anypoint = {
         source  = "AlejandroOteroFreire/anypoint"
         version = "2.0.0"
       }
     }
   }
   ```

### Local development (no install)

Add to your `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "AlejandroOteroFreire/anypoint" = "/absolute/path/to/terraform-provider-anypoint-fork"
  }
  direct {}
}
```

Then run `go build -o terraform-provider-anypoint` inside the fork — Terraform picks up the local binary directly, no `terraform init` needed.

## Quick Start

```hcl
terraform {
  required_providers {
    anypoint = {
      source  = "AlejandroOteroFreire/anypoint"
      version = "~> 2.0"
    }
  }
}

provider "anypoint" {
  client_id     = var.anypoint_client_id      # or ANYPOINT_CLIENT_ID env var
  client_secret = var.anypoint_client_secret  # or ANYPOINT_CLIENT_SECRET env var
  cplane        = "us"                        # "us" (default), "eu" or "gov"
}

# Example: list all Private Spaces in your org
data "anypoint_private_spaces" "all" {
  org_id = var.org_id
}

output "spaces" {
  value = data.anypoint_private_spaces.all.private_spaces
}
```

## Authentication

The provider supports **three** authentication flows. Use the one that matches your scenario.

### 1. Connected App — `auth_type = "connected_app"` (default)

OAuth2 `client_credentials` grant. Best for automation, CI/CD pipelines, and any scenario where the provider acts on its own behalf.

```hcl
provider "anypoint" {
  client_id     = var.anypoint_client_id
  client_secret = var.anypoint_client_secret
  cplane        = "us"
  # auth_type defaults to "connected_app"
}
```

**Requirements**: a Connected App in your root org with **"App acts on its own behalf"** type and the scopes required for the resources you manage.

### 2. User on behalf — `auth_type = "user"`

OAuth2 password grant. Required for **Access Management** operations (teams, team roles, team members, connected app scopes) — these endpoints don't accept connected-app-only tokens.

```hcl
provider "anypoint" {
  alias         = "admin"
  auth_type     = "user"
  client_id     = var.anypoint_admin_client_id      # Connected App with "App acts on behalf of a user"
  client_secret = var.anypoint_admin_client_secret
  username      = var.anypoint_admin_username       # a service-account user (MFA disabled)
  password      = var.anypoint_admin_password
  cplane        = "us"
}

# Use it explicitly on Access Management resources:
resource "anypoint_team_roles" "roles" {
  provider = anypoint.admin
  # ...
}
```

**Requirements**:
- Connected App of type **"App acts on behalf of a user"** with `Resource Owner Password Credentials` grant enabled
- A user (service account) with **MFA disabled** — password grant doesn't work with MFA
- The user must have the roles needed for the operations (typically `Organization Administrator`)

### 3. Pre-signed token

If you already have a bearer token (e.g., from an external auth flow):

```hcl
provider "anypoint" {
  access_token = var.anypoint_access_token   # or ANYPOINT_ACCESS_TOKEN env var
  cplane       = "us"
}
```

### Control plane (`cplane`)

| Value | Anypoint Region |
|-------|-----------------|
| `us` (default) | `https://anypoint.mulesoft.com` |
| `eu`           | `https://eu1.anypoint.mulesoft.com` |
| `gov`          | `https://gov.anypoint.mulesoft.com` |

## Architecture / How it works

The provider is a thin Terraform adapter over a set of generated Go client libraries — one per Anypoint API surface.

![alt text](drive/imgs/provider-arch.png)

The Go client (`anypoint-client-go`) is **generated from OpenAPI 3 specifications** that the community contributes via [`anypoint-automation-client-generator`](https://github.com/mulesoft-anypoint/anypoint-automation-client-generator). Each spec ships at least `GET`, `POST` and `DELETE` operations and gets pushed to [`anypoint-client-go`](https://github.com/mulesoft-anypoint/anypoint-client-go).

![alt text](drive/imgs/provider-cycle.png)

For new contributions, the cycle is:
1. Pick an Anypoint resource and reverse-engineer its API (Postman + browser inspector + Anypoint docs)
2. Write an OpenAPI 3 spec and submit it to the generator repo — a Go module gets generated and published
3. Implement the resource and matching data sources in this provider using the generated module

## Resources by Category

> See [`docs/resources/`](docs/resources/) and [`docs/data-sources/`](docs/data-sources/) for the full reference of every resource and data source.

### Access Management
`anypoint_bg`, `anypoint_env`, `anypoint_user`, `anypoint_team`, `anypoint_team_member`, `anypoint_team_roles`, `anypoint_team_group_mappings`, `anypoint_rolegroup`, `anypoint_rolegroup_roles`, `anypoint_user_rolegroup`, `anypoint_connected_app`, `anypoint_connected_app_scopes`, `anypoint_idp_oidc`, `anypoint_idp_saml`

### CloudHub 2.0 / Private Spaces
`anypoint_private_space`, `anypoint_private_space_tlscontext_pem`, `anypoint_private_space_tlscontext_jks`, `anypoint_private_space_connection`, `anypoint_private_space_transit_gateway`, `anypoint_private_space_association`, `anypoint_private_space_advanced_config`, `anypoint_cloudhub2_shared_space_deployment`

### CloudHub 1.0 / VPC
`anypoint_vpc`, `anypoint_vpn`, `anypoint_dlb`

### Runtime Fabric
`anypoint_fabrics`, `anypoint_fabrics_associations`, `anypoint_rtf_deployment`

### API Manager
`anypoint_apim_mule4`, `anypoint_self_managed_omni_gateway` (formerly `apim_flexgateway`), `anypoint_apim_sla_tier`, `anypoint_apim_policy_basic_auth`, `anypoint_apim_policy_client_id_enforcement`, `anypoint_apim_policy_custom`, `anypoint_apim_policy_jwt_validation`, `anypoint_apim_policy_message_logging`, `anypoint_apim_policy_rate_limiting`

### Anypoint MQ
`anypoint_amq`, `anypoint_ame`, `anypoint_ame_binding`

### Secrets Manager
`anypoint_secretgroup`, `anypoint_secretgroup_keystore`, `anypoint_secretgroup_truststore`, `anypoint_secretgroup_certificate`, `anypoint_secretgroup_certificatepinset`, `anypoint_secretgroup_sharedsecret`, `anypoint_secretgroup_crldistrib_cfgs`, `anypoint_secretgroup_tlscontext_mule`, `anypoint_secretgroup_tlscontext_self_managed_omni_gateway` (formerly `secretgroup_tlscontext_flexgateway`), `anypoint_secretgroup_tlscontext_securityfabric`

## Building from source

### Prerequisites

- Go 1.20+
- Terraform 1.0+

```bash
go build -o terraform-provider-anypoint
```

If you are pulling private modules:

```bash
go env -w GOPRIVATE=github.com/mulesoft-anypoint
```

### Local install for development

Use `dev_overrides` in your `~/.terraformrc` to test your local build without publishing:

```hcl
provider_installation {
  dev_overrides {
    "mulesoft-anypoint/anypoint" = "/absolute/path/to/terraform-provider-anypoint-fork"
  }
  direct {}
}
```

With `dev_overrides` you skip `terraform init` — your local binary is picked up directly.

### Run a sample

```bash
cd examples
# edit variables / credentials in main.tf
terraform plan
```

## Documentation

- Per-resource reference: [`docs/resources/`](docs/resources/)
- Per-data-source reference: [`docs/data-sources/`](docs/data-sources/)
- Provider configuration: [`docs/index.md`](docs/index.md)

### How docs are organized

The current docs (per-resource `.md` files) are **handwritten and committed** to `docs/`. The provider-level `docs/index.md` is generated from [`templates/index.md.tmpl`](templates/index.md.tmpl) by [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs) — keep them in sync if you edit one.

### Regenerating `docs/index.md`

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
tfplugindocs generate
```

This (re)generates `docs/index.md` from the template — current per-resource docs are **left untouched** because there are no per-resource templates yet.

### Adding a new resource — what to write

When you add a new `resource_*.go` or `data_source_*.go`, also add by hand:

1. `docs/resources/<name>.md` (or `docs/data-sources/<name>.md`) — follow the format of an existing doc (e.g. [`docs/resources/private_space_association.md`](docs/resources/private_space_association.md)). YAML front-matter + `## Example Usage` + `## Schema` + `## Import`.
2. `examples/resources/<full_name>/resource.tf` (or `examples/data-sources/<full_name>/data-source.tf`) — a real, working configuration.
3. `examples/resources/<full_name>/import.sh` if the resource supports import.

> Long-term goal: convert handwritten docs to `templates/{resources,data-sources}/<name>.md.tmpl` so `tfplugindocs generate` can keep the schema sections in sync with the Go code automatically. PRs welcome.

## Contributing

You can help by:

- **Testing** the provider and opening issues for any problem you find
- Contributing **OpenAPI specs** for new resources [here](https://github.com/mulesoft-anypoint/anypoint-automation-client-generator)
- Contributing **code** in this provider for new resources/data sources

### Release

Follow the [Terraform Registry publishing guide](https://www.terraform.io/docs/registry/providers/publishing.html#using-goreleaser-locally).

## Credits

Built with love by the MuleSoft community.

## Disclaimer

**This is [Open Source Software — please review the considerations](LICENSE.md).** This is an open source project. It does not form part of the official MuleSoft product stack and is therefore not included in MuleSoft support SLAs. Issues should be directed to the community, who will try to assist on a best-endeavours basis. This application is distributed **as is**.

Let's automate, simplify, and supercharge your Anypoint deployments with Terraform!
