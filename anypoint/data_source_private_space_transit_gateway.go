package anypoint

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/anypoint-client-go/private_space"
)

func dataSourcePrivateSpaceTransitGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateSpaceTransitGatewayRead,
		Description: `Reads a single Transit Gateway attached to a Private Space.`,
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
			"tgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the transit gateway.",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_share_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePrivateSpaceTransitGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgId := d.Get("org_id").(string)
	psId := d.Get("private_space_id").(string)
	tgwId := d.Get("tgw_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceTransitGateway(authctx, orgId, psId, tgwId).Execute()
	if err != nil {
		var details string
		if httpr != nil && httpr.StatusCode >= 400 {
			defer httpr.Body.Close()
			b, _ := io.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to read Transit Gateway %s", tgwId),
			Detail:   details,
		})
	}
	defer httpr.Body.Close()

	d.SetId(res.GetId())
	d.Set("name", res.GetName())
	d.Set("resource_share_id", res.GetResourceShareId())
	d.Set("resource_share_account", res.GetResourceShareAccount())
	d.Set("routes", res.GetRoutes())
	d.Set("status", res.GetStatus())
	d.Set("attachment", res.GetAttachment())
	d.Set("region", res.GetRegion())

	return diags
}

func flattenTransitGatewayItem(item *private_space.PrivateSpaceTransitGateway) map[string]interface{} {
	return map[string]interface{}{
		"id":                     item.GetId(),
		"name":                   item.GetName(),
		"resource_share_id":      item.GetResourceShareId(),
		"resource_share_account": item.GetResourceShareAccount(),
		"routes":                 item.GetRoutes(),
		"status":                 item.GetStatus(),
		"attachment":             item.GetAttachment(),
		"region":                 item.GetRegion(),
	}
}
