package anypoint

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mog "github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/managed_omni_gateway"
)

var regexpGatewayName = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,39}[a-z0-9])?)?$`)

func resourceManagedOmniGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceManagedOmniGatewayCreate,
		ReadContext:   resourceManagedOmniGatewayRead,
		UpdateContext: resourceManagedOmniGatewayUpdate,
		DeleteContext: resourceManagedOmniGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Creates and manages a CloudHub 2.0 Managed Omni Gateway instance.`,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of this gateway.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The organization id where the gateway is deployed.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The environment id where the gateway is deployed.",
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The name of the gateway.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexpGatewayName, "must match ^(?![-0-9])([a-z0-9-]{2,41}[a-z0-9])?$")),
			},
			"target_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target (private space) id to deploy this gateway to.",
			},
			"release_channel": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "lts",
				Description:      "The Omni Gateway release channel. Valid values: lts, edge.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"lts", "edge"}, false)),
			},
			"runtime_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Omni Gateway runtime version. Auto-selected by the platform if omitted.",
			},
			"size": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "small",
				Description:      "The gateway instance size. Valid values: small, large.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"small", "large"}, false)),
			},
			"desired_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The desired runtime status of the gateway (e.g. started, stopped).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current deployment status of the gateway.",
			},
			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the target (private space).",
			},
			"api_limit": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The API limit assigned to this gateway.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The gateway creation date.",
			},
			"last_updated_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the gateway was updated on the platform.",
			},
			"ingress": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Ingress configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public ingress URL.",
						},
						"internal_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The internal ingress URL.",
						},
						"forward_ssl_session": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether to forward SSL sessions.",
						},
						"last_mile_security": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether to enable TLS to the upstream (last mile security).",
						},
					},
				},
			},
			"properties": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Runtime properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upstream_response_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     15,
							Description: "Upstream response timeout in seconds.",
						},
						"connection_idle_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     60,
							Description: "Connection idle timeout in seconds.",
						},
					},
				},
			},
			"logging": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Logging configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "info",
							Description:      "Log level. Valid values: debug, info, warn, error.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"debug", "info", "warn", "error"}, false)),
						},
						"forward_logs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether to forward logs to Anypoint Monitoring.",
						},
					},
				},
			},
			"tracing": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Distributed tracing configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether distributed tracing is enabled.",
						},
					},
				},
			},
		},
	}
}

func resourceManagedOmniGatewayCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("organization_id").(string)
	envid := d.Get("environment_id").(string)
	authctx := getManagedOmniGatewayAuthCtx(ctx, &pco)

	body := newManagedOmniGatewayPostBody(d)
	res, httpr, err := pco.managedomnigatewayclient.DefaultApi.CreateGateway(authctx, orgid, envid).GatewayPostBody(*body).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Managed Omni Gateway",
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId(res.GetId())

	if ds, ok := d.GetOk("desired_status"); ok {
		if err := setManagedOmniGatewayDesiredStatus(ctx, &pco, orgid, envid, res.GetId(), ds.(string)); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to set desired status for Managed Omni Gateway " + res.GetId(),
				Detail:   err.Error(),
			})
			return diags
		}
	}

	return resourceManagedOmniGatewayRead(ctx, d, m)
}

func resourceManagedOmniGatewayRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("organization_id").(string)
	envid := d.Get("environment_id").(string)
	gwid := d.Id()
	authctx := getManagedOmniGatewayAuthCtx(ctx, &pco)

	res, httpr, err := pco.managedomnigatewayclient.DefaultApi.GetGateway(authctx, orgid, envid, gwid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Managed Omni Gateway " + gwid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	if err := setManagedOmniGatewayResourceData(d, res); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set Managed Omni Gateway " + gwid,
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(gwid)

	return diags
}

func resourceManagedOmniGatewayUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("organization_id").(string)
	envid := d.Get("environment_id").(string)
	gwid := d.Id()
	authctx := getManagedOmniGatewayAuthCtx(ctx, &pco)

	if d.HasChanges("release_channel", "runtime_version", "size", "ingress", "properties", "logging", "tracing") {
		body := newManagedOmniGatewayPostBody(d)
		_, httpr, err := pco.managedomnigatewayclient.DefaultApi.UpdateGateway(authctx, orgid, envid, gwid).GatewayPostBody(*body).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update Managed Omni Gateway " + gwid,
				Detail:   parseApiError(httpr, err),
			})
			return diags
		}
		defer httpr.Body.Close()
	}

	if d.HasChange("desired_status") {
		if ds, ok := d.GetOk("desired_status"); ok {
			if err := setManagedOmniGatewayDesiredStatus(ctx, &pco, orgid, envid, gwid, ds.(string)); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to set desired status for Managed Omni Gateway " + gwid,
					Detail:   err.Error(),
				})
				return diags
			}
		}
	}

	return resourceManagedOmniGatewayRead(ctx, d, m)
}

func resourceManagedOmniGatewayDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("organization_id").(string)
	envid := d.Get("environment_id").(string)
	gwid := d.Id()
	authctx := getManagedOmniGatewayAuthCtx(ctx, &pco)

	httpr, err := pco.managedomnigatewayclient.DefaultApi.DeleteGateway(authctx, orgid, envid, gwid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Managed Omni Gateway " + gwid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func setManagedOmniGatewayDesiredStatus(ctx context.Context, pco *ProviderConfOutput, orgid string, envid string, gwid string, desiredStatus string) error {
	authctx := getManagedOmniGatewayAuthCtx(ctx, pco)
	body := mog.NewGatewayDesiredStatusBody(desiredStatus)
	httpr, err := pco.managedomnigatewayclient.DefaultApi.UpdateGatewayDesiredStatus(authctx, orgid, envid, gwid).GatewayDesiredStatusBody(*body).Execute()
	if err != nil {
		return err
	}
	defer httpr.Body.Close()
	return nil
}

func newManagedOmniGatewayPostBody(d *schema.ResourceData) *mog.GatewayPostBody {
	body := mog.NewGatewayPostBody(
		d.Get("target_id").(string),
		d.Get("release_channel").(string),
		d.Get("runtime_version").(string),
		d.Get("size").(string),
	)
	body.SetName(d.Get("name").(string))
	body.SetConfiguration(*newManagedOmniGatewayConfiguration(d))
	return body
}

func newManagedOmniGatewayConfiguration(d *schema.ResourceData) *mog.GatewayConfiguration {
	config := mog.NewGatewayConfiguration()

	if v, ok := d.GetOk("ingress"); ok {
		list := v.([]any)
		if len(list) > 0 {
			m := list[0].(map[string]any)
			ingress := mog.NewGatewayIngress()
			ingress.SetForwardSslSession(m["forward_ssl_session"].(bool))
			ingress.SetLastMileSecurity(m["last_mile_security"].(bool))
			config.SetIngress(*ingress)
		}
	}
	if v, ok := d.GetOk("properties"); ok {
		list := v.([]any)
		if len(list) > 0 {
			m := list[0].(map[string]any)
			properties := mog.NewGatewayProperties()
			properties.SetUpstreamResponseTimeout(int32(m["upstream_response_timeout"].(int)))
			properties.SetConnectionIdleTimeout(int32(m["connection_idle_timeout"].(int)))
			config.SetProperties(*properties)
		}
	}
	if v, ok := d.GetOk("logging"); ok {
		list := v.([]any)
		if len(list) > 0 {
			m := list[0].(map[string]any)
			logging := mog.NewGatewayLogging()
			logging.SetLevel(m["level"].(string))
			logging.SetForwardLogs(m["forward_logs"].(bool))
			config.SetLogging(*logging)
		}
	}
	if v, ok := d.GetOk("tracing"); ok {
		list := v.([]any)
		if len(list) > 0 {
			m := list[0].(map[string]any)
			tracing := mog.NewGatewayTracing()
			tracing.SetEnabled(m["enabled"].(bool))
			config.SetTracing(*tracing)
		}
	}

	return config
}

func setManagedOmniGatewayResourceData(d *schema.ResourceData, gw *mog.Gateway) error {
	d.Set("name", gw.GetName())
	d.Set("target_id", gw.GetTargetId())
	d.Set("target_name", gw.GetTargetName())
	d.Set("status", gw.GetStatus())
	d.Set("desired_status", gw.GetDesiredStatus())
	d.Set("release_channel", gw.GetReleaseChannel())
	d.Set("runtime_version", gw.GetRuntimeVersion())
	d.Set("size", gw.GetSize())
	d.Set("api_limit", gw.GetApiLimit())
	d.Set("date_created", gw.GetDateCreated())
	d.Set("last_updated_platform", gw.GetLastUpdated())

	config := gw.GetConfiguration()
	ingress := config.GetIngress()
	d.Set("ingress", []map[string]any{
		{
			"public_url":          ingress.GetPublicUrl(),
			"internal_url":        ingress.GetInternalUrl(),
			"forward_ssl_session": ingress.GetForwardSslSession(),
			"last_mile_security":  ingress.GetLastMileSecurity(),
		},
	})
	properties := config.GetProperties()
	d.Set("properties", []map[string]any{
		{
			"upstream_response_timeout": int(properties.GetUpstreamResponseTimeout()),
			"connection_idle_timeout":   int(properties.GetConnectionIdleTimeout()),
		},
	})
	logging := config.GetLogging()
	d.Set("logging", []map[string]any{
		{
			"level":        logging.GetLevel(),
			"forward_logs": logging.GetForwardLogs(),
		},
	})
	tracing := config.GetTracing()
	d.Set("tracing", []map[string]any{
		{
			"enabled": tracing.GetEnabled(),
		},
	})

	return nil
}

func getManagedOmniGatewayAuthCtx(ctx context.Context, pco *ProviderConfOutput) context.Context {
	tmp := context.WithValue(ctx, mog.ContextAccessToken, pco.access_token)
	return context.WithValue(tmp, mog.ContextServerIndex, pco.server_index)
}
