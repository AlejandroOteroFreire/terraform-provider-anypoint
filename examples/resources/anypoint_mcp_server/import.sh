# In order for the import to work, you should provide a composite ID:
#  {ORG_ID}/{ENV_ID}/{INSTANCE_ID}

terraform import \
  -var-file params.tfvars.json \                     #variables file
  anypoint_mcp_server.server \                       #resource name
  47ec5a2c-c3ce-4994-af2e-11ed525c5b78/e60d1779-13eb-4892-9fb3-803087c14988/12346 #composite id
