data "anypoint_managed_omni_gateway" "gateway" {
  organization_id = var.root_org
  environment_id  = var.environment_id
  gateway_id      = var.gateway_id
}
