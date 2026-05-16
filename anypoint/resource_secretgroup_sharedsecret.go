package anypoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/secretgroup_sharedsecret"
)

func resourceSecretGroupSharedSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretGroupSharedSecretCreate,
		ReadContext:   resourceSecretGroupSharedSecretRead,
		UpdateContext: resourceSecretGroupSharedSecretUpdate,
		DeleteContext: resourceSecretGroupSharedSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages a Shared Secret within a Secret Group.`,
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
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"UsernamePassword", "S3Credential", "SymmetricKey", "Blob"}, false)),
				Description:      "Type of shared secret. One of: UsernamePassword, S3Credential, SymmetricKey, Blob.",
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username (for UsernamePassword type).",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password (for UsernamePassword type).",
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS Access Key ID (for S3Credential type).",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS Secret Access Key (for S3Credential type).",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Base64-encoded key (for SymmetricKey type).",
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Content (for Blob type).",
			},
		},
	}
}

func resourceSecretGroupSharedSecretCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	sgid := d.Get("secret_group_id").(string)
	authctx := getSgSharedSecretAuthCtx(ctx, &pco)

	body := buildSharedSecretBody(d)
	res, httpr, err := pco.sgsharedsecretclient.DefaultApi.CreateSharedSecret(authctx, orgid, envid, sgid).SharedSecretBody(*body).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Shared Secret",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId(res.GetId())
	return resourceSecretGroupSharedSecretRead(ctx, d, m)
}

func resourceSecretGroupSharedSecretRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	authctx := getSgSharedSecretAuthCtx(ctx, &pco)

	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	sgid := d.Get("secret_group_id").(string)
	ssid := d.Id()

	if isComposedResourceId(ssid) {
		var decomposeDiags diag.Diagnostics
		orgid, envid, sgid, ssid, decomposeDiags = decomposeSharedSecretId(d)
		if decomposeDiags.HasError() {
			return decomposeDiags
		}
	}

	res, httpr, err := pco.sgsharedsecretclient.DefaultApi.GetSharedSecret(authctx, orgid, envid, sgid, ssid).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Shared Secret",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.Set("org_id", orgid)
	d.Set("env_id", envid)
	d.Set("secret_group_id", sgid)
	d.Set("name", res.GetName())
	d.Set("type", res.GetType())
	d.Set("expiration_date", res.GetExpirationDate())
	d.Set("username", res.GetUsername())
	d.Set("access_key_id", res.GetAccessKeyId())
	d.SetId(res.GetId())
	return diags
}

func resourceSecretGroupSharedSecretUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	sgid := d.Get("secret_group_id").(string)
	ssid := d.Id()
	authctx := getSgSharedSecretAuthCtx(ctx, &pco)

	body := buildSharedSecretBody(d)
	_, httpr, err := pco.sgsharedsecretclient.DefaultApi.UpdateSharedSecret(authctx, orgid, envid, sgid, ssid).SharedSecretBody(*body).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Shared Secret",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	return resourceSecretGroupSharedSecretRead(ctx, d, m)
}

func resourceSecretGroupSharedSecretDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	sgid := d.Get("secret_group_id").(string)
	ssid := d.Id()
	authctx := getSgSharedSecretAuthCtx(ctx, &pco)

	httpr, err := pco.sgsharedsecretclient.DefaultApi.DeleteSharedSecret(authctx, orgid, envid, sgid, ssid).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete Shared Secret",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func buildSharedSecretBody(d *schema.ResourceData) *secretgroup_sharedsecret.SharedSecret {
	body := secretgroup_sharedsecret.NewSharedSecret()
	name := d.Get("name").(string)
	body.Name = &name
	t := d.Get("type").(string)
	body.Type = &t

	if v, ok := d.GetOk("expiration_date"); ok {
		s := v.(string)
		body.ExpirationDate = &s
	}
	if v, ok := d.GetOk("username"); ok {
		s := v.(string)
		body.Username = &s
	}
	if v, ok := d.GetOk("password"); ok {
		s := v.(string)
		body.Password = &s
	}
	if v, ok := d.GetOk("access_key_id"); ok {
		s := v.(string)
		body.AccessKeyId = &s
	}
	if v, ok := d.GetOk("secret_access_key"); ok {
		s := v.(string)
		body.SecretAccessKey = &s
	}
	if v, ok := d.GetOk("key"); ok {
		s := v.(string)
		body.Key = &s
	}
	if v, ok := d.GetOk("content"); ok {
		s := v.(string)
		body.Content = &s
	}
	return body
}

func decomposeSharedSecretId(d *schema.ResourceData) (string, string, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	s := DecomposeResourceId(d.Id())
	if len(s) != 4 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid Shared Secret ID format",
			Detail:   fmt.Sprintf("Expected ORG_ID/ENV_ID/SECRET_GROUP_ID/SECRET_ID, got %s", d.Id()),
		})
		return "", "", "", "", diags
	}
	return s[0], s[1], s[2], s[3], diags
}

func getSgSharedSecretAuthCtx(ctx context.Context, pco *ProviderConfOutput) context.Context {
	tmp := context.WithValue(ctx, secretgroup_sharedsecret.ContextAccessToken, pco.access_token)
	return context.WithValue(tmp, secretgroup_sharedsecret.ContextServerIndex, pco.server_index)
}
