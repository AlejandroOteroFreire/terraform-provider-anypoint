---
page_title: "anypoint_private_space_transit_gateway Data Source - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Reads a single AWS Transit Gateway attached to a Private Space.
---

# anypoint_private_space_transit_gateway (Data Source)

Reads a single **AWS Transit Gateway attachment** of a Private Space by id. Useful to look up its status, attachment state and routes.

For listing all TGWs, use [`anypoint_private_space_transit_gateways`](private_space_transit_gateways.md).

## Example Usage

```terraform
data "anypoint_private_space_transit_gateway" "main" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
  tgw_id           = "a9909a9b-8ebc-457e-82ec-f02428d69395"
}

output "tgw_attachment_status" {
  value = data.anypoint_private_space_transit_gateway.main.attachment
}
```

## Schema

### Required

- `org_id` (String) The organization id where the private space is defined.
- `private_space_id` (String) The unique identifier of the private space.
- `tgw_id` (String) The unique identifier of the transit gateway.

### Read-Only

- `id` (String) Same as `tgw_id`.
- `name` (String)
- `resource_share_id` (String) AWS RAM resource share ARN/id.
- `resource_share_account` (String) AWS account id that owns the TGW.
- `routes` (List of String) Propagated CIDR routes.
- `status` (String)
- `attachment` (String)
- `region` (String)
