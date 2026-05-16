package private_space

import (
	"encoding/json"
	"strconv"
)

type PrivateSpaceConnection struct {
	Id   *string           `json:"id,omitempty"`
	Name *string           `json:"name,omitempty"`
	Vpns []PrivateSpaceVpn `json:"vpns,omitempty"`
}

type PrivateSpaceVpn struct {
	Id                  *string                 `json:"id,omitempty"`
	Name                *string                 `json:"name,omitempty"`
	VpnId               *string                 `json:"vpnId,omitempty"`
	ConnectionId        *string                 `json:"connectionId,omitempty"`
	ConnectionName      *string                 `json:"connectionName,omitempty"`
	VpnConnectionStatus *string                 `json:"vpnConnectionStatus,omitempty"`
	LocalAsn            *int32                  `json:"localAsn,omitempty"`
	RemoteAsn           *int32                  `json:"remoteAsn,omitempty"`
	RemoteIpAddress     *string                 `json:"remoteIpAddress,omitempty"`
	StaticRoutes        []string                `json:"staticRoutes,omitempty"`
	VpnTunnels          []PrivateSpaceVpnTunnel `json:"vpnTunnels,omitempty"`
	VpnTunnelsConfig    []interface{}           `json:"vpnTunnelsConfig,omitempty"`
}

type PrivateSpaceVpnTunnel struct {
	// Configurable
	Psk           *string `json:"psk,omitempty"`
	PtpCidr       *string `json:"ptpCidr,omitempty"`
	StartupAction *string `json:"startupAction,omitempty"`
	IsLogsEnabled *bool   `json:"isLogsEnabled,omitempty"`
	// Read-only status fields
	LocalExternalIpAddress *string `json:"localExternalIpAddress,omitempty"`
	LocalPtpIpAddress      *string `json:"localPtpIpAddress,omitempty"`
	RemotePtpIpAddress     *string `json:"remotePtpIpAddress,omitempty"`
	AcceptedRouteCount     *int32  `json:"acceptedRouteCount,omitempty"`
	LastStatusChange       *string `json:"lastStatusChange,omitempty"`
	Status                 *string `json:"status,omitempty"`
	StatusMessage          *string `json:"statusMessage,omitempty"`
	RekeyMarginInSeconds   *int32  `json:"rekeyMarginInSeconds"`
	RekeyFuzz              *int32  `json:"rekeyFuzz"`
}

func NewPrivateSpaceConnection() *PrivateSpaceConnection {
	return &PrivateSpaceConnection{}
}

func (o *PrivateSpaceConnection) GetId() string {
	if o == nil || o.Id == nil {
		return ""
	}
	return *o.Id
}

func (o *PrivateSpaceConnection) GetName() string {
	if o == nil || o.Name == nil {
		return ""
	}
	return *o.Name
}

func (o *PrivateSpaceConnection) GetVpns() []PrivateSpaceVpn {
	if o == nil {
		return []PrivateSpaceVpn{}
	}
	return o.Vpns
}

func (o *PrivateSpaceVpn) GetId() string {
	if o == nil || o.Id == nil {
		return ""
	}
	return *o.Id
}

func (o *PrivateSpaceVpn) GetName() string {
	if o == nil || o.Name == nil {
		return ""
	}
	return *o.Name
}

func (o *PrivateSpaceVpn) GetVpnId() string {
	if o == nil || o.VpnId == nil {
		return ""
	}
	return *o.VpnId
}

func (o *PrivateSpaceVpn) GetConnectionId() string {
	if o == nil || o.ConnectionId == nil {
		return ""
	}
	return *o.ConnectionId
}

func (o *PrivateSpaceVpn) GetConnectionName() string {
	if o == nil || o.ConnectionName == nil {
		return ""
	}
	return *o.ConnectionName
}

func (o *PrivateSpaceVpn) GetVpnConnectionStatus() string {
	if o == nil || o.VpnConnectionStatus == nil {
		return ""
	}
	return *o.VpnConnectionStatus
}

func (o *PrivateSpaceVpn) GetLocalAsn() int32 {
	if o == nil || o.LocalAsn == nil {
		return 0
	}
	return *o.LocalAsn
}

func (o *PrivateSpaceVpn) GetRemoteAsn() int32 {
	if o == nil || o.RemoteAsn == nil {
		return 0
	}
	return *o.RemoteAsn
}

func (o *PrivateSpaceVpn) GetRemoteIpAddress() string {
	if o == nil || o.RemoteIpAddress == nil {
		return ""
	}
	return *o.RemoteIpAddress
}

func (o *PrivateSpaceVpn) GetStaticRoutes() []string {
	if o == nil {
		return []string{}
	}
	return o.StaticRoutes
}

func (o *PrivateSpaceVpn) GetVpnTunnels() []PrivateSpaceVpnTunnel {
	if o == nil {
		return []PrivateSpaceVpnTunnel{}
	}
	return o.VpnTunnels
}

func (o *PrivateSpaceVpnTunnel) GetPsk() string {
	if o == nil || o.Psk == nil {
		return ""
	}
	return *o.Psk
}

func (o *PrivateSpaceVpnTunnel) GetPtpCidr() string {
	if o == nil || o.PtpCidr == nil {
		return ""
	}
	return *o.PtpCidr
}

func (o *PrivateSpaceVpnTunnel) GetStartupAction() string {
	if o == nil || o.StartupAction == nil {
		return ""
	}
	return *o.StartupAction
}

func (o *PrivateSpaceVpnTunnel) GetIsLogsEnabled() bool {
	if o == nil || o.IsLogsEnabled == nil {
		return false
	}
	return *o.IsLogsEnabled
}

func (o *PrivateSpaceVpnTunnel) GetLocalExternalIpAddress() string {
	if o == nil || o.LocalExternalIpAddress == nil {
		return ""
	}
	return *o.LocalExternalIpAddress
}

func (o *PrivateSpaceVpnTunnel) GetLocalPtpIpAddress() string {
	if o == nil || o.LocalPtpIpAddress == nil {
		return ""
	}
	return *o.LocalPtpIpAddress
}

func (o *PrivateSpaceVpnTunnel) GetRemotePtpIpAddress() string {
	if o == nil || o.RemotePtpIpAddress == nil {
		return ""
	}
	return *o.RemotePtpIpAddress
}

func (o *PrivateSpaceVpnTunnel) GetAcceptedRouteCount() int32 {
	if o == nil || o.AcceptedRouteCount == nil {
		return 0
	}
	return *o.AcceptedRouteCount
}

func (o *PrivateSpaceVpnTunnel) GetLastStatusChange() string {
	if o == nil || o.LastStatusChange == nil {
		return ""
	}
	return *o.LastStatusChange
}

func (o *PrivateSpaceVpnTunnel) GetStatus() string {
	if o == nil || o.Status == nil {
		return ""
	}
	return *o.Status
}

func (o *PrivateSpaceVpnTunnel) GetStatusMessage() string {
	if o == nil || o.StatusMessage == nil {
		return ""
	}
	return *o.StatusMessage
}

func (o *PrivateSpaceVpnTunnel) GetRekeyMarginInSeconds() int32 {
	if o == nil || o.RekeyMarginInSeconds == nil {
		return 0
	}
	return *o.RekeyMarginInSeconds
}

func (o *PrivateSpaceVpnTunnel) GetRekeyFuzz() int32 {
	if o == nil || o.RekeyFuzz == nil {
		return 0
	}
	return *o.RekeyFuzz
}

func (o PrivateSpaceConnection) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if len(o.Vpns) > 0 {
		toSerialize["vpns"] = o.Vpns
	}
	return json.Marshal(toSerialize)
}

// MarshalJSON serializes ASN fields as strings to match the API's POST format.
// The GET response returns ASN as integers, but POST expects strings like "64512".
func (o PrivateSpaceVpn) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.RemoteIpAddress != nil {
		toSerialize["remoteIpAddress"] = o.RemoteIpAddress
	}
	if o.LocalAsn != nil {
		toSerialize["localAsn"] = strconv.Itoa(int(*o.LocalAsn))
	}
	if o.RemoteAsn != nil {
		toSerialize["remoteAsn"] = strconv.Itoa(int(*o.RemoteAsn))
	}
	if len(o.StaticRoutes) > 0 {
		toSerialize["staticRoutes"] = o.StaticRoutes
	}
	if len(o.VpnTunnels) > 0 {
		toSerialize["vpnTunnels"] = o.VpnTunnels
	}
	return json.Marshal(toSerialize)
}
