package anypoint

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/private_space"
)

var PRIVATE_SPACE_CONNECTION_SCHEMA = map[string]*schema.Schema{
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
}

func dataSourcePrivateSpaceConnection() *schema.Resource {
	psc_schema := cloneSchema(PRIVATE_SPACE_CONNECTION_SCHEMA)

	psc_schema["org_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The master organization id where the private space is defined.",
	}
	psc_schema["private_space_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The private space id.",
	}
	// Cambiar 'id' de Computed a Required para poder buscarlo
	psc_schema["connection_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The connection id to retrieve.",
	}

	return &schema.Resource{
		ReadContext: dataSourcePrivateSpaceConnectionRead,
		Description: `
        Reads a ` + "`" + `private space connection` + "`" + ` in your business group.
        `,
		Schema: psc_schema,
	}
}

func dataSourcePrivateSpaceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgId := d.Get("org_id").(string)
	psId := d.Get("private_space_id").(string)
	connId := d.Get("connection_id").(string) // Ahora obtendrá el valor correctamente
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceConnection(authctx, orgId, psId, connId).Execute()
	if err != nil {
		var details string
		if httpr != nil && httpr.StatusCode >= 400 {
			defer httpr.Body.Close()
			b, _ := io.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to Get Private Space Connection for org " + orgId + " and psId " + psId + " and connId " + connId,
			Detail:   details,
		})
		return diags
	}
	defer httpr.Body.Close()
	connection, err := flattenPrivateSpaceConnection(res)

	if err != nil {
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to parse Connection data",
			Detail:   err.Error(),
		})
		return diags
	}
	//save in data source schema
	if err := setPrivateSpaceConnectionCoreAttributesToResourceData(d, connection); err != nil {
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set Connection",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(res.GetId())

	return diags
}

func flattenPrivateSpaceConnection(item *private_space.PrivateSpaceConnection) (map[string]any, error) {
	result := make(map[string]any)

	result["id"] = item.GetId()
	result["name"] = item.GetName()

	// Procesar VPNs
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

		// Procesar Tunnels
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
	result["vpns"] = vpns

	return result, nil
}

func setPrivateSpaceConnectionCoreAttributesToResourceData(d *schema.ResourceData, connection map[string]any) error {
	attributes := getConnectionCoreAttributes()
	if connection != nil {
		for _, attr := range attributes {
			if err := d.Set(attr, connection[attr]); err != nil {
				return fmt.Errorf("unable to set Connection attribute %s\n details: %s", attr, err)
			}
		}
	}
	return nil
}

func getConnectionCoreAttributes() []string {
	attributes := [...]string{
		"id",
		"name",
		"vpns",
	}
	return attributes[:]
}
