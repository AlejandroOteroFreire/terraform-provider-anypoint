package anypoint

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/private_space"
)

func dataSourcePrivateSpaceConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateSpaceConnectionsRead,
		Description: `Reads all connections (VPN, Peering, TGW) for a specific Private Space.`,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization id where the private space is defined.",
			},
			"private_space_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the private space.",
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   {Type: schema.TypeString, Computed: true},
						"name": {Type: schema.TypeString, Computed: true},
						"vpns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpn_id":            {Type: schema.TypeString, Computed: true},
									"name":              {Type: schema.TypeString, Computed: true},
									"connection_status": {Type: schema.TypeString, Computed: true},
									"local_asn":         {Type: schema.TypeInt, Computed: true},
									"remote_asn":        {Type: schema.TypeInt, Computed: true},
									"remote_ip_address": {Type: schema.TypeString, Computed: true},
									"static_routes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"tunnels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"psk":                     {Type: schema.TypeString, Computed: true, Sensitive: true},
												"status":                  {Type: schema.TypeString, Computed: true},
												"local_external_ip":       {Type: schema.TypeString, Computed: true},
												"local_ptp_ip_address":    {Type: schema.TypeString, Computed: true},
												"remote_ptp_ip_address":   {Type: schema.TypeString, Computed: true},
												"status_message":          {Type: schema.TypeString, Computed: true},
												"accepted_route_count":    {Type: schema.TypeInt, Computed: true},
												"last_status_change":      {Type: schema.TypeString, Computed: true},
												"rekey_margin_in_seconds": {Type: schema.TypeInt, Computed: true},
												"rekey_fuzz":              {Type: schema.TypeInt, Computed: true},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourcePrivateSpaceConnectionsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceConnections(authctx, orgid, psid).Execute()

	if err != nil {
		var details string
		if httpr != nil && httpr.StatusCode >= 400 {
			defer httpr.Body.Close()
			b, _ := io.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Unable to Get Connections for PS " + psid,
			Detail:   details,
		}}
	}
	defer httpr.Body.Close()

	if err := d.Set("connections", flattenConnectionsData(res)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func flattenConnectionsData(connections []private_space.PrivateSpaceConnection) []interface{} {
	res := make([]interface{}, len(connections))
	for i, item := range connections {
		conn := make(map[string]interface{})
		conn["id"] = item.GetId()
		conn["name"] = item.GetName()

		vpns := make([]interface{}, len(item.GetVpns()))
		for j, vpnItem := range item.GetVpns() {
			v := make(map[string]interface{})
			v["name"] = vpnItem.GetName()
			v["vpn_id"] = vpnItem.GetVpnId()
			v["connection_status"] = vpnItem.GetVpnConnectionStatus()
			v["local_asn"] = vpnItem.GetLocalAsn()
			v["remote_asn"] = vpnItem.GetRemoteAsn()
			v["remote_ip_address"] = vpnItem.GetRemoteIpAddress()
			v["static_routes"] = vpnItem.GetStaticRoutes()

			tunnels := make([]interface{}, len(vpnItem.GetVpnTunnels()))
			for k, tunnelItem := range vpnItem.GetVpnTunnels() {
				t := make(map[string]interface{})
				t["psk"] = tunnelItem.GetPsk()
				t["local_external_ip"] = tunnelItem.GetLocalExternalIpAddress()
				t["status"] = tunnelItem.GetStatus()
				t["local_ptp_ip_address"] = tunnelItem.GetLocalPtpIpAddress()
				t["remote_ptp_ip_address"] = tunnelItem.GetRemotePtpIpAddress()
				t["status_message"] = tunnelItem.GetStatusMessage()
				t["accepted_route_count"] = tunnelItem.GetAcceptedRouteCount()
				t["last_status_change"] = tunnelItem.GetLastStatusChange()
				t["rekey_margin_in_seconds"] = tunnelItem.GetRekeyMarginInSeconds()
				t["rekey_fuzz"] = tunnelItem.GetRekeyFuzz()

				tunnels[k] = t
			}
			v["tunnels"] = tunnels
			vpns[j] = v
		}
		conn["vpns"] = vpns
		res[i] = conn
	}
	return res
}
