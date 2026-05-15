# IMPORTANT: This resource requires the user-on-behalf provider (auth_type = "user").
# See the docs/index.md for how to configure the `anypoint.admin` alias.

resource "anypoint_connected_app_scopes" "datadog" {
  provider = anypoint.admin

  org_id           = var.root_org
  connected_app_id = anypoint_connected_app.datadog.id

  # Org-level scope — only `context_params.org` needed
  scopes {
    scope = "view:monitoring"
    context_params {
      org = var.root_org
    }
  }

  # Org-level
  scopes {
    scope = "read:runtime_fabrics"
    context_params {
      org = var.root_org
    }
  }

  # Env-level scope — `context_params.env_id` required
  scopes {
    scope = "read:applications"
    context_params {
      org    = var.root_org
      env_id = var.sandbox_env_id
    }
  }
}
