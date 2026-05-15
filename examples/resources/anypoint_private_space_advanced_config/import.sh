# In order for the import to work, you should provide the Private Space ID.

terraform import \
  -var-file params.tfvars.json \                                  #variables file
  anypoint_private_space_advanced_config.demo \                   #resource name
  e60d1779-13eb-4892-9fb3-803087c14988                            #PRIVATE_SPACE_ID
