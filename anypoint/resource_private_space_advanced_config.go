package anypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/anypoint-client-go/private_space"
)

func resourcePrivateSpaceAdvancedConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateSpaceAdvancedConfigUpdate,
		ReadContext:   resourcePrivateSpaceAdvancedConfigRead,
		UpdateContext: resourcePrivateSpaceAdvancedConfigUpdate,
		DeleteContext: resourcePrivateSpaceAdvancedConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages advanced configuration for a Private Space: ingress settings and IAM role.`,
		Schema: map[string]*schema.Schema{
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
			"enable_iam_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable IAM role for the private space.",
			},
			"ingress_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_response_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "300",
							Description: "Read response timeout in seconds.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "https-redirect",
							Description: "Ingress protocol (e.g. 'https-redirect', 'https').",
						},
						"port_log_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "ERROR",
							Description: "Port log level (e.g. 'ERROR', 'INFO', 'DEBUG').",
						},
						"log_filters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "IP address for the log filter.",
									},
									"level": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Log level for this IP filter.",
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

func resourcePrivateSpaceAdvancedConfigUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	body := buildAdvancedConfigBody(d)
	httpr, err := pco.privatespaceclient.DefaultApi.UpdatePrivateSpaceAdvancedConfig(authctx, orgid, psid).PrivateSpaceAdvancedConfigBody(*body).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Private Space Advanced Config",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId(psid)
	return resourcePrivateSpaceAdvancedConfigRead(ctx, d, m)
}

func resourcePrivateSpaceAdvancedConfigRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpace(authctx, orgid, psid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Private Space for Advanced Config",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.Set("enable_iam_role", res.GetEnableIAMRole())
	d.Set("private_space_id", psid)
	d.SetId(psid)
	return diags
}

func resourcePrivateSpaceAdvancedConfigDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// Reset to defaults on destroy
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	defaults := private_space.NewPrivateSpaceAdvancedConfig()
	defaults.EnableIAMRole = private_space.PtrBool(false)

	httpr, err := pco.privatespaceclient.DefaultApi.UpdatePrivateSpaceAdvancedConfig(authctx, orgid, psid).PrivateSpaceAdvancedConfigBody(*defaults).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to reset Private Space Advanced Config",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func buildAdvancedConfigBody(d *schema.ResourceData) *private_space.PrivateSpaceAdvancedConfig {
	body := private_space.NewPrivateSpaceAdvancedConfig()
	enableIAM := d.Get("enable_iam_role").(bool)
	body.EnableIAMRole = &enableIAM

	if v, ok := d.GetOk("ingress_configuration"); ok {
		ingressList := v.([]interface{})
		if len(ingressList) > 0 {
			ingData := ingressList[0].(map[string]interface{})
			ingress := &private_space.PrivateSpaceIngressConfiguration{}
			ingress.ReadResponseTimeout = private_space.PtrString(ingData["read_response_timeout"].(string))
			ingress.Protocol = private_space.PtrString(ingData["protocol"].(string))

			logs := &private_space.PrivateSpaceIngressLogs{}
			logs.PortLogLevel = private_space.PtrString(ingData["port_log_level"].(string))

			if filters, ok := ingData["log_filters"].([]interface{}); ok && len(filters) > 0 {
				filterList := make([]private_space.PrivateSpaceIngressFilter, len(filters))
				for i, f := range filters {
					fMap := f.(map[string]interface{})
					filterList[i] = private_space.PrivateSpaceIngressFilter{
						IP:    fMap["ip"].(string),
						Level: fMap["level"].(string),
					}
				}
				logs.Filters = filterList
			}
			ingress.Logs = logs
			body.IngressConfiguration = ingress
		}
	}
	return body
}
