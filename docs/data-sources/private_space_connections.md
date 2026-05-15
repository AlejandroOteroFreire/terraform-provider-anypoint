---
page_title: "anypoint_private_space_connections Data Source - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Lists all VPN Connections of a Private Space.
---

# anypoint_private_space_connections (Data Source)

Lists all **VPN Connections** of a Private Space — each with its nested VPNs and tunnels.

For a single connection lookup by id, use [`anypoint_private_space_connection`](private_space_connection.md).

## Example Usage

```terraform
data "anypoint_private_space_connections" "all" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
}

# List the names of all configured VPN connections
output "connection_names" {
  value = [for c in data.anypoint_private_space_connections.all.connections : c.name]
}

# Map by name for downstream lookup
output "connection_by_name" {
  value = { for c in data.anypoint_private_space_connections.all.connections : c.name => c.id }
}
```

## Schema

### Required

- `org_id` (String) The organization id where the private space is defined.
- `private_space_id` (String) The unique identifier of the private space.

### Read-Only

- `id` (String) Time-based identifier — refreshed on every read.
- `connections` (List of Object) See [below](#nested-schema-for-connections).

### Nested Schema for `connections`

Read-Only:

- `id` (String) Connection id.
- `name` (String) Connection name.
- `vpns` (List of Object) — same shape as [`anypoint_private_space_connection.vpns`](private_space_connection.md#nested-schema-for-vpns).
