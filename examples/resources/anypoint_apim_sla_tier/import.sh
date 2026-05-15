# In order for the import to work, you should provide a composite ID:
#  {ORG_ID}/{ENV_ID}/{API_INSTANCE_ID}/{TIER_ID}

terraform import \
  -var-file params.tfvars.json \                                                                                  #variables file
  anypoint_apim_sla_tier.bronze \                                                                                 #resource name
  aa1f55d6-213d-4f60-845c-201282484cd1/7074fcdd-9b23-4ab3-97c8-5db5f4adf17d/12345/678                              #composite id
