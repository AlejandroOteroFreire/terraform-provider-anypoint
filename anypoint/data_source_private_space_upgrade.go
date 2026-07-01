package anypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivateSpaceUpgrade() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateSpaceUpgradeRead,
		Description: `
		Reads the ` + "`" + `upgrade status` + "`" + ` of a private space.
		`,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This resource's unique id, equal to the private space id.",
			},
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
			"scheduled_update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the upgrade is scheduled for.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The upgrade status (e.g. QUEUED).",
			},
		},
	}
}

func dataSourcePrivateSpaceUpgradeRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetUpgradeStatus(authctx, orgid, psid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Private Space upgrade status " + psid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.Set("scheduled_update_time", res.GetScheduledUpdateTime())
	d.Set("status", res.GetStatus())
	d.SetId(psid)

	return diags
}
