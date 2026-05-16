package anypoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/private_space"
)

func resourcePrivateSpaceAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateSpaceAssociationCreate,
		ReadContext:   resourcePrivateSpaceAssociationRead,
		UpdateContext: resourcePrivateSpaceAssociationUpdate,
		DeleteContext: resourcePrivateSpaceAssociationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages environment associations for a Private Space. Controls which organizations and environments can use the private space.`,
		Schema: map[string]*schema.Schema{
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
			"associations": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of environment associations to create.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The organization ID to associate.",
						},
						"environment": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Environment to associate: a UUID, 'all', 'production', or 'sandbox'.",
						},
					},
				},
			},
			"created_associations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The associations as created by the API, including their assigned IDs.",
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

func resourcePrivateSpaceAssociationCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Get("private_space_id").(string)
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	body := buildAssociationBody(d)
	res, httpr, err := pco.privatespaceclient.DefaultApi.CreatePrivateSpaceAssociations(authctx, orgid, psid).PrivateSpaceAssociationBody(*body).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Private Space Associations",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId(psid)
	d.Set("created_associations", flattenAssociations(res))
	return diags
}

func resourcePrivateSpaceAssociationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	res, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceAssociations(authctx, orgid, psid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Private Space Associations",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.Set("private_space_id", psid)
	d.Set("created_associations", flattenAssociations(res))
	return diags
}

func resourcePrivateSpaceAssociationUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	// Get current associations to find IDs to delete
	current, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceAssociations(authctx, orgid, psid).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read current associations for update",
			Detail:   parseApiError(httpr, err),
		})
	}
	httpr.Body.Close()

	// Delete all existing associations
	for _, assoc := range current {
		delhttpr, err := pco.privatespaceclient.DefaultApi.DeletePrivateSpaceAssociation(authctx, orgid, psid, assoc.GetId()).Execute()
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to delete association %s during update", assoc.GetId()),
				Detail:   parseApiError(delhttpr, err),
			})
		}
		delhttpr.Body.Close()
	}

	// Recreate with new values
	return resourcePrivateSpaceAssociationCreate(ctx, d, m)
}

func resourcePrivateSpaceAssociationDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	psid := d.Id()
	authctx := getPrivateSpaceAuthCtx(ctx, &pco)

	current, httpr, err := pco.privatespaceclient.DefaultApi.GetPrivateSpaceAssociations(authctx, orgid, psid).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read associations for deletion",
			Detail:   parseApiError(httpr, err),
		})
	}
	httpr.Body.Close()

	for _, assoc := range current {
		delhttpr, err := pco.privatespaceclient.DefaultApi.DeletePrivateSpaceAssociation(authctx, orgid, psid, assoc.GetId()).Execute()
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unable to delete association %s", assoc.GetId()),
				Detail:   parseApiError(delhttpr, err),
			})
		}
		delhttpr.Body.Close()
	}

	d.SetId("")
	return diags
}

func buildAssociationBody(d *schema.ResourceData) *private_space.PrivateSpaceAssociationRequest {
	body := private_space.NewPrivateSpaceAssociationRequest()
	raw := d.Get("associations").([]interface{})
	items := make([]private_space.PrivateSpaceAssociationItem, len(raw))
	for i, v := range raw {
		a := v.(map[string]interface{})
		items[i] = private_space.PrivateSpaceAssociationItem{
			OrganizationId: a["organization_id"].(string),
			Environment:    a["environment"].(string),
		}
	}
	body.Associations = items
	return body
}

func flattenAssociations(assocs []private_space.PrivateSpaceAssociation) []interface{} {
	result := make([]interface{}, len(assocs))
	for i, a := range assocs {
		result[i] = map[string]interface{}{
			"id":              a.GetId(),
			"organization_id": a.GetOrganizationId(),
			"environment_id":  a.GetEnvironmentId(),
		}
	}
	return result
}
