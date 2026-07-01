# In order for the import to work, you should provide the gateway ID:
#  {GATEWAY_ID}
# NOTE: organization_id and environment_id are not part of the import ID and must be set in the resource config.

terraform import \
  -var-file params.tfvars.json \       #variables file
  anypoint_managed_omni_gateway.gateway \  #resource name
  aa1f55d6-213d-4f60-845c-201282484cd1 #gateway id
