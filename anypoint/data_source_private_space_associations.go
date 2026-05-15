package anypoint

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivateSpaceAssociations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePrivateSpaceAssociationsRead,
		Description: `Reads all environment associations for a Private Space.`,
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
			"associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":              {Type: schema.TypeString, Computed: true},
						"organization_id": {Type: schema.TypeString, Computed: true},
						"environment_id":  {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourcePrivateSpaceAssociationsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceAssociations(authctx, orgid, psid).Execute()
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
			Summary:  "Unable to list associations for private space " + psid,
			Detail:   details,
		})
	}
	defer httpr.Body.Close()

	d.Set("associations", flattenAssociations(res))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
