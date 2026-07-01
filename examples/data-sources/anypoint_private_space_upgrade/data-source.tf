data "anypoint_private_space_upgrade" "upgrade" {
  org_id           = var.root_org
  private_space_id = var.private_space_id
}
