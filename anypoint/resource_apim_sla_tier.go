package anypoint

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/apim"
)

func resourceApimSlaTier() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApimSlaTierCreate,
		ReadContext:   resourceApimSlaTierRead,
		UpdateContext: resourceApimSlaTierUpdate,
		DeleteContext: resourceApimSlaTierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages an SLA Tier for an API Manager instance.`,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_approve": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "ACTIVE",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false)),
			},
			"limits": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_period_in_milliseconds": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"maximum_requests": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"visible": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
	}
}

func resourceApimSlaTierCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	apiid := d.Get("api_instance_id").(string)
	authctx := getApimAuthCtx(ctx, &pco)

	body := buildSlaTierBody(d)
	res, httpr, err := pco.apimclient.DefaultApi.CreateSlaTier(authctx, orgid, envid, apiid).SlaTierBody(*body).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create SLA Tier",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId(strconv.FormatInt(res.GetId(), 10))
	return resourceApimSlaTierRead(ctx, d, m)
}

func resourceApimSlaTierRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	authctx := getApimAuthCtx(ctx, &pco)

	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	apiid := d.Get("api_instance_id").(string)
	tierId := d.Id()

	if isComposedResourceId(tierId) {
		var decomposeDiags diag.Diagnostics
		orgid, envid, apiid, tierId, decomposeDiags = decomposeApimSlaTierId(d)
		if decomposeDiags.HasError() {
			return decomposeDiags
		}
	}

	res, httpr, err := pco.apimclient.DefaultApi.GetSlaTiers(authctx, orgid, envid, apiid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read SLA Tiers",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	tierIdInt, _ := strconv.ParseInt(tierId, 10, 64)
	var found *apim.SlaTier
	for i, t := range res.Tiers {
		if t.GetId() == tierIdInt {
			found = &res.Tiers[i]
			break
		}
	}
	if found == nil {
		d.SetId("")
		return diags
	}

	d.Set("name", found.GetName())
	d.Set("description", found.GetDescription())
	d.Set("auto_approve", found.GetAutoApprove())
	d.Set("status", found.GetStatus())
	d.Set("limits", flattenSlaTierLimits(found.GetLimits()))
	d.Set("org_id", orgid)
	d.Set("env_id", envid)
	d.Set("api_instance_id", apiid)
	d.SetId(strconv.FormatInt(found.GetId(), 10))
	return diags
}

func resourceApimSlaTierUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	apiid := d.Get("api_instance_id").(string)
	tierId := d.Id()
	authctx := getApimAuthCtx(ctx, &pco)

	body := buildSlaTierBody(d)
	_, httpr, err := pco.apimclient.DefaultApi.UpdateSlaTier(authctx, orgid, envid, apiid, tierId).SlaTierBody(*body).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update SLA Tier",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	return resourceApimSlaTierRead(ctx, d, m)
}

func resourceApimSlaTierDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	apiid := d.Get("api_instance_id").(string)
	tierId := d.Id()
	authctx := getApimAuthCtx(ctx, &pco)

	httpr, err := pco.apimclient.DefaultApi.DeleteSlaTier(authctx, orgid, envid, apiid, tierId).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete SLA Tier",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func buildSlaTierBody(d *schema.ResourceData) *apim.SlaTier {
	body := apim.NewSlaTier()
	body.Name = apim.PtrString(d.Get("name").(string))
	autoApprove := d.Get("auto_approve").(bool)
	body.AutoApprove = &autoApprove
	status := d.Get("status").(string)
	body.Status = &status
	if v, ok := d.GetOk("description"); ok {
		body.Description = apim.PtrString(v.(string))
	}

	rawLimits := d.Get("limits").([]interface{})
	limits := make([]apim.SlaTierLimit, len(rawLimits))
	for i, l := range rawLimits {
		lm := l.(map[string]interface{})
		tp := int64(lm["time_period_in_milliseconds"].(int))
		mr := int64(lm["maximum_requests"].(int))
		vis := lm["visible"].(bool)
		limits[i] = apim.SlaTierLimit{
			TimePeriodInMilliseconds: &tp,
			MaximumRequests:          &mr,
			Visible:                  &vis,
		}
	}
	body.Limits = limits
	return body
}

func flattenSlaTierLimits(limits []apim.SlaTierLimit) []interface{} {
	result := make([]interface{}, len(limits))
	for i, l := range limits {
		result[i] = map[string]interface{}{
			"time_period_in_milliseconds": int(l.GetTimePeriodInMilliseconds()),
			"maximum_requests":            int(l.GetMaximumRequests()),
			"visible":                     l.GetVisible(),
		}
	}
	return result
}

func decomposeApimSlaTierId(d *schema.ResourceData) (string, string, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	s := DecomposeResourceId(d.Id())
	if len(s) != 4 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid SLA Tier ID format",
			Detail:   fmt.Sprintf("Expected ORG_ID/ENV_ID/API_ID/TIER_ID, got %s", d.Id()),
		})
		return "", "", "", "", diags
	}
	return s[0], s[1], s[2], s[3], diags
}
