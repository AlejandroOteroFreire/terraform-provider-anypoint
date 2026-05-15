---
page_title: "anypoint_private_space_association Resource - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Manages environment associations for a Private Space. Controls which organizations and environments can use the private space.
---

# anypoint_private_space_association (Resource)

Manages **environment associations** for a Private Space. An association tells Anypoint that a given environment (in the same or a sub-organization) is allowed to deploy workloads on the Private Space.

Use this resource to declaratively control which orgs/environments can use the Private Space. Re-applying with a different `associations` set replaces all bindings.

## Example Usage

### Associate a single environment

```terraform
resource "anypoint_private_space_association" "main" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id

  associations {
    organization_id = var.org_id
    environment     = anypoint_env.production.id
  }
}
```

### Associate all environments of the org

```terraform
resource "anypoint_private_space_association" "all_envs" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id

  associations {
    organization_id = var.org_id
    environment     = "all"
  }
}
```

### Multiple associations (cross-org, multiple envs)

```terraform
resource "anypoint_private_space_association" "multi" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id

  associations {
    organization_id = var.org_id
    environment     = "production"
  }

  associations {
    organization_id = var.suborg_id
    environment     = anypoint_env.suborg_sandbox.id
  }
}
```

## Schema

### Required

- `org_id` (String, ForceNew) The organization id where the private space is defined.
- `private_space_id` (String, ForceNew) The unique identifier of the private space.
- `associations` (Block List) Environment associations to create. See [below](#nested-schema-for-associations).

### Read-Only

- `id` (String) Same as `private_space_id`.
- `created_associations` (List of Object) Associations as created by the API, with their assigned IDs. See [below](#nested-schema-for-created_associations).

### Nested Schema for `associations`

Required:

- `organization_id` (String) The organization ID to associate.
- `environment` (String) Environment to associate. Accepts either an environment UUID, or one of the special values: `"all"`, `"production"`, `"sandbox"`.

### Nested Schema for `created_associations`

Read-Only:

- `id` (String) Assigned association ID.
- `organization_id` (String)
- `environment_id` (String)

## Import

Import is supported using the private space ID as the import key:

```shell
terraform import anypoint_private_space_association.main <private_space_id>
```
