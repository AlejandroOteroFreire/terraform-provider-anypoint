package anypoint

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/iancoleman/strcase"
	apim "github.com/mulesoft-anypoint/terraform-provider-anypoint/internal/clients/apim"
)

/*
Shared schema/expand/flatten logic for anypoint_agent_instance and anypoint_mcp_server.

Both resources are, under the hood, API Manager instances (the same `.../apis` endpoint
used by anypoint_apim_mule4) distinguished only by `endpoint.type` ("a2a" for agents,
"mcp" for MCP servers) and by the `technology` value, which the Anypoint UI/API expects
as "flexGateway" on the wire even though it is surfaced as "omniGateway" everywhere else
(see technologyToAPI/technologyFromAPI below - ported from the official
mulesoft/terraform-provider-anypoint agentstools package).
*/

const agentsToolsTechnologyAPIValue = "flexGateway"
const agentsToolsTechnologyUserValue = "omniGateway"

func technologyToAPI(t string) string {
	if t == agentsToolsTechnologyUserValue {
		return agentsToolsTechnologyAPIValue
	}
	return t
}

func technologyFromAPI(t string) string {
	if t == agentsToolsTechnologyAPIValue {
		return agentsToolsTechnologyUserValue
	}
	return t
}

func agentsToolsSchema(endpointType, endpointTypeDescription string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"last_updated": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The last time this resource has been updated locally.",
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The numeric identifier of this instance (stored as string for Terraform compatibility).",
		},
		"org_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The organization id where this instance is defined.",
		},
		"env_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The environment id where this instance is defined.",
		},
		"technology": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          agentsToolsTechnologyUserValue,
			Description:      "The gateway technology. Only 'omniGateway' is currently supported.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{agentsToolsTechnologyUserValue}, false)),
		},
		"provider_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The identity provider id to use for this instance.",
		},
		"instance_label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A human-readable label for this instance.",
		},
		"approval_method": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Client approval method. Valid values: 'manual'. Omit for no approval required.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"manual"}, false)),
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The current status of this instance.",
		},
		"asset_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The Exchange asset id (mirrors spec.asset_id after creation).",
		},
		"asset_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The Exchange asset version (mirrors spec.version after creation).",
		},
		"product_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The product (major) version, computed from the asset version.",
		},
		"consumer_endpoint": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Consumer-facing endpoint URI (the public URL clients use to reach this instance).",
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"upstream_uri": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Shorthand for a single-upstream routing configuration with weight 100. Mutually exclusive with 'routing'.",
		},
		"gateway_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The Managed Omni Gateway id. When provided (and 'deployment' is omitted), the deployment block is auto-populated from this gateway's details.",
		},
		"spec": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "The Exchange asset specification backing this instance.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"asset_id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The Exchange asset id.",
					},
					"group_id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The Exchange group (organization) id.",
					},
					"version": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The asset version.",
					},
				},
			},
		},
		"endpoint": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "Endpoint / proxy configuration for this instance.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"deployment_type": {
						Type:             schema.TypeString,
						Optional:         true,
						Default:          "HY",
						Description:      "Deployment type. Valid values: 'HY' (hybrid), 'CH' (CloudHub), 'CH2' (CloudHub 2.0), 'RF' (Runtime Fabric).",
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HY", "CH", "CH2", "RF"}, false)),
					},
					"type": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: endpointTypeDescription,
					},
					"base_path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Base path (e.g. 'my-agent'). The provider constructs the full proxy URI as http://0.0.0.0:8081/<base_path>. Mutually exclusive with 'uri'.",
					},
					"uri": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Direct implementation URI (e.g. 'http://www.google.com'). Mutually exclusive with 'base_path'.",
					},
					"response_timeout": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Response timeout in milliseconds.",
					},
				},
			},
		},
		"deployment": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "Deployment target configuration. Auto-populated when 'gateway_id' is set.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"environment_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "The environment id for deployment (usually matches the top-level env_id).",
					},
					"type": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "HY",
						Description: "Deployment type. Valid values: 'HY', 'CH', 'RF'.",
					},
					"expected_status": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "deployed",
						Description: "Expected deployment status. Valid values: 'deployed', 'undeployed'.",
					},
					"overwrite": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether to overwrite an existing deployment.",
					},
					"target_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "The target gateway id to deploy to.",
					},
					"target_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The target gateway name.",
					},
					"gateway_version": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The Omni Gateway runtime version.",
					},
				},
			},
		},
		"routing": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Routing rules with weighted upstream backends. Mutually exclusive with 'upstream_uri'.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"label": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "A label for this route.",
					},
					"rules": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Match conditions for this route (methods, path, headers).",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"methods": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Pipe-separated HTTP methods (e.g. 'GET', 'POST|PUT').",
								},
								"path": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "URL path pattern to match (e.g. '/api/*').",
								},
								"host": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Host header value to match.",
								},
								"headers": {
									Type:        schema.TypeMap,
									Optional:    true,
									Description: "Header key-value pairs to match.",
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
							},
						},
					},
					"upstreams": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "Weighted upstream backends for this route.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"weight": {
									Type:        schema.TypeInt,
									Optional:    true,
									Default:     100,
									Description: "Traffic weight percentage (0-100). Weights across upstreams should sum to 100.",
								},
								"uri": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "The upstream backend URI.",
								},
								"label": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "A label for this upstream.",
								},
								"tls_context_id": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "TLS context for upstream connections. Format: 'secretGroupId/tlsContextId'.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// newAgentsToolsPostBody builds the create request body shared by agent instances and MCP servers.
func newAgentsToolsPostBody(d *schema.ResourceData, endpointType string) (*apim.ApimInstancePostBody, error) {
	body := apim.NewApimInstancePostBody()
	body.SetTechnology(technologyToAPI(d.Get("technology").(string)))

	if val, ok := d.GetOk("instance_label"); ok {
		body.SetInstanceLabel(val.(string))
	}

	body.SetSpec(*newAgentsToolsSpec(d))

	endpoint, err := newAgentsToolsEndpoint(d, endpointType)
	if err != nil {
		return nil, err
	}
	body.SetEndpoint(*endpoint)

	if routing := newAgentsToolsRouting(d); routing != nil {
		body.SetRouting(routing)
	}

	if deployment, ok := newAgentsToolsDeployment(d); ok {
		body.SetDeployment(*deployment)
	}

	return body, nil
}

func newAgentsToolsSpec(d *schema.ResourceData) *apim.Spec {
	spec := apim.NewSpec()
	list := d.Get("spec").([]any)
	if len(list) > 0 {
		m := list[0].(map[string]any)
		spec.SetAssetId(m["asset_id"].(string))
		spec.SetGroupId(m["group_id"].(string))
		spec.SetVersion(m["version"].(string))
	}
	return spec
}

func newAgentsToolsEndpoint(d *schema.ResourceData, endpointType string) (*apim.EndpointPostBody, error) {
	endpoint := apim.NewEndpointPostBody()
	endpoint.SetType(endpointType)

	list := d.Get("endpoint").([]any)
	deploymentType := "HY"
	var basePath, directUri string
	var responseTimeout int
	if len(list) > 0 && list[0] != nil {
		m := list[0].(map[string]any)
		if v, ok := m["deployment_type"].(string); ok && v != "" {
			deploymentType = v
		}
		basePath, _ = m["base_path"].(string)
		directUri, _ = m["uri"].(string)
		responseTimeout, _ = m["response_timeout"].(int)
	}
	endpoint.SetDeploymentType(deploymentType)

	if directUri != "" {
		endpoint.SetUri(directUri)
	} else {
		trimmed := strings.TrimPrefix(basePath, "/")
		endpoint.SetProxyUri("http://0.0.0.0:8081/" + trimmed)
	}

	if responseTimeout > 0 {
		endpoint.SetResponseTimeout(strconv.Itoa(responseTimeout))
	}

	return endpoint, nil
}

func newAgentsToolsRouting(d *schema.ResourceData) []apim.RoutingPostBodyInner {
	if upstreamUri, ok := d.GetOk("upstream_uri"); ok {
		u := apim.NewRoutingPostBodyInnerUpstreams()
		u.SetWeight(100)
		u.SetUri(upstreamUri.(string))
		return []apim.RoutingPostBodyInner{{Upstreams: []apim.RoutingPostBodyInnerUpstreams{*u}}}
	}

	list, ok := d.GetOk("routing")
	if !ok {
		return nil
	}
	routes := list.([]any)
	if len(routes) == 0 {
		return nil
	}

	result := make([]apim.RoutingPostBodyInner, len(routes))
	for i, r := range routes {
		rm := r.(map[string]any)
		route := apim.NewRoutingPostBodyInner()
		if label, ok := rm["label"].(string); ok && label != "" {
			route.SetLabel(label)
		}
		if rulesList, ok := rm["rules"].([]any); ok && len(rulesList) > 0 && rulesList[0] != nil {
			rulesMap := rulesList[0].(map[string]any)
			rules := apim.NewRoutingRules()
			if v, ok := rulesMap["methods"].(string); ok && v != "" {
				rules.SetMethods(v)
			}
			if v, ok := rulesMap["path"].(string); ok && v != "" {
				rules.SetPath(v)
			}
			if v, ok := rulesMap["host"].(string); ok && v != "" {
				rules.SetHost(v)
			}
			if headers, ok := rulesMap["headers"].(map[string]any); ok && len(headers) > 0 {
				rules.SetHeaders(headers)
			}
			route.SetRules(*rules)
		}
		upstreamsList, _ := rm["upstreams"].([]any)
		upstreams := make([]apim.RoutingPostBodyInnerUpstreams, len(upstreamsList))
		for j, u := range upstreamsList {
			um := u.(map[string]any)
			upstream := apim.NewRoutingPostBodyInnerUpstreams()
			upstream.SetWeight(int32(um["weight"].(int)))
			upstream.SetUri(um["uri"].(string))
			if v, ok := um["label"].(string); ok && v != "" {
				upstream.SetLabel(v)
			}
			if v, ok := um["tls_context_id"].(string); ok && v != "" {
				upstream.SetTlsContextId(v)
			}
			upstreams[j] = *upstream
		}
		route.SetUpstreams(upstreams)
		result[i] = *route
	}
	return result
}

// newAgentsToolsDeployment builds the deployment block, either from the explicit 'deployment'
// config block or (when 'gateway_id' is set and 'deployment' is empty) by resolving the gateway
// via the Managed Omni Gateway client.
func newAgentsToolsDeployment(d *schema.ResourceData) (*apim.DeploymentPostBody, bool) {
	list := d.Get("deployment").([]any)
	if len(list) > 0 && list[0] != nil {
		m := list[0].(map[string]any)
		deployment := apim.NewDeploymentPostBody()
		if v, ok := m["environment_id"].(string); ok && v != "" {
			deployment.SetEnvironmentId(v)
		} else {
			deployment.SetEnvironmentId(d.Get("env_id").(string))
		}
		if v, ok := m["type"].(string); ok && v != "" {
			deployment.SetType(v)
		}
		if v, ok := m["expected_status"].(string); ok && v != "" {
			deployment.SetExpectedStatus(v)
		}
		if v, ok := m["overwrite"].(bool); ok {
			deployment.SetOverwrite(v)
		}
		if v, ok := m["target_id"].(string); ok && v != "" {
			deployment.SetTargetId(v)
		}
		return deployment, true
	}
	return nil, false
}

// resolveGatewayDeployment fetches gateway details and returns a deployment block populated
// with target_id/target_name/gateway_version, mirroring the official provider's gateway_id shortcut.
func resolveGatewayDeployment(ctx context.Context, pco *ProviderConfOutput, orgid, envid, gatewayid string) (*apim.DeploymentPostBody, diag.Diagnostics) {
	var diags diag.Diagnostics
	authctx := getManagedOmniGatewayAuthCtx(ctx, pco)
	gw, httpr, err := pco.managedomnigatewayclient.DefaultApi.GetGateway(authctx, orgid, envid, gatewayid).Execute()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to resolve gateway_id " + gatewayid,
			Detail:   parseApiError(httpr, err),
		})
		return nil, diags
	}
	defer httpr.Body.Close()

	deployment := apim.NewDeploymentPostBody()
	deployment.SetEnvironmentId(envid)
	deployment.SetType("HY")
	deployment.SetExpectedStatus("deployed")
	deployment.SetTargetId(gw.GetId())
	deployment.SetTargetName(gw.GetName())
	deployment.SetGatewayVersion(gw.GetRuntimeVersion())

	return deployment, diags
}

func setAgentsToolsResourceData(d *schema.ResourceData, details *apim.ApimInstanceDetails) error {
	d.Set("instance_label", details.GetInstanceLabel())
	d.Set("asset_id", details.GetAssetId())
	d.Set("asset_version", details.GetAssetVersion())
	d.Set("product_version", details.GetProductVersion())
	d.Set("status", details.GetStatus())
	d.Set("technology", technologyFromAPI(details.GetTechnology()))
	d.Set("consumer_endpoint", details.GetEndpointUri())

	spec := []map[string]any{
		{
			"asset_id": details.GetAssetId(),
			"group_id": details.GetGroupId(),
			"version":  details.GetAssetVersion(),
		},
	}
	d.Set("spec", spec)

	if endpoint := details.Endpoint; endpoint != nil {
		responseTimeout := 0
		if rt, ok := endpoint.GetResponseTimeoutOk(); ok && rt != nil {
			if v, err := strconv.Atoi(*rt); err == nil {
				responseTimeout = v
			}
		}
		d.Set("endpoint", []map[string]any{
			{
				"deployment_type":  endpoint.GetDeploymentType(),
				"type":             endpoint.GetType(),
				"base_path":        "",
				"uri":              endpoint.GetUri(),
				"response_timeout": responseTimeout,
			},
		})
	}

	if dep, ok := details.GetDeploymentOk(); ok && dep != nil {
		d.Set("deployment", []map[string]any{
			{
				"environment_id":  dep.GetEnvironmentId(),
				"type":            dep.GetType(),
				"expected_status": dep.GetExpectedStatus(),
				"target_id":       dep.GetTargetId(),
				"target_name":     dep.GetTargetName(),
				"gateway_version": dep.GetGatewayVersion(),
			},
		})
	}

	if len(details.Routing) > 0 {
		routing := make([]map[string]any, len(details.Routing))
		for i, r := range details.Routing {
			upstreams := make([]map[string]any, len(r.Upstreams))
			for j, u := range r.Upstreams {
				upstreams[j] = map[string]any{
					"weight": int(u.GetWeight()),
					"uri":    u.GetUri(),
					"label":  u.GetLabel(),
				}
			}
			route := map[string]any{
				"label":     r.GetLabel(),
				"upstreams": upstreams,
			}
			if rules := r.Rules; rules != nil {
				route["rules"] = []map[string]any{
					{
						"methods": rules.GetMethods(),
						"path":    rules.GetPath(),
						"host":    rules.GetHost(),
					},
				}
			}
			routing[i] = route
		}
		d.Set("routing", routing)
	}

	return nil
}

func decomposeAgentsToolsId(d *schema.ResourceData) (string, string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	s := DecomposeResourceId(d.Id())
	if len(s) != 3 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid ID format",
			Detail:   "Expected ORG_ID/ENV_ID/INSTANCE_ID, got " + d.Id(),
		})
		return "", "", "", diags
	}
	return s[0], s[1], s[2], diags
}

func getAgentsToolsUpdatableAttributes() []string {
	return []string{
		"instance_label", "provider_id", "approval_method", "consumer_endpoint",
		"endpoint", "spec", "routing", "upstream_uri", "deployment",
	}
}

// newAgentsToolsPatchBody builds a raw PATCH body (map) for fields not covered - or covered
// differently - by the strongly typed create body, mirroring apim_mule4's PATCH pattern.
func newAgentsToolsPatchBody(d *schema.ResourceData, endpointType string) map[string]any {
	body := make(map[string]any)
	if d.HasChange("instance_label") {
		body["instanceLabel"] = d.Get("instance_label").(string)
	}
	if d.HasChange("provider_id") {
		body["providerId"] = d.Get("provider_id").(string)
	}
	if d.HasChange("approval_method") {
		body["approvalMethod"] = d.Get("approval_method").(string)
	}
	if d.HasChange("consumer_endpoint") {
		body["endpointUri"] = d.Get("consumer_endpoint").(string)
	}
	if d.HasChanges("spec") {
		spec := newAgentsToolsSpec(d)
		body[strcase.ToLowerCamel("spec")] = spec
	}
	if d.HasChanges("endpoint") {
		if endpoint, err := newAgentsToolsEndpoint(d, endpointType); err == nil {
			body["endpoint"] = endpoint
		}
	}
	if d.HasChanges("routing", "upstream_uri") {
		if routing := newAgentsToolsRouting(d); routing != nil {
			body["routing"] = routing
		}
	}
	if d.HasChanges("deployment") {
		if deployment, ok := newAgentsToolsDeployment(d); ok {
			body["deployment"] = deployment
		}
	}
	return body
}

// agentsToolsListItem is the flattened shape returned by the list data sources.
type agentsToolsListItem struct {
	id             string
	instanceLabel  string
	status         string
	assetId        string
	assetVersion   string
	productVersion string
	environmentId  string
}

// listAgentsToolsInstances lists every API Manager instance for org/env, keeps only those
// with technology == "flexGateway" (the wire value for "omniGateway"), then fetches each
// candidate's details to keep only the ones whose endpoint.type matches endpointType
// ("a2a" for agent instances, "mcp" for MCP servers). The list endpoint's response shape
// does not expose endpoint.type, so this requires one detail call per flexGateway instance.
func listAgentsToolsInstances(ctx context.Context, pco *ProviderConfOutput, orgid, envid, endpointType string) ([]agentsToolsListItem, error) {
	authctx := getApimAuthCtx(ctx, pco)
	collection, httpr, err := pco.apimclient.DefaultApi.GetEnvApimInstances(authctx, orgid, envid).Execute()
	if err != nil {
		return nil, fmt.Errorf("%s", parseApiError(httpr, err))
	}
	defer httpr.Body.Close()

	var result []agentsToolsListItem
	for _, asset := range collection.GetAssets() {
		for _, api := range asset.GetApis() {
			if api.GetTechnology() != agentsToolsTechnologyAPIValue {
				continue
			}
			id := strconv.Itoa(int(api.GetId()))

			details, dhttpr, derr := pco.apimclient.DefaultApi.GetApimInstanceDetails(authctx, orgid, envid, id).Execute()
			if derr != nil {
				return nil, fmt.Errorf("unable to read details for instance %s: %s", id, parseApiError(dhttpr, derr))
			}
			dhttpr.Body.Close()

			if details.Endpoint == nil || details.Endpoint.GetType() != endpointType {
				continue
			}

			result = append(result, agentsToolsListItem{
				id:             id,
				instanceLabel:  api.GetInstanceLabel(),
				status:         api.GetStatus(),
				assetId:        api.GetAssetId(),
				assetVersion:   api.GetAssetVersion(),
				productVersion: api.GetProductVersion(),
				environmentId:  api.GetEnvironmentId(),
			})
		}
	}
	return result, nil
}

func flattenAgentsToolsListItems(items []agentsToolsListItem) []map[string]any {
	result := make([]map[string]any, len(items))
	for i, item := range items {
		result[i] = map[string]any{
			"id":              item.id,
			"instance_label":  item.instanceLabel,
			"status":          item.status,
			"asset_id":        item.assetId,
			"asset_version":   item.assetVersion,
			"product_version": item.productVersion,
			"environment_id":  item.environmentId,
		}
	}
	return result
}
