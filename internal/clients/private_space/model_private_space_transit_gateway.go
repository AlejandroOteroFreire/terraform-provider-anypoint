package private_space

import "encoding/json"

// PrivateSpaceTransitGateway representa un Transit Gateway asociado a un Private Space.
// Body de creación: {"name":"...","resourceShareId":"...","resourceShareAccount":"...","routes":["10.0.0.0/8"]}
type PrivateSpaceTransitGateway struct {
	Id                   *string  `json:"id,omitempty"`
	Name                 *string  `json:"name,omitempty"`
	ResourceShareId      *string  `json:"resourceShareId,omitempty"`
	ResourceShareAccount *string  `json:"resourceShareAccount,omitempty"`
	Routes               []string `json:"routes,omitempty"`
	// Status fields — read-only, returned by GET
	Status     *string `json:"status,omitempty"`
	Attachment *string `json:"attachment,omitempty"`
	Region     *string `json:"region,omitempty"`
}

func NewPrivateSpaceTransitGateway() *PrivateSpaceTransitGateway {
	return &PrivateSpaceTransitGateway{}
}

func (o *PrivateSpaceTransitGateway) GetId() string {
	if o == nil || o.Id == nil {
		return ""
	}
	return *o.Id
}

func (o *PrivateSpaceTransitGateway) GetName() string {
	if o == nil || o.Name == nil {
		return ""
	}
	return *o.Name
}

func (o *PrivateSpaceTransitGateway) GetResourceShareId() string {
	if o == nil || o.ResourceShareId == nil {
		return ""
	}
	return *o.ResourceShareId
}

func (o *PrivateSpaceTransitGateway) GetResourceShareAccount() string {
	if o == nil || o.ResourceShareAccount == nil {
		return ""
	}
	return *o.ResourceShareAccount
}

func (o *PrivateSpaceTransitGateway) GetRoutes() []string {
	if o == nil {
		return []string{}
	}
	return o.Routes
}

func (o *PrivateSpaceTransitGateway) GetStatus() string {
	if o == nil || o.Status == nil {
		return ""
	}
	return *o.Status
}

func (o *PrivateSpaceTransitGateway) GetAttachment() string {
	if o == nil || o.Attachment == nil {
		return ""
	}
	return *o.Attachment
}

func (o *PrivateSpaceTransitGateway) GetRegion() string {
	if o == nil || o.Region == nil {
		return ""
	}
	return *o.Region
}

func (o PrivateSpaceTransitGateway) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.ResourceShareId != nil {
		toSerialize["resourceShareId"] = o.ResourceShareId
	}
	if o.ResourceShareAccount != nil {
		toSerialize["resourceShareAccount"] = o.ResourceShareAccount
	}
	if len(o.Routes) > 0 {
		toSerialize["routes"] = o.Routes
	}
	return json.Marshal(toSerialize)
}
