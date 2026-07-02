# Changelog

All notable changes to this provider are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

(nothing yet)

## [2.0.5] — 2026-07-02

### Fixed
- `anypoint_bg`: business-group create/update no longer fails with
  "Can not enable entitlement on a business group. It can only be set for a
  master organization". v2.0.4 sent every entitlement unconditionally via
  `d.Get()`, so the master-org-only boolean feature flags (several of which
  default to `true` in the schema — `mqadvancedfeatures`, `designcenter_api/mozart`,
  `runtimefabric`, `runtimefabriccloud`, `anypointsecuritytokenization`,
  `anypointsecurityedgepolicies`) were always included in the request body and
  rejected by the Access Management API. These boolean feature entitlements are
  now only sent when explicitly set in the configuration (read via
  `d.GetRawConfig()` so the schema default no longer leaks through). Quota /
  quantitative entitlements (vCores, VPCs, gateways, MQ, object store, partners,
  managed gateways, etc.) are still sent and remain assignable on a sub-org.

## [2.0.4] — 2026-07-01

### Changed

- `anypoint_bg` (resource and data source): migrated off the external `github.com/mulesoft-anypoint/anypoint-client-go/org` module onto a new internal client (`internal/clients/org`), matching this repo's own-client convention already used by `apim`, `private_space`, `managed_omni_gateway`, `connected_app`, etc. The external module is dropped from `go.mod` entirely.
  - This fixes `entitlements_managed_gateway_small`/`entitlements_managed_gateway_large`, which are now **configurable** (previously read-only) — the old external SDK didn't model these two fields for create/update even though the real Access Management API accepts them (confirmed against `mulesoft/mulesoft-dx`'s `apis/access-management/api.yaml`). Added corresponding read-only `entitlements_managed_gateway_small_reassigned`/`entitlements_managed_gateway_large_reassigned` fields.
  - It also fixes a **dormant bug**: roughly 40 `entitlements_*` attributes were already declared `Optional` in the schema (e.g. `entitlements_mqmessages_base`, `entitlements_designcenter_api`, `entitlements_servicemesh_enabled`, `entitlements_crowd_*`, ...) but were silently never sent to the API on create/update, because the external SDK's `EntitlementsCore` type only modeled 10 of the ~50 entitlement fields the real API supports. All of them are now genuinely applied.

## [2.0.3] — 2026-07-01

### Changed

- `anypoint_bg` (resource and data source): renamed `entitlements_vpns_assigned`/`entitlements_vpns_reassigned` to `entitlements_network_connections_assigned`/`entitlements_network_connections_reassigned` to match the Omni Gateway-era Anypoint API naming.
- Bumped `github.com/mulesoft-anypoint/anypoint-client-go/org` from `v0.4.0` to `v1.2.0`. This renames the internal API methods used by `anypoint_bg` (`OrganizationsPost`→`CreateBG`, `OrganizationsOrgIdGet`→`GetBG`, `OrganizationsOrgIdPut`→`UpdateBG`, `OrganizationsOrgIdDelete`→`DeleteBG`) and switches entitlement sub-objects to the new `*Entitlement`-suffixed types (e.g. `LoadBalancer`→`LoadBalancerEntitlement`). No schema-visible change beyond what's listed below.

### Added

- `anypoint_bg` (resource and data source): expose `entitlements_managed_gateway_small` and `entitlements_managed_gateway_large` as read-only entitlement attributes. (`omni_gateway` was investigated but is **not** an entitlement field in the org API — omitted.)
- `anypoint_private_space_upgrade` (new resource and data source): schedule, cancel and read the status of a CloudHub 2.0 Private Space runtime upgrade, matching the `PATCH/DELETE/GET .../upgrade` and `GET .../upgradestatus` endpoints. Adds `ScheduleUpgrade`, `CancelUpgrade` and `GetUpgradeStatus` to the internal `private_space` client, plus the `PrivateSpaceUpgradeStatus` model.
- `anypoint_managed_omni_gateway` (new resource) and `anypoint_managed_omni_gateway`/`anypoint_managed_omni_gateways` (new data sources): create, read, update, delete and list CloudHub 2.0 **managed** Omni Gateway instances (as opposed to the existing `anypoint_self_managed_omni_gateway`, which is customer-hosted). Backed by a new hand-written internal client (`internal/clients/managed_omni_gateway`) against the Omni Gateway Manager API (`https://anypoint.mulesoft.com/gatewaymanager/xapi/v1`), since this API has no generated SDK yet. Supports `ingress`, `properties`, `logging` and `tracing` configuration blocks, plus a `desired_status` field to start/stop the gateway via the `.../desiredstatus` endpoint.
- `anypoint_agent_instance` and `anypoint_mcp_server` (new resources, "Agents Tools"): manage Agent and MCP server instances deployed to an Omni Gateway target, with `spec`/`endpoint`/`deployment`/`routing` blocks and an optional `gateway_id` shortcut that auto-resolves the deployment target via `anypoint_managed_omni_gateway`. Both reuse the existing `internal/clients/apim` client (same `.../apis` endpoint as `anypoint_apim_mule4`), distinguished only by `endpoint.type` (`a2a` vs `mcp`) — confirmed by reading the real endpoints/schemas from the official `mulesoft/terraform-provider-anypoint` source (`internal/client/agentstools`, `internal/resource/agentstools`). Fixed a pre-existing generator bug in `internal/clients/apim/model_routing_post_body_inner.go`: `Upstreams` was typed as a single object instead of `[]RoutingPostBodyInnerUpstreams`, which would have silently dropped all but one upstream per route; also added the missing `tls_context_id` field.
- `anypoint_agent_instances` and `anypoint_mcp_servers` (new data sources): list Agent/MCP server instances for an environment. Since the API Manager list endpoint doesn't expose `endpoint.type`, this filters candidates by `technology == "flexGateway"` and then fetches each candidate's details to check `endpoint.type` (one extra API call per candidate).

### Fixed

- `anypoint_managed_omni_gateway`: fixed the Omni Gateway Manager API base path from `/gatewaymanager/api/v1` to the correct `/gatewaymanager/xapi/v1`, and added the `X-ANYPNT-ORG-ID`/`X-ANYPNT-ENV-ID` headers the platform expects on every call to this API — both confirmed against the official provider's source after it was found.

## [2.0.2] — 2026-05-18

### Fixed

- `anypoint_connected_app_scopes`: filter out Anypoint-managed scopes (e.g. `profile`) from both the state and the update payload. Anypoint attaches these automatically to every Connected App and cannot be removed via API — including them in the diff produced a permanent drift loop where every `terraform plan` showed `profile` as needing removal. The filter list lives in `anypointAutoScopes` map and can be extended for any other scopes discovered to be auto-managed.

## [2.0.1] — 2026-05-18

### Fixed

- `anypoint_connected_app_scopes`: scopes schema changed from `TypeList` to `TypeSet` with a deterministic hash on `scope + org + env_id`. The Anypoint API returns scopes in non-deterministic order, which caused a permanent drift loop where every `terraform plan` showed the same set of removed+added scopes (e.g. 17 scopes "changed" but only reordered). Now the set comparison is order-independent.
- `anypoint_connected_app_scopes`: the `context_params` block is now always emitted in the state with `env_id=""` for org-scoped scopes (instead of being omitted), matching the shape produced by `dynamic context_params { ... env_id = try(..., "") }` in HCL. Configs should set `env_id = ""` (not `null`) for org-only scopes so the Set hashes line up.

## [2.0.0] — 2026-05-16

First release of the `AlejandroOteroFreire/anypoint` fork.

Contains **breaking changes** as part of the Flex Gateway → Omni Gateway rebrand. Existing users will need to migrate state (see [Migration guide](#migration-from-18x) or the full [INSTALL.md](INSTALL.md#migrating-from-mulesoft-anypointanypoint-18x)).

Distribution moved from the public Terraform Registry (`mulesoft-anypoint/anypoint`) to GitHub Releases at https://github.com/AlejandroOteroFreire/terraform-provider-anypoint/releases.

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
