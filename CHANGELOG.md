# Changelog

All notable changes to this provider are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] — 2.0.0 (planned)

This release contains **breaking changes** as part of the Flex Gateway → Omni Gateway rebrand. Existing users will need to migrate state (see [Migration guide](#migration-from-18x)).

### Breaking changes

- **Rebrand**: `anypoint_apim_flexgateway` → `anypoint_self_managed_omni_gateway` (resource)
- **Rebrand**: `anypoint_flexgateway_target` → `anypoint_self_managed_omni_gateway_target` (data source)
- **Rebrand**: `anypoint_flexgateway_targets` → `anypoint_self_managed_omni_gateway_targets` (data source)
- **Rebrand**: `anypoint_flexgateway_registration_token` → `anypoint_self_managed_omni_gateway_registration_token` (data source)
- **Rebrand**: `anypoint_secretgroup_tlscontext_flexgateway` → `anypoint_secretgroup_tlscontext_self_managed_omni_gateway` (resource + data source)

> The Anypoint API endpoint paths (`/flex-gateway-targets`) and the response field `technology = "flexGateway"` remain unchanged — MuleSoft kept the wire format. Only Terraform-facing names changed.

### Added

#### Authentication
- New provider field `auth_type` with values `"connected_app"` (default) and `"user"`.
- **OAuth2 password grant** (`auth_type = "user"`) — authenticates on behalf of a user, required for Access Management operations (teams, team roles, team members, connected app scopes). Hits the same `/accounts/api/v2/oauth2/token` endpoint with `grant_type=password`.
- `username` and `password` fields are **no longer deprecated** when used with `auth_type = "user"`.

#### Resources & data sources — ported from `mulesoft/terraform-provider-anypoint`
- `anypoint_private_space_association` (R) + `anypoint_private_space_associations` (DS)
- `anypoint_private_space_advanced_config` (R)
- `anypoint_connected_app_scopes` (R) — manages contextual scopes for Connected Apps with full `org` / `env_id` context support
- `anypoint_apim_sla_tier` (R) — SLA tiers for API Manager instances
- `anypoint_secretgroup_sharedsecret` (R) — Shared Secrets in 4 types: `UsernamePassword`, `S3Credential`, `SymmetricKey`, `Blob`
- `anypoint_secretgroup_certificatepinset` (R) — Certificate Pinsets via `multipart/form-data`

#### Resources & data sources — original to this fork
- `anypoint_private_space_connection` (R) + `anypoint_private_space_connection` (DS) + `anypoint_private_space_connections` (DS) — full VPN connection management: BGP/static routing, dual-tunnel, automatic or custom PSK + ptp-CIDR
- `anypoint_private_space_transit_gateway` (R) + `anypoint_private_space_transit_gateway` (DS) + `anypoint_private_space_transit_gateways` (DS) — AWS Transit Gateway attachments to Private Spaces via AWS RAM

#### Go client modules
- New: `secretgroup_sharedsecret/` (full CRUD over `/sharedSecrets`)
- New: `secretgroup_certificatepinset/` (multipart/form-data create, JSON read/delete)
- Extended `private_space/`: associations, advanced config, connection (BGP/static/dual-tunnel), transit gateway (AWS RAM)
- Extended `apim/`: SLA tiers (full CRUD)
- Extended `authorization/`: `Credentials` model now includes `username`/`password` for password grant
- Extended `connected_app/`: scopes management (uses existing types)

### Changed

- **Documentation**:
  - `README.md` rewritten with TOC, 3 auth modes, architecture diagrams, resources by category
  - `docs/index.md` updated with `auth_type`, multi-provider patterns, env vars table
  - 23 new `.md` docs for new/ported resources and data sources
  - 6 docs renamed (`flexgateway*.md` → `self_managed_omni_gateway*.md`)
  - 9 new `examples/resources/<name>/` directories with `resource.tf` + `import.sh`

### Fixed

- `local_asn` in `anypoint_private_space_connection`: was declared `Required` with a `Default` value (Terraform schema invalid). Fixed to `Optional + Default`.
- `anypoint_private_space_connection`: API response handling — ASN fields are integers in the GET response but strings in the POST body. Added `MarshalJSON` on `PrivateSpaceVpn` to serialize ASN as strings while keeping `*int32` for unmarshaling.
- `anypoint_private_space_connection`: `vpnConnectionStatus` JSON tag was incorrect (`connectionStatus`), now matches the actual API response.
- `anypoint_private_space_transit_gateway`: model was over-engineered with nested `spec`/`status` blocks that don't exist in the API. Flattened to match the real wire format (`resourceShareId`, `resourceShareAccount`, `routes` as top-level fields).

## Migration from 1.8.x

If you are upgrading from `1.8.x` and have any of the renamed resources in state, run `terraform state mv` before the next `terraform apply`:

```bash
# Resource rename
terraform state mv \
  'anypoint_apim_flexgateway.example' \
  'anypoint_self_managed_omni_gateway.example'

# TLS context rename
terraform state mv \
  'anypoint_secretgroup_tlscontext_flexgateway.example' \
  'anypoint_secretgroup_tlscontext_self_managed_omni_gateway.example'
```

For `for_each` usage, repeat per key. For data sources, just update the references in your `.tf` files — no state surgery needed.

### Using `auth_type = "user"`

To use Access Management resources (`anypoint_team_roles`, `anypoint_connected_app_scopes`, etc.) declare a second provider alias:

```hcl
provider "anypoint" {
  client_id     = var.client_id
  client_secret = var.client_secret
  cplane        = "us"
}

provider "anypoint" {
  alias         = "admin"
  auth_type     = "user"
  client_id     = var.admin_client_id        # Connected App "acts on behalf of a user"
  client_secret = var.admin_client_secret
  username      = var.admin_username         # service account, MFA disabled
  password      = var.admin_password
  cplane        = "us"
}

resource "anypoint_team_roles" "roles" {
  provider = anypoint.admin
  # ...
}
```

See [`docs/index.md`](docs/index.md) for the full description.

---

## [1.8.2] and earlier

See the [GitHub release notes](https://github.com/mulesoft-anypoint/terraform-provider-anypoint/releases) for `v1.8.2` and prior versions.
