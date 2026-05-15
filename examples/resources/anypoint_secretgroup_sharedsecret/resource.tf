# UsernamePassword type — DB credentials
resource "anypoint_secretgroup_sharedsecret" "db_creds" {
  org_id          = var.root_org
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.sg.id

  name = "db-creds"
  type = "UsernamePassword"

  username = "appuser"
  password = var.db_password
}

# S3Credential type — AWS access keys
resource "anypoint_secretgroup_sharedsecret" "s3" {
  org_id          = var.root_org
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.sg.id

  name              = "s3-uploader"
  type              = "S3Credential"
  access_key_id     = var.aws_access_key_id
  secret_access_key = var.aws_secret_access_key
  expiration_date   = "2026-12-31T23:59:59Z"
}

# SymmetricKey type — base64-encoded AES key
resource "anypoint_secretgroup_sharedsecret" "aes" {
  org_id          = var.root_org
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.sg.id

  name = "payload-aes-key"
  type = "SymmetricKey"
  key  = filebase64("${path.module}/keys/aes.key")
}

# Blob type — opaque content
resource "anypoint_secretgroup_sharedsecret" "license" {
  org_id          = var.root_org
  env_id          = var.env_id
  secret_group_id = anypoint_secretgroup.sg.id

  name    = "vendor-license"
  type    = "Blob"
  content = file("${path.module}/license.txt")
}
