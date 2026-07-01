package anypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceManagedOmniGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceManagedOmniGatewayRead,
		Description: `
		Reads a specific ` + "`" + `managed omni gateway` + "`" + ` instance.
		`,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization id where the gateway is deployed.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment id where the gateway is deployed.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The gateway id to retrieve.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of this gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the gateway.",
			},
			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target (private space) id this gateway is deployed to.",
			},
			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the target (private space).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current deployment status of the gateway.",
			},
			"desired_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The desired runtime status of the gateway.",
			},
			"release_channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Omni Gateway release channel in use.",
			},
			"runtime_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The active Omni Gateway runtime version.",
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The gateway instance size.",
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
		},
	}
}

func dataSourceManagedOmniGatewayRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("organization_id").(string)
	envid := d.Get("environment_id").(string)
	gwid := d.Get("gateway_id").(string)
	authctx := getManagedOmniGatewayAuthCtx(ctx, &pco)

	res, httpr, err := pco.managedomnigatewayclient.DefaultApi.GetGateway(authctx, orgid, envid, gwid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Managed Omni Gateway " + gwid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.Set("name", res.GetName())
	d.Set("target_id", res.GetTargetId())
	d.Set("target_name", res.GetTargetName())
	d.Set("status", res.GetStatus())
	d.Set("desired_status", res.GetDesiredStatus())
	d.Set("release_channel", res.GetReleaseChannel())
	d.Set("runtime_version", res.GetRuntimeVersion())
	d.Set("size", res.GetSize())
	d.Set("api_limit", res.GetApiLimit())
	d.Set("date_created", res.GetDateCreated())
	d.Set("last_updated_platform", res.GetLastUpdated())
	d.SetId(res.GetId())

	return diags
}
