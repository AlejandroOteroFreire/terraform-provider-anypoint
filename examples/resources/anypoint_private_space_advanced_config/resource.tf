resource "anypoint_private_space_advanced_config" "demo" {
  org_id           = var.root_org
  private_space_id = var.private_space_id

  enable_iam_role = true

  ingress_configuration {
    read_response_timeout = "600"
    protocol              = "https-redirect"
    port_log_level        = "INFO"

    log_filters {
      ip    = "203.0.113.42"
      level = "DEBUG"
    }
  }
}
