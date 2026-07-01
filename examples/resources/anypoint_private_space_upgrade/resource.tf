resource "anypoint_private_space_upgrade" "upgrade" {
  org_id           = var.root_org
  private_space_id = var.private_space_id
  opt_in           = true
  date             = "2026-08-01T00:00:00Z"
}
