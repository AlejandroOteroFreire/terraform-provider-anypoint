resource "anypoint_mcp_server" "server" {
  org_id = var.root_org
  env_id = var.environment_id

  instance_label = "my-mcp-server"
  gateway_id     = anypoint_managed_omni_gateway.gateway.id

  spec {
    asset_id = "my-mcp-server-spec"
    group_id = var.root_org
    version  = "1.0.0"
  }

  endpoint {
    deployment_type = "CH2"
    base_path       = "my-mcp-server"
  }

  routing {
    label = "default"
    upstreams {
      weight = 100
      uri    = "http://my-mcp-backend.internal:8081"
    }
  }
}
