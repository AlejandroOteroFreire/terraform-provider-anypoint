# In order for the import to work, you should provide a composite ID:
#  {ORG_ID}/{PRIVATE_SPACE_ID}/{TGW_ID}

terraform import \
  -var-file params.tfvars.json \                                                                                #variables file
  anypoint_private_space_transit_gateway.main \                                                                 #resource name
  47ec5a2c-c3ce-4994-af2e-11ed525c5b78/e60d1779-13eb-4892-9fb3-803087c14988/a9909a9b-8ebc-457e-82ec-f02428d69395 #composite id
