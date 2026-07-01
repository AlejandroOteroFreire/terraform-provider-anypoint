data "anypoint_managed_omni_gateways" "gateways" {
  organization_id = var.root_org
  environment_id  = var.environment_id
}
