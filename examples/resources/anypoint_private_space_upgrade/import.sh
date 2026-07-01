# In order for the import to work, you should provide the private space ID:
#  {PRIVATE_SPACE_ID}
# NOTE: org_id is not part of the import ID and must be set in the resource config.

terraform import \
  -var-file params.tfvars.json \        #variables file
  anypoint_private_space_upgrade.upgrade \  #resource name
  e60d1779-13eb-4892-9fb3-803087c14988  #private space id
