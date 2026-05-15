package anypoint

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSelfManagedOmniGatewayRegistrationToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSelfManagedOmniGatewayRegistrationTokenRead,
		Description: `
		Retrieve a Self-Managed Omni Gateway registration token used to register a new Self-Managed Omni Gateway instance.
		`,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization id where the Self-Managed Omni Gateway targets are defined.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment id where the Self-Managed Omni Gateway targets are defined.",
			},
			"registration_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration token that can be used to register a new Self-Managed Omni Gateway",
			},
		},
	}
}

func dataSourceSelfManagedOmniGatewayRegistrationTokenRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	authctx := getSelfManagedOmniGatewayAuthCtx(ctx, &pco)
	//perform request
	res, httpr, err := pco.flexgatewayclient.DefaultApi.GetFlexGatewayRegistrationToken(authctx, orgid, envid).Execute()
	if err != nil {
		var details string
		if httpr != nil && httpr.StatusCode >= 400 {
			defer httpr.Body.Close()
			b, _ := io.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Self-Managed Omni Gateway registration token ",
			Detail:   details,
		})
		return diags
	}
	defer httpr.Body.Close()
	//process response data
	if err := d.Set("registration_token", res.GetRegistrationToken()); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set registration token of fex-gateway org " + orgid + " and env " + envid,
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
