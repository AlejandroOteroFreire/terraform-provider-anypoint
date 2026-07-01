package anypoint

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const mcpServerEndpointType = "mcp"

func resourceMCPServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMCPServerCreate,
		ReadContext:   resourceMCPServerRead,
		UpdateContext: resourceMCPServerUpdate,
		DeleteContext: resourceMCPServerDelete,
		Description: `
		Manages an MCP (Model Context Protocol) server instance in Anypoint API Manager,
		deployed to an Omni Gateway target with routing rules and upstream backends.
		`,
		Schema: agentsToolsSchema(mcpServerEndpointType, "Endpoint protocol type. Always 'mcp' for MCP server instances."),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMCPServerCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	authctx := getApimAuthCtx(ctx, &pco)

	if _, hasDeployment := d.GetOk("deployment"); !hasDeployment {
		if gatewayid, ok := d.GetOk("gateway_id"); ok {
			deployment, ddiags := resolveGatewayDeployment(ctx, &pco, orgid, envid, gatewayid.(string))
			if ddiags.HasError() {
				return ddiags
			}
			d.Set("deployment", []map[string]any{
				{
					"environment_id":  deployment.GetEnvironmentId(),
					"type":            deployment.GetType(),
					"expected_status": deployment.GetExpectedStatus(),
					"target_id":       deployment.GetTargetId(),
					"target_name":     deployment.GetTargetName(),
					"gateway_version": deployment.GetGatewayVersion(),
				},
			})
		}
	}

	body, err := newAgentsToolsPostBody(d, mcpServerEndpointType)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to build MCP server request",
			Detail:   err.Error(),
		})
		return diags
	}

	res, httpr, err := pco.apimclient.DefaultApi.PostApimInstance(authctx, orgid, envid).ApimInstancePostBody(*body).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create MCP server for org " + orgid + " and env " + envid,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId(strconv.Itoa(int(res.GetId())))

	return resourceMCPServerRead(ctx, d, m)
}

func resourceMCPServerRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	id := d.Id()
	if isComposedResourceId(id) {
		orgid, envid, id, diags = decomposeAgentsToolsId(d)
	}
	if diags.HasError() {
		return diags
	}
	authctx := getApimAuthCtx(ctx, &pco)

	res, httpr, err := pco.apimclient.DefaultApi.GetApimInstanceDetails(authctx, orgid, envid, id).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to read MCP server " + id,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	if err := setAgentsToolsResourceData(d, res); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set MCP server " + id,
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(id)
	d.Set("org_id", orgid)
	d.Set("env_id", envid)

	return diags
}

func resourceMCPServerUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	id := d.Id()

	if d.HasChanges(getAgentsToolsUpdatableAttributes()...) {
		body := newAgentsToolsPatchBody(d, mcpServerEndpointType)
		authctx := getApimAuthCtx(ctx, &pco)
		_, httpr, err := pco.apimclient.DefaultApi.PatchApimInstance(authctx, orgid, envid, id).Body(body).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update MCP server " + id,
				Detail:   parseApiError(httpr, err),
			})
			return diags
		}
		defer httpr.Body.Close()
		return resourceMCPServerRead(ctx, d, m)
	}
	return diags
}

func resourceMCPServerDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	id := d.Id()
	authctx := getApimAuthCtx(ctx, &pco)

	httpr, err := pco.apimclient.DefaultApi.DeleteApimInstance(authctx, orgid, envid, id).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete MCP server " + id,
			Detail:   parseApiError(httpr, err),
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId("")
	return diags
}
