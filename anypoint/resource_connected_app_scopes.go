package anypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/connected_app"
)

func resourceConnectedAppScopes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectedAppScopesUpdate,
		ReadContext:   resourceConnectedAppScopesRead,
		UpdateContext: resourceConnectedAppScopesUpdate,
		DeleteContext: resourceConnectedAppScopesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages scopes for a Connected App.`,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The organization id where the connected app is defined.",
			},
			"connected_app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the connected app.",
			},
			"scopes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of scopes to assign to the connected app.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The scope name (e.g. 'read:full', 'manage:exchange').",
						},
						"context_params": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Context parameters for the scope (org/env).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"org": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"env_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceConnectedAppScopesUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgId := d.Get("org_id").(string)
	connAppId := d.Get("connected_app_id").(string)
	authctx := getConnectedAppAuthCtx(ctx, &pco)

	body := buildConnectedAppScopesBody(d)
	httpr, err := pco.connectedappclient.DefaultApi.UpdateConnectedAppScopes(authctx, orgId, connAppId).ConnectedAppScopesPutBody(*body).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update Connected App Scopes",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId(connAppId)
	return resourceConnectedAppScopesRead(ctx, d, m)
}

func resourceConnectedAppScopesRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgId := d.Get("org_id").(string)
	connAppId := d.Id()
	authctx := getConnectedAppAuthCtx(ctx, &pco)

	res, httpr, err := pco.connectedappclient.DefaultApi.GetConnectedAppScopes(authctx, orgId, connAppId).Execute()
	if err != nil {
		if httpr != nil && httpr.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read Connected App Scopes",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.Set("scopes", flattenConnectedAppScopes(res.GetData()))
	d.Set("connected_app_id", connAppId)
	d.SetId(connAppId)
	return diags
}

func resourceConnectedAppScopesDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgId := d.Get("org_id").(string)
	connAppId := d.Id()
	authctx := getConnectedAppAuthCtx(ctx, &pco)

	empty := connected_app.NewConnectedAppScopesPutBody()
	empty.Scopes = []connected_app.ScopeCore{}

	httpr, err := pco.connectedappclient.DefaultApi.UpdateConnectedAppScopes(authctx, orgId, connAppId).ConnectedAppScopesPutBody(*empty).Execute()
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to clear Connected App Scopes",
			Detail:   parseApiError(httpr, err),
		})
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}

func buildConnectedAppScopesBody(d *schema.ResourceData) *connected_app.ConnectedAppScopesPutBody {
	body := connected_app.NewConnectedAppScopesPutBody()
	raw := d.Get("scopes").([]interface{})
	scopes := make([]connected_app.ScopeCore, len(raw))
	for i, v := range raw {
		s := v.(map[string]interface{})
		scope := connected_app.NewScopeCore()
		scopeName := s["scope"].(string)
		scope.Scope = &scopeName

		if cpList, ok := s["context_params"].([]interface{}); ok && len(cpList) > 0 {
			cp := cpList[0].(map[string]interface{})
			cps := connected_app.NewContextParams()
			if org, ok := cp["org"].(string); ok && org != "" {
				cps.Org = &org
			}
			if env, ok := cp["env_id"].(string); ok && env != "" {
				cps.EnvId = &env
			}
			scope.ContextParams = cps
		}
		scopes[i] = *scope
	}
	body.Scopes = scopes
	return body
}

func flattenConnectedAppScopes(scopes []connected_app.ScopeCore) []interface{} {
	result := make([]interface{}, len(scopes))
	for i, s := range scopes {
		m := map[string]interface{}{
			"scope": s.GetScope(),
		}
		cp := s.GetContextParams()
		if cp.Org != nil || cp.EnvId != nil {
			m["context_params"] = []interface{}{
				map[string]interface{}{
					"org":    cp.GetOrg(),
					"env_id": cp.GetEnvId(),
				},
			}
		}
		result[i] = m
	}
	return result
}
