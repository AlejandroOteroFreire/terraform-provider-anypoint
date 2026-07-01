resource "anypoint_managed_omni_gateway" "gateway" {
  organization_id = var.root_org
  environment_id  = var.environment_id
  name            = "prod-omni-gw"
  target_id       = var.private_space_id
  release_channel = "lts"
  size            = "small"

  ingress {
    forward_ssl_session = true
    last_mile_security  = true
  }

  properties {
    upstream_response_timeout = 15
    connection_idle_timeout   = 60
  }

  logging {
    level        = "info"
    forward_logs = true
  }

  tracing {
    enabled = false
  }
}
