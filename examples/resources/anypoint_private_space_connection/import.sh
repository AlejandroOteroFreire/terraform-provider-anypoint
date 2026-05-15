# In order for the import to work, you should provide a composite ID:
#  {ORG_ID}/{PRIVATE_SPACE_ID}/{CONNECTION_ID}

terraform import \
  -var-file params.tfvars.json \                                                                                #variables file
  anypoint_private_space_connection.bgp_auto \                                                                  #resource name
  47ec5a2c-c3ce-4994-af2e-11ed525c5b78/e60d1779-13eb-4892-9fb3-803087c14988/71f41a82-a1d0-4ff0-bb92-92a3456bfa03 #composite id
