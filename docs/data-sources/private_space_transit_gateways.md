---
page_title: "anypoint_private_space_transit_gateways Data Source - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Lists all AWS Transit Gateways attached to a Private Space.
---

# anypoint_private_space_transit_gateways (Data Source)

Lists all **AWS Transit Gateway attachments** of a Private Space.

For a single TGW lookup by id, use [`anypoint_private_space_transit_gateway`](private_space_transit_gateway.md).

## Example Usage

```terraform
data "anypoint_private_space_transit_gateways" "all" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
}

output "tgw_names" {
  value = [for t in data.anypoint_private_space_transit_gateways.all.transit_gateways : t.name]
}

output "tgw_status_map" {
  value = { for t in data.anypoint_private_space_transit_gateways.all.transit_gateways : t.name => t.status }
}
```

## Schema

### Required

- `org_id` (String) The organization id where the private space is defined.
- `private_space_id` (String) The unique identifier of the private space.

### Read-Only

- `id` (String) Time-based identifier — refreshed on every read.
- `transit_gateways` (List of Object) See [below](#nested-schema-for-transit_gateways).

### Nested Schema for `transit_gateways`

Read-Only:

- `id` (String) TGW attachment id.
- `name` (String)
- `resource_share_id` (String)
- `resource_share_account` (String)
- `routes` (List of String)
- `status` (String)
- `attachment` (String)
- `region` (String)
