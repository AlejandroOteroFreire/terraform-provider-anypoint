package private_space

import "encoding/json"

type PrivateSpaceAssociationRequest struct {
	Associations []PrivateSpaceAssociationItem `json:"associations,omitempty"`
}

type PrivateSpaceAssociationItem struct {
	OrganizationId string `json:"organizationId"`
	Environment    string `json:"environment"`
}

type PrivateSpaceAssociation struct {
	Id             *string `json:"id,omitempty"`
	EnvironmentId  *string `json:"environmentId,omitempty"`
	OrganizationId *string `json:"organizationId,omitempty"`
}

func NewPrivateSpaceAssociationRequest() *PrivateSpaceAssociationRequest {
	return &PrivateSpaceAssociationRequest{}
}

func (o *PrivateSpaceAssociation) GetId() string {
	if o == nil || o.Id == nil {
		return ""
	}
	return *o.Id
}

func (o *PrivateSpaceAssociation) GetEnvironmentId() string {
	if o == nil || o.EnvironmentId == nil {
		return ""
	}
	return *o.EnvironmentId
}

func (o *PrivateSpaceAssociation) GetOrganizationId() string {
	if o == nil || o.OrganizationId == nil {
		return ""
	}
	return *o.OrganizationId
}

func (o PrivateSpaceAssociationRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{
		"associations": o.Associations,
	}
	return json.Marshal(toSerialize)
}
