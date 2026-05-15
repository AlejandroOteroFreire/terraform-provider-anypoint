# In order for the import to work, you should provide the Connected App ID.

terraform import \
  -var-file params.tfvars.json \                          #variables file
  anypoint_connected_app_scopes.datadog \                 #resource name
  9d756ed0b0ae4ab1976ba52ad2b44752                        #CONNECTED_APP_ID
