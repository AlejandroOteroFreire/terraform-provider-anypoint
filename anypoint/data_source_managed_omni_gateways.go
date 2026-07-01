package anypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mog "github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/managed_omni_gateway"
)

func dataSourceManagedOmniGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceManagedOmniGatewaysRead,
		Description: `
		Lists all ` + "`" + `managed omni gateway` + "`" + ` instances in an environment.
		`,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization id where the gateways are deployed.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment id where the gateways are deployed.",
			},
			"gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of managed omni gateway instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The current deployment status of the gateway.",
						},
						"runtime_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The active Omni Gateway runtime version.",
						},
						"release_channel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Omni Gateway release channel in use.",
						},
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The gateway instance size.",
						},
					},
				},
			},
		},
	}
}

func dataSourceManagedOmniGatewaysRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("organization_id").(string)
	envid := d.Get("environment_id").(string)
	authctx := getManagedOmniGatewayAuthCtx(ctx, &pco)

	res, httpr, err := pco.managedomnigatewayclient.DefaultApi.GetGateways(authctx, orgid, envid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to list Managed Omni Gateways",
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.Set("gateways", flattenManagedOmniGatewaysList(res.GetContent()))
	d.SetId(orgid + "/" + envid)

	return diags
}

func flattenManagedOmniGatewaysList(items []mog.Gateway) []map[string]any {
	result := make([]map[string]any, len(items))
	for i, gw := range items {
		result[i] = map[string]any{
			"id":              gw.GetId(),
			"name":            gw.GetName(),
			"status":          gw.GetStatus(),
			"runtime_version": gw.GetRuntimeVersion(),
			"release_channel": gw.GetReleaseChannel(),
			"size":            gw.GetSize(),
		}
	}
	return result
}
