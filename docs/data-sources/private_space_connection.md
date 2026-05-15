---
page_title: "anypoint_private_space_connection Data Source - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Reads a single VPN Connection of a Private Space.
---

# anypoint_private_space_connection (Data Source)

Reads a single **VPN Connection** of a Private Space — useful when you need to look up an existing connection (created manually or by another Terraform module) by its ID.

For listing all connections, use [`anypoint_private_space_connections`](private_space_connections.md).

## Example Usage

```terraform
data "anypoint_private_space_connection" "branch" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
  connection_id    = "3457b7a9-76c3-4e1d-b53b-1bf2f496d87e"
}

output "branch_status" {
  value = [for v in data.anypoint_private_space_connection.branch.vpns : v.connection_status]
}

output "branch_public_ips" {
  value = flatten([
    for v in data.anypoint_private_space_connection.branch.vpns : [
      for t in v.tunnels : t.local_external_ip
    ]
  ])
}
```

## Schema

### Required

- `org_id` (String) The organization id where the private space is defined.
- `private_space_id` (String) The private space id.
- `connection_id` (String) The connection id to retrieve.

### Read-Only

- `id` (String) Connection id.
- `name` (String) Connection name.
- `vpns` (List of Object) VPN configurations. See [below](#nested-schema-for-vpns).

### Nested Schema for `vpns`

Read-Only:

- `vpn_id`, `name`, `connection_status` (String)
- `local_asn`, `remote_asn` (Int)
- `remote_ip_address` (String)
- `static_routes` (List of String)
- `tunnels` (List of Object) See [below](#nested-schema-for-tunnels).

#### Nested Schema for `tunnels`

Read-Only:

- `psk` (String, Sensitive)
- `status`, `local_external_ip`, `local_ptp_ip_address`, `remote_ptp_ip_address`, `status_message`, `last_status_change` (String)
- `accepted_route_count`, `rekey_margin_in_seconds`, `rekey_fuzz` (Int)
