---
page_title: "anypoint_private_space_associations Data Source - terraform-provider-anypoint"
subcategory: "CloudHub 2.0"
description: |-
  Reads all environment associations for a Private Space.
---

# anypoint_private_space_associations (Data Source)

Lists all environment associations attached to a Private Space — including their IDs, the originating organization, and the bound environment.

Useful to inspect existing bindings (for example, to feed downstream resources or for drift detection).

## Example Usage

```terraform
data "anypoint_private_space_associations" "current" {
  org_id           = var.org_id
  private_space_id = anypoint_private_space.demo.id
}

output "associated_envs" {
  value = [for a in data.anypoint_private_space_associations.current.associations : a.environment_id]
}
```

## Schema

### Required

- `org_id` (String) The organization id where the private space is defined.
- `private_space_id` (String) The unique identifier of the private space.

### Read-Only

- `id` (String) Time-based identifier — refreshed on every read.
- `associations` (List of Object) See [below](#nested-schema-for-associations).

### Nested Schema for `associations`

Read-Only:

- `id` (String) Association ID.
- `organization_id` (String) The org that this association grants access to.
- `environment_id` (String) The environment that this association grants access to.
