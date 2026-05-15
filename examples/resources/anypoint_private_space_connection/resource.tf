# BGP with automatic tunnels (Anypoint generates PSKs + ptp CIDRs)
resource "anypoint_private_space_connection" "bgp_auto" {
  org_id           = var.root_org
  private_space_id = var.private_space_id
  name             = "prod-vpn"

  vpns {
    remote_ip_address = "200.10.10.1"
    local_asn         = 64512
    remote_asn        = 65001
    startup_action    = "start"
  }
}

# BGP with custom tunnels (provide your own PSK + ptp CIDR)
resource "anypoint_private_space_connection" "bgp_custom" {
  org_id           = var.root_org
  private_space_id = var.private_space_id
  name             = "prod-vpn-custom"

  vpns {
    remote_ip_address = "200.10.10.1"
    local_asn         = 64512
    remote_asn        = 65001
    startup_action    = "add"

    tunnels {
      psk      = var.tunnel_1_psk
      ptp_cidr = "169.254.10.0/30"
    }

    tunnels {
      psk      = var.tunnel_2_psk
      ptp_cidr = "169.254.11.0/30"
    }
  }
}

# Static routing with custom tunnels
resource "anypoint_private_space_connection" "static" {
  org_id           = var.root_org
  private_space_id = var.private_space_id
  name             = "branch-office-vpn"

  vpns {
    remote_ip_address = "200.155.100.10"
    local_asn         = 64512
    static_routes     = ["10.0.0.0/8", "192.168.0.0/16"]
    startup_action    = "add"

    tunnels {
      psk      = var.tunnel_1_psk
      ptp_cidr = "169.254.20.0/30"
    }

    tunnels {
      psk      = var.tunnel_2_psk
      ptp_cidr = "169.254.21.0/30"
    }
  }
}
