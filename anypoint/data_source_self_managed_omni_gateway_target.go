package anypoint

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	flexgateway "github.com/mulesoft-anypoint/anypoint-client-go/flexgateway"
)

func dataSourceSelfManagedOmniGatewayTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSelfManagedOmniGatewayTargetRead,
		Description: `
		Read all Self-Managed Omni Gateway targets.
		`,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Self-Managed Omni Gateway target's unique id",
			},
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
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Self-Managed Omni Gateway target",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Self-Managed Omni Gateway target",
			},
			"replicas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of replicas by status type",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the Self-Managed Omni Gateway replicas",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of the Self-Managed Omni Gateway replicas",
						},
						"certificate_expiration_dates": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Certificate expiration dates for the given replicas",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of tags",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update date-time",
			},
			"versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of version numbers",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the version number",
			},
		},
	}
}

func dataSourceSelfManagedOmniGatewayTargetRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	id := d.Get("id").(string)
	authctx := getSelfManagedOmniGatewayAuthCtx(ctx, &pco)
	//perform request
	res, httpr, err := pco.flexgatewayclient.DefaultApi.GetFlexGatewayTargetById(authctx, orgid, envid, id).Execute()
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
			Summary:  "Unable to get Self-Managed Omni Gateway target " + id,
			Detail:   details,
		})
		return diags
	}
	defer httpr.Body.Close()
	//parse data
	data := flattenSelfManagedOmniGatewayTargetDetails(res)
	if err := setSelfManagedOmniGatewayTargetAttributesToResourceData(d, data); err != nil {
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set Self-Managed Omni Gateway target attributes",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(string(res.GetId()))
	return diags
}

func flattenSelfManagedOmniGatewayTargetDetails(target *flexgateway.FlexGatewayTargetDetails) map[string]any {
	elem := make(map[string]any)
	if val, ok := target.GetNameOk(); ok && val != nil {
		elem["name"] = *val
	}
	if val, ok := target.GetStatusOk(); ok && val != nil {
		elem["status"] = *val
	}
	if val, ok := target.GetReplicasOk(); ok && val != nil {
		elem["replicas"] = flattenSelfManagedOmniGatewayTargetReplicas(val)
	}
	if val, ok := target.GetTagsOk(); ok && val != nil {
		elem["tags"] = val
	}
	if val, ok := target.GetLastUpdateOk(); ok && val != nil {
		elem["last_update"] = val.String()
	}
	if val, ok := target.GetVersionsOk(); ok && val != nil {
		elem["versions"] = val
	}
	if val, ok := target.GetVersionOk(); ok && val != nil {
		elem["version"] = *val
	}

	return elem
}

func flattenSelfManagedOmniGatewayTargetReplicas(replicas []flexgateway.FlexGatewayTargetDetailsReplicasInner) []map[string]any {
	slice := make([]map[string]any, len(replicas))
	for i, r := range replicas {
		elem := make(map[string]any)
		if val, ok := r.GetStatusOk(); ok && val != nil {
			elem["status"] = *val
		}
		if val, ok := r.GetCountOk(); ok && val != nil {
			elem["count"] = *val
		}
		if dates, ok := r.GetCertificateExpirationDatesOk(); ok && dates != nil {
			strdates := make([]string, len(dates))
			for j, d := range dates {
				strdates[j] = d.String()
			}
			elem["certificate_expiration_dates"] = strdates
		}
		slice[i] = elem
	}
	return slice
}

func setSelfManagedOmniGatewayTargetAttributesToResourceData(d *schema.ResourceData, data map[string]any) error {
	attributes := getSelfManagedOmniGatewayTargetAttributes()
	if data != nil {
		for _, attr := range attributes {
			if val, ok := data[attr]; ok {
				if err := d.Set(attr, val); err != nil {
					return fmt.Errorf("unable to set Self-Managed Omni Gateway target attribute %s\n\tdetails: %s", attr, err)
				}
			}
		}
	}
	return nil
}

func getSelfManagedOmniGatewayTargetAttributes() []string {
	attributes := [...]string{
		"name", "status", "replicas", "tags",
		"last_update", "versions", "version",
	}
	return attributes[:]
}
