package secretgroup_certificatepinset

type CertificatePinset struct {
	Name           *string               `json:"name,omitempty"`
	ExpirationDate *string               `json:"expirationDate,omitempty"`
	Algorithm      *string               `json:"algorithm,omitempty"`
	Meta           *CertificatePinsetMeta `json:"meta,omitempty"`
}

type CertificatePinsetMeta struct {
	Id *string `json:"id,omitempty"`
}

func NewCertificatePinset() *CertificatePinset { return &CertificatePinset{} }

func (o *CertificatePinset) GetId() string {
	if o == nil || o.Meta == nil || o.Meta.Id == nil {
		return ""
	}
	return *o.Meta.Id
}

func (o *CertificatePinset) GetName() string {
	if o == nil || o.Name == nil {
		return ""
	}
	return *o.Name
}

func (o *CertificatePinset) GetExpirationDate() string {
	if o == nil || o.ExpirationDate == nil {
		return ""
	}
	return *o.ExpirationDate
}

func (o *CertificatePinset) GetAlgorithm() string {
	if o == nil || o.Algorithm == nil {
		return ""
	}
	return *o.Algorithm
}
