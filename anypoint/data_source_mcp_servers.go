package anypoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMCPServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMCPServersRead,
		Description: `
		Lists all ` + "`" + `MCP server instances` + "`" + ` (technology 'omniGateway', endpoint type 'mcp') in an environment.
		`,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization id where the MCP servers are defined.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment id where the MCP servers are defined.",
			},
			"servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of MCP server instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The numeric identifier of this instance (stored as string for Terraform compatibility).",
						},
						"instance_label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance's label.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The current status of this instance.",
						},
						"asset_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Exchange asset id.",
						},
						"asset_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Exchange asset version.",
						},
						"product_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product (major) version.",
						},
						"environment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment id this instance is deployed to.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMCPServersRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)

	items, err := listAgentsToolsInstances(ctx, &pco, orgid, envid, mcpServerEndpointType)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to list MCP servers for org " + orgid + " and env " + envid,
			Detail:   err.Error(),
		})
		return diags
	}

	d.Set("servers", flattenAgentsToolsListItems(items))
	d.SetId(orgid + "/" + envid)

	return diags
}
