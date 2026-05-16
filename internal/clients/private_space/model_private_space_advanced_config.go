package private_space

type PrivateSpaceAdvancedConfig struct {
	IngressConfiguration *PrivateSpaceIngressConfiguration `json:"ingressConfiguration,omitempty"`
	EnableIAMRole        *bool                             `json:"enableIAMRole,omitempty"`
}

type PrivateSpaceIngressConfiguration struct {
	ReadResponseTimeout *string                    `json:"readResponseTimeout,omitempty"`
	Protocol            *string                    `json:"protocol,omitempty"`
	Logs                *PrivateSpaceIngressLogs   `json:"logs,omitempty"`
}

type PrivateSpaceIngressLogs struct {
	PortLogLevel *string                      `json:"portLogLevel,omitempty"`
	Filters      []PrivateSpaceIngressFilter  `json:"filters,omitempty"`
}

type PrivateSpaceIngressFilter struct {
	IP    string `json:"ip"`
	Level string `json:"level"`
}

func NewPrivateSpaceAdvancedConfig() *PrivateSpaceAdvancedConfig {
	return &PrivateSpaceAdvancedConfig{}
}

func (o *PrivateSpaceAdvancedConfig) GetEnableIAMRole() bool {
	if o == nil || o.EnableIAMRole == nil {
		return false
	}
	return *o.EnableIAMRole
}

func (o *PrivateSpaceAdvancedConfig) GetIngressConfiguration() PrivateSpaceIngressConfiguration {
	if o == nil || o.IngressConfiguration == nil {
		return PrivateSpaceIngressConfiguration{}
	}
	return *o.IngressConfiguration
}

func (o *PrivateSpaceIngressConfiguration) GetReadResponseTimeout() string {
	if o == nil || o.ReadResponseTimeout == nil {
		return ""
	}
	return *o.ReadResponseTimeout
}

func (o *PrivateSpaceIngressConfiguration) GetProtocol() string {
	if o == nil || o.Protocol == nil {
		return ""
	}
	return *o.Protocol
}

func (o *PrivateSpaceIngressConfiguration) GetLogs() PrivateSpaceIngressLogs {
	if o == nil || o.Logs == nil {
		return PrivateSpaceIngressLogs{}
	}
	return *o.Logs
}

func (o *PrivateSpaceIngressLogs) GetPortLogLevel() string {
	if o == nil || o.PortLogLevel == nil {
		return ""
	}
	return *o.PortLogLevel
}

func (o *PrivateSpaceIngressLogs) GetFilters() []PrivateSpaceIngressFilter {
	if o == nil {
		return []PrivateSpaceIngressFilter{}
	}
	return o.Filters
}
