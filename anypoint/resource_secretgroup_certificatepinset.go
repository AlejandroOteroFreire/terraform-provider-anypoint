package anypoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/anypoint-client-go/secretgroup_certificatepinset"
)

func resourceSecretGroupCertificatePinset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretGroupCertificatePinsetCreate,
		ReadContext:   resourceSecretGroupCertificatePinsetRead,
		DeleteContext: resourceSecretGroupCertificatePinsetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages a Certificate Pinset within a Secret Group.`,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_pinset_base64": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Base64-encoded PEM content of the certificate pinset.",
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecretGroupCertificatePinsetCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	sgid := d.Get("secret_group_id").(string)
	authctx := getSgCertificatePinsetAuthCtx(ctx, &pco)

	res, httpr, err := pco.sgcertificatepinsetclient.DefaultApi.CreateCertificatePinset(authctx, orgid, envid, sgid).
		Name(d.Get("name").(string)).
		PinsetBase64(d.Get("certificate_pinset_base64").(string)).
		Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Certificate Pinset",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId(res.GetId())
	return resourceSecretGroupCertificatePinsetRead(ctx, d, m)
}

func resourceSecretGroupCertificatePinsetRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	authctx := getSgCertificatePinsetAuthCtx(ctx, &pco)

	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	sgid := d.Get("secret_group_id").(string)
	pinId := d.Id()

	if isComposedResourceId(pinId) {
		var decomposeDiags diag.Diagnostics
		orgid, envid, sgid, pinId, decomposeDiags = decomposeCertificatePinsetId(d)
		if decomposeDiags.HasError() {
			return decomposeDiags
		}
	}

	res, httpr, err := pco.sgcertificatepinsetclient.DefaultApi.GetCertificatePinset(authctx, orgid, envid, sgid, pinId).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Certificate Pinset",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.Set("org_id", orgid)
	d.Set("env_id", envid)
	d.Set("secret_group_id", sgid)
	d.Set("name", res.GetName())
	d.Set("expiration_date", res.GetExpirationDate())
	d.Set("algorithm", res.GetAlgorithm())
	d.SetId(res.GetId())
	return diags
}

func resourceSecretGroupCertificatePinsetDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	sgid := d.Get("secret_group_id").(string)
	pinId := d.Id()
	authctx := getSgCertificatePinsetAuthCtx(ctx, &pco)

	httpr, err := pco.sgcertificatepinsetclient.DefaultApi.DeleteCertificatePinset(authctx, orgid, envid, sgid, pinId).Execute()
	if err != nil {
		// Pinsets get cleaned up with parent secret group, ignore 404
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Certificate Pinset",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func decomposeCertificatePinsetId(d *schema.ResourceData) (string, string, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	s := DecomposeResourceId(d.Id())
	if len(s) != 4 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Certificate Pinset ID format",
			Detail:   fmt.Sprintf("Expected ORG_ID/ENV_ID/SECRET_GROUP_ID/PIN_ID, got %s", d.Id()),
		})
		return "", "", "", "", diags
	}
	return s[0], s[1], s[2], s[3], diags
}

func getSgCertificatePinsetAuthCtx(ctx context.Context, pco *ProviderConfOutput) context.Context {
	tmp := context.WithValue(ctx, secretgroup_certificatepinset.ContextAccessToken, pco.access_token)
	return context.WithValue(tmp, secretgroup_certificatepinset.ContextServerIndex, pco.server_index)
}
