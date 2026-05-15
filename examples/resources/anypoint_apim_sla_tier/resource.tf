resource "anypoint_apim_sla_tier" "bronze" {
  org_id          = var.root_org
  env_id          = var.env_id
  api_instance_id = anypoint_apim_mule4.my_api.id

  name         = "Bronze"
  description  = "Basic tier — 100 requests per minute"
  auto_approve = true
  status       = "ACTIVE"

  limits {
    time_period_in_milliseconds = 60000
    maximum_requests            = 100
    visible                     = true
  }
}

resource "anypoint_apim_sla_tier" "silver" {
  org_id          = var.root_org
  env_id          = var.env_id
  api_instance_id = anypoint_apim_mule4.my_api.id

  name         = "Silver"
  description  = "1000 req/min + 50000 req/day"
  auto_approve = false
  status       = "ACTIVE"

  limits {
    time_period_in_milliseconds = 60000
    maximum_requests            = 1000
  }

  limits {
    time_period_in_milliseconds = 86400000 # 1 day
    maximum_requests            = 50000
  }
}
