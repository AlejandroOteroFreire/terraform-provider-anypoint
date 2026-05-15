resource "anypoint_private_space_association" "demo" {
  org_id           = var.root_org
  private_space_id = var.private_space_id

  associations {
    organization_id = var.root_org
    environment     = "production"
  }

  associations {
    organization_id = var.root_org
    environment     = var.sandbox_env_id
  }
}

# Read back the resulting association IDs from the API
output "associations" {
  value = anypoint_private_space_association.demo.created_associations
}
