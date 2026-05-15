---
page_title: "anypoint_private_space_connection Resource - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Manages a VPN Connection within a Private Space.
---

# anypoint_private_space_connection (Resource)

Manages a **VPN Connection** within a Private Space. Each connection contains one or more VPN tunnels with BGP or static routing — Anypoint uses these connections to extend the Private Space network to your own on-prem or cloud environments via IPsec.

A connection always contains **exactly 2 IPsec tunnels** (for redundancy). You can either:
- Let Anypoint **auto-generate** tunnel parameters (PSK + point-to-point CIDR) by omitting the `tunnels` blocks
- Provide your own **custom** PSK and ptp-CIDR (when your on-prem device requires specific values)

Two routing modes:
- **BGP** — set `remote_asn`. Routes are advertised dynamically.
- **Static** — set `static_routes`. Manually declared CIDRs.

> **Unique to this fork.** Not available in the upstream provider — fully reimplemented here including BGP/static routing, dual-tunnel and custom PSK/CIDR support.

## Example Usage

### BGP with automatic tunnels

```terraform
resource "anypoint_private_space_connection" "bgp_auto" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
  name             = "prod-vpn"

  vpns {
    remote_ip_address = "200.10.10.1"
    local_asn         = 64512
    remote_asn        = 65001
    startup_action    = "start"
    # No `tunnels` blocks → Anypoint generates PSKs + ptp CIDRs automatically
  }
}
```

### BGP with custom tunnels

```terraform
resource "anypoint_private_space_connection" "bgp_custom" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
  name             = "prod-vpn-custom"

  vpns {
    remote_ip_address = "200.10.10.1"
    local_asn         = 64512
    remote_asn        = 65001
    startup_action    = "add"      # "add" = manual init from your side

    tunnels {
      psk      = var.tunnel_1_psk
      ptp_cidr = "169.254.10.0/30"
    }

    tunnels {
      psk      = var.tunnel_2_psk
      ptp_cidr = "169.254.11.0/30"
    }
  }
}
```

### Static routing with custom tunnels

```terraform
resource "anypoint_private_space_connection" "static" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
  name             = "branch-office-vpn"

  vpns {
    remote_ip_address = "200.155.100.10"
    local_asn         = 64512
    static_routes     = ["10.0.0.0/8", "192.168.0.0/16"]
    startup_action    = "add"

    tunnels {
      psk      = var.tunnel_1_psk
      ptp_cidr = "169.254.20.0/30"
    }

    tunnels {
      psk      = var.tunnel_2_psk
      ptp_cidr = "169.254.21.0/30"
    }
  }
}
```

## Schema

### Required

- `org_id` (String, ForceNew) The organization id where the private space is defined.
- `private_space_id` (String, ForceNew) The unique identifier of the private space.
- `name` (String) The connection name.
- `vpns` (Block List) VPN configurations (typically one). See [below](#nested-schema-for-vpns).

### Read-Only

- `id` (String) Connection id.

### Nested Schema for `vpns`

Required:

- `remote_ip_address` (String) Remote IP address of the VPN peer.

Optional:

- `local_asn` (Int, Default `64512`) Local ASN. Used for both BGP and Static routing.
- `remote_asn` (Int) Remote ASN. **Required for BGP. Omit for Static.**
- `static_routes` (List of String) Static CIDR routes. **Required for Static. Omit for BGP.**
- `startup_action` (String, Default `"start"`) Tunnel initiation mode. Applies to **both** tunnels:
  - `"start"` — Anypoint initiates the IPsec negotiation (automatic mode)
  - `"add"` — your side initiates (manual mode — useful when behind firewalls / NAT)
- `tunnels` (Block List, Max 2) Custom tunnel parameters. **Omit entirely** for automatic generation, or **provide exactly 2 blocks** for custom config.

Read-Only:

- `vpn_id`, `vpn_name`, `connection_id`, `connection_name`, `connection_status` (String)

#### Nested Schema for `tunnels`

Optional:

- `psk` (String, Sensitive) Pre-Shared Key.
- `ptp_cidr` (String) Point-to-Point CIDR (e.g. `169.254.10.0/30`). Must be a `/30` from the `169.254.0.0/16` link-local range.

Read-Only:

- `is_logs_enabled` (Bool)
- `status` (String)
- `local_external_ip` (String) Public IP of the Anypoint endpoint.
- `local_ptp_ip_address`, `remote_ptp_ip_address` (String)
- `status_message` (String)
- `accepted_route_count` (Int)
- `last_status_change` (String)
- `rekey_margin_in_seconds`, `rekey_fuzz` (Int)

## Import

Import key is composite `<ORG_ID>/<PRIVATE_SPACE_ID>/<CONNECTION_ID>`:

```shell
terraform import anypoint_private_space_connection.bgp_auto \
  47ec5a2c.../e60d1779.../71f41a82...
```
