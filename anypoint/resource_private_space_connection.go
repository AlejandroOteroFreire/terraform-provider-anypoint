package anypoint

/*
BGP con tunel automatico:
{"name":"Prueba",
"vpns":[{"localAsn":"64512","remoteAsn":"65001","remoteIpAddress":"200.155.100.10","vpnTunnels":[{"psk":"","ptpCidr":"","startupAction":"start"},{"psk":"","ptpCidr":"","startupAction":"start"}]}]}

BGP con tunel custom, sin auto inicio:
{"name":"Prueba","vpns":[{"localAsn":"64512","remoteAsn":"65001","remoteIpAddress":"200.155.100.10","vpnTunnels":[{"psk":"12313132132132132132132","ptpCidr":"10.10.20.0/24","startupAction":"add"},{"psk":"13131321321321321331313","ptpCidr":"10.10.21.0/24","startupAction":"add"}]}]}

STatic routes con tunel custom, sin auto inicio:
{"name":"Prueba","vpns":[{"localAsn":"64512","remoteIpAddress":"200.155.100.10","staticRoutes":["10.0.0.0/8"],"vpnTunnels":[{"psk":"12313132132132132132132","ptpCidr":"10.10.20.0/24","startupAction":"add"},{"psk":"13131321321321321331313","ptpCidr":"10.10.21.0/24","startupAction":"add"}]}]}
*/

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/private_space"
)

func resourcePrivateSpaceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateSpaceConnectionCreate,
		ReadContext:   resourcePrivateSpaceConnectionRead,
		UpdateContext: resourcePrivateSpaceConnectionUpdate,
		DeleteContext: resourcePrivateSpaceConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages a VPN Connection within a Private Space.`,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The organization id where the private space is defined.",
			},
			"private_space_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the private space.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the connection.",
			},
			"vpns": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "VPN configuration details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_ip_address": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Remote IP address of the VPN peer.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
						},
						"local_asn": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     64512,
							Description: "Local ASN. Used for both BGP and Static routing.",
						},
						"remote_asn": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Remote ASN. Required for BGP routing. Omit for Static routing.",
						},
						"static_routes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Static CIDR routes. Required for Static routing. Omit for BGP.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"startup_action": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "start",
							Description:      "Tunnel initiation mode: 'start' (automatic) or 'add' (manual). Applies to both tunnels.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"start", "add"}, false)),
						},
						"tunnels": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    2,
							Description: "Custom tunnel configuration. Omit for automatic tunnel setup. If specified, provide exactly 2 tunnels.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"psk": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Pre-Shared Key for the tunnel.",
									},
									"ptp_cidr": {
										Type:             schema.TypeString,
										Optional:         true,
										Description:      "Point-to-Point CIDR for the tunnel (e.g. 169.254.10.0/30).",
										ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
									},
									"is_logs_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_external_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"local_ptp_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_ptp_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status_message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"accepted_route_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"last_status_change": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rekey_margin_in_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rekey_fuzz": {
										Type:     schema.TypeInt,
										Computed: true,
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

func resourcePrivateSpaceConnectionCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	body := buildConnectionBody(d)

	res, httpr, err := pco.privatespaceclient.DefaultApi.CreatePrivateSpaceConnection(authctx, orgid, psid).PrivateSpaceConnectionPostBody(*body).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Connection",
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId(res.GetId())
	return resourcePrivateSpaceConnectionRead(ctx, d, m)
}

func resourcePrivateSpaceConnectionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	connid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	if isComposedResourceId(connid) {
		var decomposeDiags diag.Diagnostics
		orgid, psid, connid, decomposeDiags = decomposePrivateSpaceConnectionId(d)
		if decomposeDiags.HasError() {
			return decomposeDiags
		}
	}

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceConnection(authctx, orgid, psid, connid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Connection " + connid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.Set("name", res.GetName())
	d.Set("org_id", orgid)
	d.Set("private_space_id", psid)
	d.Set("vpns", flattenVpns(res.GetVpns()))
	d.SetId(res.GetId())

	return diags
}

func resourcePrivateSpaceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	connid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	body := buildConnectionBody(d)

	_, httpr, err := pco.privatespaceclient.DefaultApi.UpdatePrivateSpaceConnection(authctx, orgid, psid, connid).PrivateSpaceConnection(*body).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Connection " + connid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	return resourcePrivateSpaceConnectionRead(ctx, d, m)
}

func resourcePrivateSpaceConnectionDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	connid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	httpr, err := pco.privatespaceclient.DefaultApi.DeletePrivateSpaceConnection(authctx, orgid, psid, connid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Connection " + connid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func buildConnectionBody(d *schema.ResourceData) *private_space.PrivateSpaceConnection {
	body := private_space.NewPrivateSpaceConnection()
	body.Name = private_space.PtrString(d.Get("name").(string))

	vpnsList := d.Get("vpns").([]interface{})
	vpns := make([]private_space.PrivateSpaceVpn, len(vpnsList))

	for i, vpnData := range vpnsList {
		vpn := vpnData.(map[string]interface{})

		localAsn := int32(vpn["local_asn"].(int))
		vpnItem := private_space.PrivateSpaceVpn{
			RemoteIpAddress: private_space.PtrString(vpn["remote_ip_address"].(string)),
			LocalAsn:        &localAsn,
		}

		// BGP: include remote_asn only when explicitly set
		if remoteAsn, ok := vpn["remote_asn"].(int); ok && remoteAsn != 0 {
			asn := int32(remoteAsn)
			vpnItem.RemoteAsn = &asn
		}

		// Static routing: include routes when provided
		if routes, ok := vpn["static_routes"].([]interface{}); ok && len(routes) > 0 {
			routeList := make([]string, len(routes))
			for j, r := range routes {
				routeList[j] = r.(string)
			}
			vpnItem.StaticRoutes = routeList
		}

		// Always send 2 tunnels — automatic (empty psk/ptpCidr) or custom
		startupAction := vpn["startup_action"].(string)
		tunnelsList := vpn["tunnels"].([]interface{})

		if len(tunnelsList) == 0 {
			// Automatic tunnel setup
			vpnItem.VpnTunnels = []private_space.PrivateSpaceVpnTunnel{
				{Psk: private_space.PtrString(""), PtpCidr: private_space.PtrString(""), StartupAction: private_space.PtrString(startupAction)},
				{Psk: private_space.PtrString(""), PtpCidr: private_space.PtrString(""), StartupAction: private_space.PtrString(startupAction)},
			}
		} else {
			// Custom tunnel setup
			tunnels := make([]private_space.PrivateSpaceVpnTunnel, len(tunnelsList))
			for j, tunnelData := range tunnelsList {
				t := tunnelData.(map[string]interface{})
				tunnels[j] = private_space.PrivateSpaceVpnTunnel{
					Psk:           private_space.PtrString(t["psk"].(string)),
					PtpCidr:       private_space.PtrString(t["ptp_cidr"].(string)),
					StartupAction: private_space.PtrString(startupAction),
				}
			}
			vpnItem.VpnTunnels = tunnels
		}

		vpns[i] = vpnItem
	}

	body.Vpns = vpns
	return body
}

func flattenVpns(vpns []private_space.PrivateSpaceVpn) []interface{} {
	result := make([]interface{}, len(vpns))

	for i, vpn := range vpns {
		// startup_action is the same on both tunnels — read from first tunnel
		startupAction := "start"
		if len(vpn.GetVpnTunnels()) > 0 {
			startupAction = vpn.VpnTunnels[0].GetStartupAction()
		}

		vpnMap := map[string]interface{}{
			"vpn_id":            vpn.GetVpnId(),
			"vpn_name":          vpn.GetName(),
			"connection_id":     vpn.GetConnectionId(),
			"connection_name":   vpn.GetConnectionName(),
			"remote_ip_address": vpn.GetRemoteIpAddress(),
			"connection_status": vpn.GetVpnConnectionStatus(),
			"static_routes":     vpn.GetStaticRoutes(),
			"startup_action":    startupAction,
			"local_asn":         int(vpn.GetLocalAsn()),
			"remote_asn":        int(vpn.GetRemoteAsn()),
		}

		// Tunnels: only include when custom (non-empty psk or ptp_cidr)
		var tunnels []interface{}
		for _, t := range vpn.GetVpnTunnels() {
			if t.GetPsk() != "" || t.GetPtpCidr() != "" {
				tunnels = append(tunnels, map[string]interface{}{
					"psk":                     t.GetPsk(),
					"ptp_cidr":                t.GetPtpCidr(),
					"is_logs_enabled":         t.GetIsLogsEnabled(),
					"status":                  t.GetStatus(),
					"local_external_ip":       t.GetLocalExternalIpAddress(),
					"local_ptp_ip_address":    t.GetLocalPtpIpAddress(),
					"remote_ptp_ip_address":   t.GetRemotePtpIpAddress(),
					"status_message":          t.GetStatusMessage(),
					"accepted_route_count":    int(t.GetAcceptedRouteCount()),
					"last_status_change":      t.GetLastStatusChange(),
					"rekey_margin_in_seconds": int(t.GetRekeyMarginInSeconds()),
					"rekey_fuzz":              int(t.GetRekeyFuzz()),
				})
			}
		}
		vpnMap["tunnels"] = tunnels

		result[i] = vpnMap
	}

	return result
}

func decomposePrivateSpaceConnectionId(d *schema.ResourceData) (string, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	s := DecomposeResourceId(d.Id())
	if len(s) != 3 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Connection ID format",
			Detail:   fmt.Sprintf("Expected ORG_ID/PRIVATE_SPACE_ID/CONNECTION_ID, got %s", d.Id()),
		})
		return "", "", "", diags
	}
	return s[0], s[1], s[2], diags
}

func parseApiError(httpr *http.Response, err error) string {
	if httpr != nil && httpr.Body != nil {
		b, _ := io.ReadAll(httpr.Body)
		return string(b)
	}
	return err.Error()
}
