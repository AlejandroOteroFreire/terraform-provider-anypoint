resource "anypoint_secretgroup_certificatepinset" "vendor_api_pins" {
  org_id          = var.root_org
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.sg.id

  name = "vendor-api-pinset"

  # PEM file with one or more concatenated trusted certificates
  certificate_pinset_base64 = filebase64("${path.module}/certs/vendor-pins.pem")
}

output "pinset_algorithm" {
  value = anypoint_secretgroup_certificatepinset.vendor_api_pins.algorithm
}

output "pinset_expiration" {
  value = anypoint_secretgroup_certificatepinset.vendor_api_pins.expiration_date
}
