---
page_title: "anypoint_private_space_transit_gateway Resource - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Manages an AWS Transit Gateway attachment to a Private Space.
---

# anypoint_private_space_transit_gateway (Resource)

Manages an **AWS Transit Gateway attachment** to a Private Space. This connects your Anypoint Private Space network to your existing AWS network via AWS Resource Access Manager (RAM) shared Transit Gateway.

The flow is:

1. In your AWS account, create or pick an existing **AWS Transit Gateway**.
2. Share it with the Anypoint AWS account via **AWS Resource Access Manager (RAM)**.
3. Use this resource to declare the attachment from the Private Space side — Anypoint accepts the share, creates the TGW attachment and propagates routes.

> **Unique to this fork.** Not available in upstream nor in the official `mulesoft/terraform-provider-anypoint`.

## Example Usage

```terraform
resource "anypoint_private_space_transit_gateway" "main" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id

  name                   = "prod-aws-tgw"
  resource_share_id      = "arn:aws:ram:us-east-1:123456789012:resource-share/abc12345-..."
  resource_share_account = "123456789012"
  routes                 = ["10.0.0.0/8", "192.168.0.0/16"]
}

# Status is asynchronously computed by AWS — read it back to monitor
output "tgw_status" {
  value = anypoint_private_space_transit_gateway.main.status
}

output "tgw_attachment" {
  value = anypoint_private_space_transit_gateway.main.attachment
}
```

### Tie routes to your VPC CIDRs

```terraform
locals {
  vpc_cidrs = ["10.10.0.0/16", "10.20.0.0/16"]
}

resource "anypoint_private_space_transit_gateway" "main" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id

  name                   = "corp-network"
  resource_share_id      = aws_ram_resource_share.tgw_share.arn
  resource_share_account = data.aws_caller_identity.current.account_id
  routes                 = local.vpc_cidrs
}
```

## Schema

### Required

- `org_id` (String, ForceNew) The organization id where the private space is defined.
- `private_space_id` (String, ForceNew) The unique identifier of the private space.
- `name` (String) Display name of the transit gateway attachment.
- `resource_share_id` (String) The AWS Resource Share ARN/ID (from AWS Resource Access Manager).
- `resource_share_account` (String) The AWS account ID that owns the Transit Gateway and the share.
- `routes` (List of String) CIDR routes to propagate through this transit gateway into the Private Space VPC.

### Read-Only

- `id` (String) Anypoint TGW attachment id.
- `status` (String) Overall status (e.g. `pending`, `available`).
- `attachment` (String) AWS attachment status (e.g. `Pending PrivateSpace Attachment`, `available`).
- `region` (String) AWS region of the transit gateway.

## Import

Import key is composite `<ORG_ID>/<PRIVATE_SPACE_ID>/<TGW_ID>`:

```shell
terraform import anypoint_private_space_transit_gateway.main \
  47ec5a2c.../e60d1779.../a9909a9b-8ebc-457e-82ec-f02428d69395
```
