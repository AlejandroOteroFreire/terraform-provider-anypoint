package anypoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/private_space"
)

func resourcePrivateSpaceTransitGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateSpaceTransitGatewayCreate,
		ReadContext:   resourcePrivateSpaceTransitGatewayRead,
		UpdateContext: resourcePrivateSpaceTransitGatewayUpdate,
		DeleteContext: resourcePrivateSpaceTransitGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages a Transit Gateway attachment within a Private Space.`,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
				Description: "The unique identifier of the private space.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the transit gateway.",
			},
			"resource_share_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS Resource Share ID (from AWS Resource Access Manager).",
			},
			"resource_share_account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS account ID that owns the transit gateway.",
			},
			"routes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of CIDR routes to propagate through this transit gateway.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// Read-only status fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the transit gateway attachment.",
			},
			"attachment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The attachment status.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AWS region of the transit gateway.",
			},
		},
	}
}

func resourcePrivateSpaceTransitGatewayCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	body := buildTransitGatewayBody(d)

	res, httpr, err := pco.privatespaceclient.DefaultApi.CreatePrivateSpaceTransitGateway(authctx, orgid, psid).PrivateSpaceTransitGatewayPostBody(*body).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Transit Gateway",
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId(res.GetId())
	return resourcePrivateSpaceTransitGatewayRead(ctx, d, m)
}

func resourcePrivateSpaceTransitGatewayRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	tgwid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	if isComposedResourceId(tgwid) {
		var decomposeDiags diag.Diagnostics
		orgid, psid, tgwid, decomposeDiags = decomposePrivateSpaceTransitGatewayId(d)
		if decomposeDiags.HasError() {
			return decomposeDiags
		}
	}

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceTransitGateway(authctx, orgid, psid, tgwid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Transit Gateway " + tgwid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.Set("name", res.GetName())
	d.Set("resource_share_id", res.GetResourceShareId())
	d.Set("resource_share_account", res.GetResourceShareAccount())
	d.Set("routes", res.GetRoutes())
	d.Set("status", res.GetStatus())
	d.Set("attachment", res.GetAttachment())
	d.Set("region", res.GetRegion())
	d.Set("org_id", orgid)
	d.Set("private_space_id", psid)
	d.SetId(res.GetId())

	return diags
}

func resourcePrivateSpaceTransitGatewayUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	tgwid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	body := buildTransitGatewayBody(d)

	_, httpr, err := pco.privatespaceclient.DefaultApi.UpdatePrivateSpaceTransitGateway(authctx, orgid, psid, tgwid).PrivateSpaceTransitGatewayPatchBody(*body).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Transit Gateway " + tgwid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	return resourcePrivateSpaceTransitGatewayRead(ctx, d, m)
}

func resourcePrivateSpaceTransitGatewayDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	tgwid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	httpr, err := pco.privatespaceclient.DefaultApi.DeletePrivateSpaceTransitGateway(authctx, orgid, psid, tgwid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Transit Gateway " + tgwid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func buildTransitGatewayBody(d *schema.ResourceData) *private_space.PrivateSpaceTransitGateway {
	body := private_space.NewPrivateSpaceTransitGateway()
	body.Name = private_space.PtrString(d.Get("name").(string))
	body.ResourceShareId = private_space.PtrString(d.Get("resource_share_id").(string))
	body.ResourceShareAccount = private_space.PtrString(d.Get("resource_share_account").(string))

	if routes, ok := d.GetOk("routes"); ok {
		rawRoutes := routes.([]interface{})
		routeList := make([]string, len(rawRoutes))
		for i, r := range rawRoutes {
			routeList[i] = r.(string)
		}
		body.Routes = routeList
	}

	return body
}

func decomposePrivateSpaceTransitGatewayId(d *schema.ResourceData) (string, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	s := DecomposeResourceId(d.Id())
	if len(s) != 3 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Transit Gateway ID format",
			Detail:   fmt.Sprintf("Expected ORG_ID/PRIVATE_SPACE_ID/TGW_ID, got %s", d.Id()),
		})
		return "", "", "", diags
	}
	return s[0], s[1], s[2], diags
}
