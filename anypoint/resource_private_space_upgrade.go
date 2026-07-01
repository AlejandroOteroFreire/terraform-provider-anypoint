package anypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePrivateSpaceUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateSpaceUpgradeCreate,
		ReadContext:   resourcePrivateSpaceUpgradeRead,
		UpdateContext: resourcePrivateSpaceUpgradeUpdate,
		DeleteContext: resourcePrivateSpaceUpgradeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Schedules (or cancels) a runtime upgrade for a Private Space.`,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This resource's unique id, equal to the private space id.",
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
				Description: "The unique identifier of the private space to upgrade.",
			},
			"opt_in": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to opt in to the scheduled upgrade.",
			},
			"date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The scheduled upgrade date. If omitted, the platform default schedule is used.",
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

func resourcePrivateSpaceUpgradeCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	req := pco.privatespaceclient.DefaultApi.ScheduleUpgrade(authctx, orgid, psid).OptIn(d.Get("opt_in").(bool))
	if date, ok := d.GetOk("date"); ok {
		req = req.Date(date.(string))
	}

	_, httpr, err := req.Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to schedule Private Space upgrade",
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId(psid)
	return resourcePrivateSpaceUpgradeRead(ctx, d, m)
}

func resourcePrivateSpaceUpgradeRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetUpgradeStatus(authctx, orgid, psid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Private Space upgrade status " + psid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.Set("private_space_id", psid)
	d.Set("scheduled_update_time", res.GetScheduledUpdateTime())
	d.Set("status", res.GetStatus())
	d.SetId(psid)

	return diags
}

func resourcePrivateSpaceUpgradeUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	req := pco.privatespaceclient.DefaultApi.ScheduleUpgrade(authctx, orgid, psid).OptIn(d.Get("opt_in").(bool))
	if date, ok := d.GetOk("date"); ok {
		req = req.Date(date.(string))
	}

	_, httpr, err := req.Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Private Space upgrade " + psid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	return resourcePrivateSpaceUpgradeRead(ctx, d, m)
}

func resourcePrivateSpaceUpgradeDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	httpr, err := pco.privatespaceclient.DefaultApi.CancelUpgrade(authctx, orgid, psid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to cancel Private Space upgrade " + psid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}
