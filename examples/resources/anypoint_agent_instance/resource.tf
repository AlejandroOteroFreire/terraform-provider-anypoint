resource "anypoint_agent_instance" "agent" {
  org_id = var.root_org
  env_id = var.environment_id

  instance_label = "my-agent"
  gateway_id     = anypoint_managed_omni_gateway.gateway.id

  spec {
    asset_id = "my-agent-spec"
    group_id = var.root_org
    version  = "1.0.0"
  }

  endpoint {
    deployment_type = "CH2"
    base_path       = "my-agent"
  }

  upstream_uri = "http://my-agent-backend.internal:8080"
}
