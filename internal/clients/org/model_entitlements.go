package org

// Entitlements models the `entitlements` object shared by the GET response and the
// POST/PUT request bodies of the Access Management "organizations" endpoints
// (mulesoft/mulesoft-dx apis/access-management/api.yaml). The real API accepts every
// field below on create/update - this replaces the external anypoint-client-go/org
// module, whose EntitlementsCore type only modeled 10 of them.
//
// Every field is a pointer so a Terraform "not configured" value can be told apart
// from an explicit zero, matching how resource_bg.go/data_source_bg.go already use it.
type Entitlements struct {
	Hybrid                       *HybridEntitlement             `json:"hybrid,omitempty"`
	HybridInsight                *bool                          `json:"hybridInsight,omitempty"`
	HybridAutoDiscoverProperties *bool                          `json:"hybridAutoDiscoverProperties,omitempty"`
	CreateEnvironments           *bool                          `json:"createEnvironments,omitempty"`
	GlobalDeployment             *bool                          `json:"globalDeployment,omitempty"`
	CreateSubOrgs                *bool                          `json:"createSubOrgs,omitempty"`
	VCoresProduction             *VCoresEntitlement             `json:"vCoresProduction,omitempty"`
	VCoresSandbox                *VCoresEntitlement             `json:"vCoresSandbox,omitempty"`
	VCoresDesign                 *VCoresEntitlement             `json:"vCoresDesign,omitempty"`
	StaticIps                    *AssignedReassignedEntitlement `json:"staticIps,omitempty"`
	Vpcs                         *AssignedReassignedEntitlement `json:"vpcs,omitempty"`
	NetworkConnections           *AssignedReassignedEntitlement `json:"networkConnections,omitempty"`
	WorkerLoggingOverride        *EnabledEntitlement            `json:"workerLoggingOverride,omitempty"`
	Messaging                    *AssignedEntitlement           `json:"messaging,omitempty"`
	MqMessages                   *BaseAddOnEntitlement          `json:"mqMessages,omitempty"`
	MqRequests                   *BaseAddOnEntitlement          `json:"mqRequests,omitempty"`
	MqAdvancedFeatures           *EnabledEntitlement            `json:"mqAdvancedFeatures,omitempty"`
	ObjectStoreRequestUnits      *BaseAddOnEntitlement          `json:"objectStoreRequestUnits,omitempty"`
	ObjectStoreKeys              *BaseAddOnEntitlement          `json:"objectStoreKeys,omitempty"`
	Gateways                     *AssignedEntitlement           `json:"gateways,omitempty"`
	DesignCenter                 *DesignCenterEntitlement       `json:"designCenter,omitempty"`
	PartnersProduction           *AssignedEntitlement           `json:"partnersProduction,omitempty"`
	PartnersSandbox              *AssignedEntitlement           `json:"partnersSandbox,omitempty"`
	TradingPartnersProduction    *AssignedEntitlement           `json:"tradingPartnersProduction,omitempty"`
	TradingPartnersSandbox       *AssignedEntitlement           `json:"tradingPartnersSandbox,omitempty"`
	LoadBalancer                 *AssignedReassignedEntitlement `json:"loadBalancer,omitempty"`
	ExternalIdentity             *bool                          `json:"externalIdentity,omitempty"`
	Autoscaling                  *bool                          `json:"autoscaling,omitempty"`
	ArmAlerts                    *bool                          `json:"armAlerts,omitempty"`
	Apis                         *EnabledEntitlement            `json:"apis,omitempty"`
	ApiMonitoring                *ApiMonitoringEntitlement      `json:"apiMonitoring,omitempty"`
	ApiCommunityManager          *EnabledEntitlement            `json:"apiCommunityManager,omitempty"`
	MonitoringCenter             *MonitoringCenterEntitlement   `json:"monitoringCenter,omitempty"`
	ApiQuery                     *ApiQueryEntitlement           `json:"apiQuery,omitempty"`
	ApiQueryC360                 *EnabledEntitlement            `json:"apiQueryC360,omitempty"`
	AngGovernance                *AngGovernanceEntitlement      `json:"apiGovernance,omitempty"`
	Crowd                        *CrowdEntitlement              `json:"crowd,omitempty"`
	Cam                          *EnabledEntitlement            `json:"cam,omitempty"`
	Exchange2                    *EnabledEntitlement            `json:"exchange2,omitempty"`
	CrowdSelfServiceMigration    *EnabledEntitlement            `json:"crowdSelfServiceMigration,omitempty"`
	KpiDashboard                 *EnabledEntitlement            `json:"kpiDashboard,omitempty"`
	Pcf                          *bool                          `json:"pcf,omitempty"`
	AppViz                       *bool                          `json:"appViz,omitempty"`
	RuntimeFabric                *bool                          `json:"runtimeFabric,omitempty"`
	AnypointSecurityTokenization *EnabledEntitlement            `json:"anypointSecurityTokenization,omitempty"`
	AnypointSecurityEdgePolicies *EnabledEntitlement            `json:"anypointSecurityEdgePolicies,omitempty"`
	RuntimeFabricCloud           *EnabledEntitlement            `json:"runtimeFabricCloud,omitempty"`
	ServiceMesh                  *EnabledEntitlement            `json:"serviceMesh,omitempty"`
	WorkerClouds                 *AssignedReassignedEntitlement `json:"workerClouds,omitempty"`
	ManagedGatewaySmall          *AssignedReassignedEntitlement `json:"managedGatewaySmall,omitempty"`
	ManagedGatewayLarge          *AssignedReassignedEntitlement `json:"managedGatewayLarge,omitempty"`
}

func NewEntitlementsCore() *Entitlements { return &Entitlements{} }

// --- simple boolean/top-level accessors ---

func (o *Entitlements) GetHybrid() HybridEntitlement {
	if o == nil || o.Hybrid == nil {
		return HybridEntitlement{}
	}
	return *o.Hybrid
}
func (o *Entitlements) SetHybrid(v HybridEntitlement) { o.Hybrid = &v }

func (o *Entitlements) GetHybridInsight() bool {
	if o == nil || o.HybridInsight == nil {
		return false
	}
	return *o.HybridInsight
}
func (o *Entitlements) SetHybridInsight(v bool) { o.HybridInsight = &v }

func (o *Entitlements) GetHybridAutoDiscoverProperties() bool {
	if o == nil || o.HybridAutoDiscoverProperties == nil {
		return false
	}
	return *o.HybridAutoDiscoverProperties
}
func (o *Entitlements) SetHybridAutoDiscoverProperties(v bool) { o.HybridAutoDiscoverProperties = &v }

func (o *Entitlements) GetCreateEnvironments() bool {
	if o == nil || o.CreateEnvironments == nil {
		return false
	}
	return *o.CreateEnvironments
}
func (o *Entitlements) SetCreateEnvironments(v bool) { o.CreateEnvironments = &v }

func (o *Entitlements) GetGlobalDeployment() bool {
	if o == nil || o.GlobalDeployment == nil {
		return false
	}
	return *o.GlobalDeployment
}
func (o *Entitlements) SetGlobalDeployment(v bool) { o.GlobalDeployment = &v }

func (o *Entitlements) GetCreateSubOrgs() bool {
	if o == nil || o.CreateSubOrgs == nil {
		return false
	}
	return *o.CreateSubOrgs
}
func (o *Entitlements) SetCreateSubOrgs(v bool) { o.CreateSubOrgs = &v }

func (o *Entitlements) GetExternalIdentity() bool {
	if o == nil || o.ExternalIdentity == nil {
		return false
	}
	return *o.ExternalIdentity
}
func (o *Entitlements) SetExternalIdentity(v bool) { o.ExternalIdentity = &v }

func (o *Entitlements) GetAutoscaling() bool {
	if o == nil || o.Autoscaling == nil {
		return false
	}
	return *o.Autoscaling
}
func (o *Entitlements) SetAutoscaling(v bool) { o.Autoscaling = &v }

func (o *Entitlements) GetArmAlerts() bool {
	if o == nil || o.ArmAlerts == nil {
		return false
	}
	return *o.ArmAlerts
}
func (o *Entitlements) SetArmAlerts(v bool) { o.ArmAlerts = &v }

func (o *Entitlements) GetPcf() bool {
	if o == nil || o.Pcf == nil {
		return false
	}
	return *o.Pcf
}
func (o *Entitlements) SetPcf(v bool) { o.Pcf = &v }

func (o *Entitlements) GetAppViz() bool {
	if o == nil || o.AppViz == nil {
		return false
	}
	return *o.AppViz
}
func (o *Entitlements) SetAppViz(v bool) { o.AppViz = &v }

func (o *Entitlements) GetRuntimeFabric() bool {
	if o == nil || o.RuntimeFabric == nil {
		return false
	}
	return *o.RuntimeFabric
}
func (o *Entitlements) SetRuntimeFabric(v bool) { o.RuntimeFabric = &v }

// --- nested-object accessors ---

func (o *Entitlements) GetVCoresProduction() VCoresEntitlement {
	if o == nil || o.VCoresProduction == nil {
		return VCoresEntitlement{}
	}
	return *o.VCoresProduction
}
func (o *Entitlements) SetVCoresProduction(v VCoresEntitlement) { o.VCoresProduction = &v }

func (o *Entitlements) GetVCoresSandbox() VCoresEntitlement {
	if o == nil || o.VCoresSandbox == nil {
		return VCoresEntitlement{}
	}
	return *o.VCoresSandbox
}
func (o *Entitlements) SetVCoresSandbox(v VCoresEntitlement) { o.VCoresSandbox = &v }

func (o *Entitlements) GetVCoresDesign() VCoresEntitlement {
	if o == nil || o.VCoresDesign == nil {
		return VCoresEntitlement{}
	}
	return *o.VCoresDesign
}
func (o *Entitlements) SetVCoresDesign(v VCoresEntitlement) { o.VCoresDesign = &v }

func (o *Entitlements) GetStaticIps() AssignedReassignedEntitlement {
	if o == nil || o.StaticIps == nil {
		return AssignedReassignedEntitlement{}
	}
	return *o.StaticIps
}
func (o *Entitlements) SetStaticIps(v AssignedReassignedEntitlement) { o.StaticIps = &v }

func (o *Entitlements) GetVpcs() AssignedReassignedEntitlement {
	if o == nil || o.Vpcs == nil {
		return AssignedReassignedEntitlement{}
	}
	return *o.Vpcs
}
func (o *Entitlements) SetVpcs(v AssignedReassignedEntitlement) { o.Vpcs = &v }

func (o *Entitlements) GetVpns() AssignedReassignedEntitlement  { return o.GetNetworkConnections() }
func (o *Entitlements) SetVpns(v AssignedReassignedEntitlement) { o.SetNetworkConnections(v) }

func (o *Entitlements) GetNetworkConnections() AssignedReassignedEntitlement {
	if o == nil || o.NetworkConnections == nil {
		return AssignedReassignedEntitlement{}
	}
	return *o.NetworkConnections
}
func (o *Entitlements) SetNetworkConnections(v AssignedReassignedEntitlement) {
	o.NetworkConnections = &v
}

func (o *Entitlements) GetWorkerLoggingOverride() EnabledEntitlement {
	if o == nil || o.WorkerLoggingOverride == nil {
		return EnabledEntitlement{}
	}
	return *o.WorkerLoggingOverride
}
func (o *Entitlements) SetWorkerLoggingOverride(v EnabledEntitlement) { o.WorkerLoggingOverride = &v }

func (o *Entitlements) GetMessaging() AssignedEntitlement {
	if o == nil || o.Messaging == nil {
		return AssignedEntitlement{}
	}
	return *o.Messaging
}
func (o *Entitlements) SetMessaging(v AssignedEntitlement) { o.Messaging = &v }

func (o *Entitlements) GetMqMessages() BaseAddOnEntitlement {
	if o == nil || o.MqMessages == nil {
		return BaseAddOnEntitlement{}
	}
	return *o.MqMessages
}
func (o *Entitlements) SetMqMessages(v BaseAddOnEntitlement) { o.MqMessages = &v }

func (o *Entitlements) GetMqRequests() BaseAddOnEntitlement {
	if o == nil || o.MqRequests == nil {
		return BaseAddOnEntitlement{}
	}
	return *o.MqRequests
}
func (o *Entitlements) SetMqRequests(v BaseAddOnEntitlement) { o.MqRequests = &v }

func (o *Entitlements) GetMqAdvancedFeatures() EnabledEntitlement {
	if o == nil || o.MqAdvancedFeatures == nil {
		return EnabledEntitlement{}
	}
	return *o.MqAdvancedFeatures
}
func (o *Entitlements) SetMqAdvancedFeatures(v EnabledEntitlement) { o.MqAdvancedFeatures = &v }

func (o *Entitlements) GetObjectStoreRequestUnits() BaseAddOnEntitlement {
	if o == nil || o.ObjectStoreRequestUnits == nil {
		return BaseAddOnEntitlement{}
	}
	return *o.ObjectStoreRequestUnits
}
func (o *Entitlements) SetObjectStoreRequestUnits(v BaseAddOnEntitlement) {
	o.ObjectStoreRequestUnits = &v
}

func (o *Entitlements) GetObjectStoreKeys() BaseAddOnEntitlement {
	if o == nil || o.ObjectStoreKeys == nil {
		return BaseAddOnEntitlement{}
	}
	return *o.ObjectStoreKeys
}
func (o *Entitlements) SetObjectStoreKeys(v BaseAddOnEntitlement) { o.ObjectStoreKeys = &v }

func (o *Entitlements) GetGateways() AssignedEntitlement {
	if o == nil || o.Gateways == nil {
		return AssignedEntitlement{}
	}
	return *o.Gateways
}
func (o *Entitlements) SetGateways(v AssignedEntitlement) { o.Gateways = &v }

func (o *Entitlements) GetDesignCenter() DesignCenterEntitlement {
	if o == nil || o.DesignCenter == nil {
		return DesignCenterEntitlement{}
	}
	return *o.DesignCenter
}
func (o *Entitlements) SetDesignCenter(v DesignCenterEntitlement) { o.DesignCenter = &v }

func (o *Entitlements) GetPartnersProduction() AssignedEntitlement {
	if o == nil || o.PartnersProduction == nil {
		return AssignedEntitlement{}
	}
	return *o.PartnersProduction
}
func (o *Entitlements) SetPartnersProduction(v AssignedEntitlement) { o.PartnersProduction = &v }

func (o *Entitlements) GetPartnersSandbox() AssignedEntitlement {
	if o == nil || o.PartnersSandbox == nil {
		return AssignedEntitlement{}
	}
	return *o.PartnersSandbox
}
func (o *Entitlements) SetPartnersSandbox(v AssignedEntitlement) { o.PartnersSandbox = &v }

func (o *Entitlements) GetTradingPartnersProduction() AssignedEntitlement {
	if o == nil || o.TradingPartnersProduction == nil {
		return AssignedEntitlement{}
	}
	return *o.TradingPartnersProduction
}
func (o *Entitlements) SetTradingPartnersProduction(v AssignedEntitlement) {
	o.TradingPartnersProduction = &v
}

func (o *Entitlements) GetTradingPartnersSandbox() AssignedEntitlement {
	if o == nil || o.TradingPartnersSandbox == nil {
		return AssignedEntitlement{}
	}
	return *o.TradingPartnersSandbox
}
func (o *Entitlements) SetTradingPartnersSandbox(v AssignedEntitlement) {
	o.TradingPartnersSandbox = &v
}

func (o *Entitlements) GetLoadBalancer() AssignedReassignedEntitlement {
	if o == nil || o.LoadBalancer == nil {
		return AssignedReassignedEntitlement{}
	}
	return *o.LoadBalancer
}
func (o *Entitlements) SetLoadBalancer(v AssignedReassignedEntitlement) { o.LoadBalancer = &v }

func (o *Entitlements) GetApis() EnabledEntitlement {
	if o == nil || o.Apis == nil {
		return EnabledEntitlement{}
	}
	return *o.Apis
}
func (o *Entitlements) SetApis(v EnabledEntitlement) { o.Apis = &v }

func (o *Entitlements) GetApiMonitoring() ApiMonitoringEntitlement {
	if o == nil || o.ApiMonitoring == nil {
		return ApiMonitoringEntitlement{}
	}
	return *o.ApiMonitoring
}
func (o *Entitlements) SetApiMonitoring(v ApiMonitoringEntitlement) { o.ApiMonitoring = &v }

func (o *Entitlements) GetApiCommunityManager() EnabledEntitlement {
	if o == nil || o.ApiCommunityManager == nil {
		return EnabledEntitlement{}
	}
	return *o.ApiCommunityManager
}
func (o *Entitlements) SetApiCommunityManager(v EnabledEntitlement) { o.ApiCommunityManager = &v }

func (o *Entitlements) GetMonitoringCenter() MonitoringCenterEntitlement {
	if o == nil || o.MonitoringCenter == nil {
		return MonitoringCenterEntitlement{}
	}
	return *o.MonitoringCenter
}
func (o *Entitlements) SetMonitoringCenter(v MonitoringCenterEntitlement) { o.MonitoringCenter = &v }

func (o *Entitlements) GetApiQuery() ApiQueryEntitlement {
	if o == nil || o.ApiQuery == nil {
		return ApiQueryEntitlement{}
	}
	return *o.ApiQuery
}
func (o *Entitlements) SetApiQuery(v ApiQueryEntitlement) { o.ApiQuery = &v }

func (o *Entitlements) GetApiQueryC360() EnabledEntitlement {
	if o == nil || o.ApiQueryC360 == nil {
		return EnabledEntitlement{}
	}
	return *o.ApiQueryC360
}
func (o *Entitlements) SetApiQueryC360(v EnabledEntitlement) { o.ApiQueryC360 = &v }

func (o *Entitlements) GetAngGovernance() AngGovernanceEntitlement {
	if o == nil || o.AngGovernance == nil {
		return AngGovernanceEntitlement{}
	}
	return *o.AngGovernance
}
func (o *Entitlements) SetAngGovernance(v AngGovernanceEntitlement) { o.AngGovernance = &v }

func (o *Entitlements) GetCrowd() CrowdEntitlement {
	if o == nil || o.Crowd == nil {
		return CrowdEntitlement{}
	}
	return *o.Crowd
}
func (o *Entitlements) SetCrowd(v CrowdEntitlement) { o.Crowd = &v }

func (o *Entitlements) GetCam() EnabledEntitlement {
	if o == nil || o.Cam == nil {
		return EnabledEntitlement{}
	}
	return *o.Cam
}
func (o *Entitlements) SetCam(v EnabledEntitlement) { o.Cam = &v }

func (o *Entitlements) GetExchange2() EnabledEntitlement {
	if o == nil || o.Exchange2 == nil {
		return EnabledEntitlement{}
	}
	return *o.Exchange2
}
func (o *Entitlements) SetExchange2(v EnabledEntitlement) { o.Exchange2 = &v }

func (o *Entitlements) GetCrowdSelfServiceMigration() EnabledEntitlement {
	if o == nil || o.CrowdSelfServiceMigration == nil {
		return EnabledEntitlement{}
	}
	return *o.CrowdSelfServiceMigration
}
func (o *Entitlements) SetCrowdSelfServiceMigration(v EnabledEntitlement) {
	o.CrowdSelfServiceMigration = &v
}

func (o *Entitlements) GetKpiDashboard() EnabledEntitlement {
	if o == nil || o.KpiDashboard == nil {
		return EnabledEntitlement{}
	}
	return *o.KpiDashboard
}
func (o *Entitlements) SetKpiDashboard(v EnabledEntitlement) { o.KpiDashboard = &v }

func (o *Entitlements) GetAnypointSecurityTokenization() EnabledEntitlement {
	if o == nil || o.AnypointSecurityTokenization == nil {
		return EnabledEntitlement{}
	}
	return *o.AnypointSecurityTokenization
}
func (o *Entitlements) SetAnypointSecurityTokenization(v EnabledEntitlement) {
	o.AnypointSecurityTokenization = &v
}

func (o *Entitlements) GetAnypointSecurityEdgePolicies() EnabledEntitlement {
	if o == nil || o.AnypointSecurityEdgePolicies == nil {
		return EnabledEntitlement{}
	}
	return *o.AnypointSecurityEdgePolicies
}
func (o *Entitlements) SetAnypointSecurityEdgePolicies(v EnabledEntitlement) {
	o.AnypointSecurityEdgePolicies = &v
}

func (o *Entitlements) GetRuntimeFabricCloud() EnabledEntitlement {
	if o == nil || o.RuntimeFabricCloud == nil {
		return EnabledEntitlement{}
	}
	return *o.RuntimeFabricCloud
}
func (o *Entitlements) SetRuntimeFabricCloud(v EnabledEntitlement) { o.RuntimeFabricCloud = &v }

func (o *Entitlements) GetServiceMesh() EnabledEntitlement {
	if o == nil || o.ServiceMesh == nil {
		return EnabledEntitlement{}
	}
	return *o.ServiceMesh
}
func (o *Entitlements) SetServiceMesh(v EnabledEntitlement) { o.ServiceMesh = &v }

func (o *Entitlements) GetWorkerClouds() AssignedReassignedEntitlement {
	if o == nil || o.WorkerClouds == nil {
		return AssignedReassignedEntitlement{}
	}
	return *o.WorkerClouds
}
func (o *Entitlements) SetWorkerClouds(v AssignedReassignedEntitlement) { o.WorkerClouds = &v }

func (o *Entitlements) GetManagedGatewaySmall() AssignedReassignedEntitlement {
	if o == nil || o.ManagedGatewaySmall == nil {
		return AssignedReassignedEntitlement{}
	}
	return *o.ManagedGatewaySmall
}
func (o *Entitlements) SetManagedGatewaySmall(v AssignedReassignedEntitlement) {
	o.ManagedGatewaySmall = &v
}

func (o *Entitlements) GetManagedGatewayLarge() AssignedReassignedEntitlement {
	if o == nil || o.ManagedGatewayLarge == nil {
		return AssignedReassignedEntitlement{}
	}
	return *o.ManagedGatewayLarge
}
func (o *Entitlements) SetManagedGatewayLarge(v AssignedReassignedEntitlement) {
	o.ManagedGatewayLarge = &v
}

// --- shared small entitlement shapes ---

type HybridEntitlement struct {
	Enabled *bool `json:"enabled,omitempty"`
}

func (o HybridEntitlement) GetEnabled() bool {
	if o.Enabled == nil {
		return false
	}
	return *o.Enabled
}
func (o *HybridEntitlement) SetEnabled(v bool) { o.Enabled = &v }

type EnabledEntitlement struct {
	Enabled *bool `json:"enabled,omitempty"`
}

func (o EnabledEntitlement) GetEnabled() bool {
	if o.Enabled == nil {
		return false
	}
	return *o.Enabled
}
func (o *EnabledEntitlement) SetEnabled(v bool) { o.Enabled = &v }

type AssignedEntitlement struct {
	Assigned *int32 `json:"assigned,omitempty"`
}

func (o AssignedEntitlement) GetAssigned() int32 {
	if o.Assigned == nil {
		return 0
	}
	return *o.Assigned
}
func (o *AssignedEntitlement) SetAssigned(v int32) { o.Assigned = &v }

type AssignedReassignedEntitlement struct {
	Assigned   *int32 `json:"assigned,omitempty"`
	Reassigned *int32 `json:"reassigned,omitempty"`
}

func (o AssignedReassignedEntitlement) GetAssigned() int32 {
	if o.Assigned == nil {
		return 0
	}
	return *o.Assigned
}
func (o *AssignedReassignedEntitlement) SetAssigned(v int32) { o.Assigned = &v }

func (o AssignedReassignedEntitlement) GetReassigned() int32 {
	if o.Reassigned == nil {
		return 0
	}
	return *o.Reassigned
}
func (o *AssignedReassignedEntitlement) SetReassigned(v int32) { o.Reassigned = &v }

// VCoresEntitlement mirrors AssignedReassignedEntitlement but with float32 units
// (vCoresProduction/Sandbox/Design are fractional, e.g. 0.5).
type VCoresEntitlement struct {
	Assigned   *float32 `json:"assigned,omitempty"`
	Reassigned *float32 `json:"reassigned,omitempty"`
}

func (o VCoresEntitlement) GetAssigned() float32 {
	if o.Assigned == nil {
		return 0
	}
	return *o.Assigned
}
func (o *VCoresEntitlement) SetAssigned(v float32) { o.Assigned = &v }

func (o VCoresEntitlement) GetReassigned() float32 {
	if o.Reassigned == nil {
		return 0
	}
	return *o.Reassigned
}
func (o *VCoresEntitlement) SetReassigned(v float32) { o.Reassigned = &v }

type BaseAddOnEntitlement struct {
	Base  *int32 `json:"base,omitempty"`
	AddOn *int32 `json:"addOn,omitempty"`
}

func (o BaseAddOnEntitlement) GetBase() int32 {
	if o.Base == nil {
		return 0
	}
	return *o.Base
}
func (o *BaseAddOnEntitlement) SetBase(v int32) { o.Base = &v }

func (o BaseAddOnEntitlement) GetAddOn() int32 {
	if o.AddOn == nil {
		return 0
	}
	return *o.AddOn
}
func (o *BaseAddOnEntitlement) SetAddOn(v int32) { o.AddOn = &v }

type DesignCenterEntitlement struct {
	Api    *bool `json:"api,omitempty"`
	Mozart *bool `json:"mozart,omitempty"`
}

func (o DesignCenterEntitlement) GetApi() bool {
	if o.Api == nil {
		return false
	}
	return *o.Api
}
func (o *DesignCenterEntitlement) SetApi(v bool) { o.Api = &v }

func (o DesignCenterEntitlement) GetMozart() bool {
	if o.Mozart == nil {
		return false
	}
	return *o.Mozart
}
func (o *DesignCenterEntitlement) SetMozart(v bool) { o.Mozart = &v }

type ApiMonitoringEntitlement struct {
	Schedules *int32 `json:"schedules,omitempty"`
}

func (o ApiMonitoringEntitlement) GetSchedules() int32 {
	if o.Schedules == nil {
		return 0
	}
	return *o.Schedules
}
func (o *ApiMonitoringEntitlement) SetSchedules(v int32) { o.Schedules = &v }

type MonitoringCenterEntitlement struct {
	ProductSKU *int32 `json:"productSKU,omitempty"`
}

func (o MonitoringCenterEntitlement) GetProductSKU() int32 {
	if o.ProductSKU == nil {
		return 0
	}
	return *o.ProductSKU
}
func (o *MonitoringCenterEntitlement) SetProductSKU(v int32) { o.ProductSKU = &v }

type ApiQueryEntitlement struct {
	Enabled    *bool  `json:"enabled,omitempty"`
	ProductSKU *int32 `json:"productSKU,omitempty"`
}

func (o ApiQueryEntitlement) GetEnabled() bool {
	if o.Enabled == nil {
		return false
	}
	return *o.Enabled
}
func (o *ApiQueryEntitlement) SetEnabled(v bool) { o.Enabled = &v }

func (o ApiQueryEntitlement) GetProductSKU() int32 {
	if o.ProductSKU == nil {
		return 0
	}
	return *o.ProductSKU
}
func (o *ApiQueryEntitlement) SetProductSKU(v int32) { o.ProductSKU = &v }

type AngGovernanceEntitlement struct {
	Level *int32 `json:"apisPerMonth,omitempty"`
}

func (o AngGovernanceEntitlement) GetLevel() int32 {
	if o.Level == nil {
		return 0
	}
	return *o.Level
}
func (o *AngGovernanceEntitlement) SetLevel(v int32) { o.Level = &v }

type CrowdEntitlement struct {
	HideApiManagerDesigner *bool `json:"hideApiManagerDesigner,omitempty"`
	HideFormerApiPlatform  *bool `json:"hideFormerApiPlatform,omitempty"`
	Environments           *bool `json:"environments,omitempty"`
}

func (o CrowdEntitlement) GetHideApiManagerDesigner() bool {
	if o.HideApiManagerDesigner == nil {
		return false
	}
	return *o.HideApiManagerDesigner
}
func (o *CrowdEntitlement) SetHideApiManagerDesigner(v bool) { o.HideApiManagerDesigner = &v }

func (o CrowdEntitlement) GetHideFormerApiPlatform() bool {
	if o.HideFormerApiPlatform == nil {
		return false
	}
	return *o.HideFormerApiPlatform
}
func (o *CrowdEntitlement) SetHideFormerApiPlatform(v bool) { o.HideFormerApiPlatform = &v }

func (o CrowdEntitlement) GetEnvironments() bool {
	if o.Environments == nil {
		return false
	}
	return *o.Environments
}
func (o *CrowdEntitlement) SetEnvironments(v bool) { o.Environments = &v }

// --- type aliases matching the external SDK's naming, to minimize call-site churn ---

type EntitlementsCore = Entitlements
type LoadBalancerEntitlement = AssignedReassignedEntitlement
type StaticIpsEntitlement = AssignedReassignedEntitlement
type VCoresSandboxEntitlement = VCoresEntitlement
type VCoresDesignEntitlement = VCoresEntitlement
type VCoresProductionEntitlement = VCoresEntitlement
type VpnsEntitlement = AssignedReassignedEntitlement
type VpcsEntitlement = AssignedReassignedEntitlement

func NewLoadBalancerEntitlementWithDefaults() *AssignedReassignedEntitlement {
	return &AssignedReassignedEntitlement{}
}
func NewStaticIpsEntitlementWithDefaults() *AssignedReassignedEntitlement {
	return &AssignedReassignedEntitlement{}
}
func NewVCoresSandboxEntitlementWithDefaults() *VCoresEntitlement { return &VCoresEntitlement{} }
func NewVCoresDesignEntitlementWithDefaults() *VCoresEntitlement  { return &VCoresEntitlement{} }
func NewVCoresProductionEntitlementWithDefaults() *VCoresEntitlement {
	return &VCoresEntitlement{}
}
func NewVpnsEntitlementWithDefaults() *AssignedReassignedEntitlement {
	return &AssignedReassignedEntitlement{}
}
func NewVpcsEntitlementWithDefaults() *AssignedReassignedEntitlement {
	return &AssignedReassignedEntitlement{}
}
