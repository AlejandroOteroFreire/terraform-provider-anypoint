resource "anypoint_private_space_transit_gateway" "main" {
  org_id           = var.root_org
  private_space_id = var.private_space_id

  name                   = "prod-aws-tgw"
  resource_share_id      = "arn:aws:ram:us-east-1:123456789012:resource-share/abc12345-6789-..."
  resource_share_account = "123456789012"
  routes                 = ["10.0.0.0/8", "192.168.0.0/16"]
}

# Status is computed asynchronously by AWS — read it back to monitor the attachment
output "tgw_status" {
  value = anypoint_private_space_transit_gateway.main.status
}

output "tgw_attachment" {
  value = anypoint_private_space_transit_gateway.main.attachment
}
