package anypoint

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivateSpaceTransitGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateSpaceTransitGatewaysRead,
		Description: `Reads all Transit Gateways attached to a Private Space.`,
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
			"transit_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":                     {Type: schema.TypeString, Computed: true},
						"name":                   {Type: schema.TypeString, Computed: true},
						"resource_share_id":      {Type: schema.TypeString, Computed: true},
						"resource_share_account": {Type: schema.TypeString, Computed: true},
						"routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status":     {Type: schema.TypeString, Computed: true},
						"attachment": {Type: schema.TypeString, Computed: true},
						"region":     {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourcePrivateSpaceTransitGatewaysRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceTransitGateways(authctx, orgid, psid).Execute()
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
			Summary:  "Unable to list Transit Gateways for private space " + psid,
			Detail:   details,
		})
	}
	defer httpr.Body.Close()

	tgws := make([]interface{}, len(res))
	for i, item := range res {
		tgws[i] = flattenTransitGatewayItem(&item)
	}

	if err := d.Set("transit_gateways", tgws); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
