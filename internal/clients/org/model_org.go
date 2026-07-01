package org

// Owner models the organization owner sub-object of MasterBGDetail.
type Owner struct {
	Id                      *string `json:"id,omitempty"`
	CreatedAt               *string `json:"createdAt,omitempty"`
	UpdatedAt               *string `json:"updatedAt,omitempty"`
	OrganizationId          *string `json:"organizationId,omitempty"`
	FirstName               *string `json:"firstName,omitempty"`
	LastName                *string `json:"lastName,omitempty"`
	Email                   *string `json:"email,omitempty"`
	PhoneNumber             *string `json:"phoneNumber,omitempty"`
	Username                *string `json:"username,omitempty"`
	IdproviderId            *string `json:"idprovider_id,omitempty"`
	Enabled                 *bool   `json:"enabled,omitempty"`
	Deleted                 *bool   `json:"deleted,omitempty"`
	LastLogin               *string `json:"lastLogin,omitempty"`
	MfaVerificationExcluded *bool   `json:"mfaVerificationExcluded,omitempty"`
	MfaVerifiersConfigured  *string `json:"mfaVerifiersConfigured,omitempty"`
	Type                    *string `json:"type,omitempty"`
}

func (o Owner) GetId() string {
	if o.Id == nil {
		return ""
	}
	return *o.Id
}
func (o Owner) GetCreatedAt() string {
	if o.CreatedAt == nil {
		return ""
	}
	return *o.CreatedAt
}
func (o Owner) GetUpdatedAt() string {
	if o.UpdatedAt == nil {
		return ""
	}
	return *o.UpdatedAt
}
func (o Owner) GetOrganizationId() string {
	if o.OrganizationId == nil {
		return ""
	}
	return *o.OrganizationId
}
func (o Owner) GetFirstName() string {
	if o.FirstName == nil {
		return ""
	}
	return *o.FirstName
}
func (o Owner) GetLastName() string {
	if o.LastName == nil {
		return ""
	}
	return *o.LastName
}
func (o Owner) GetEmail() string {
	if o.Email == nil {
		return ""
	}
	return *o.Email
}
func (o Owner) GetPhoneNumber() string {
	if o.PhoneNumber == nil {
		return ""
	}
	return *o.PhoneNumber
}
func (o Owner) GetUsername() string {
	if o.Username == nil {
		return ""
	}
	return *o.Username
}
func (o Owner) GetIdproviderId() string {
	if o.IdproviderId == nil {
		return ""
	}
	return *o.IdproviderId
}
func (o Owner) GetEnabled() bool {
	if o.Enabled == nil {
		return false
	}
	return *o.Enabled
}
func (o Owner) GetDeleted() bool {
	if o.Deleted == nil {
		return false
	}
	return *o.Deleted
}
func (o Owner) GetLastLogin() string {
	if o.LastLogin == nil {
		return ""
	}
	return *o.LastLogin
}
func (o Owner) GetMfaVerificationExcluded() bool {
	if o.MfaVerificationExcluded == nil {
		return false
	}
	return *o.MfaVerificationExcluded
}
func (o Owner) GetMfaVerifiersConfigured() string {
	if o.MfaVerifiersConfigured == nil {
		return ""
	}
	return *o.MfaVerifiersConfigured
}
func (o Owner) GetType() string {
	if o.Type == nil {
		return ""
	}
	return *o.Type
}

// Subscription models the subscription sub-object of MasterBGDetail.
type Subscription struct {
	Category   *string `json:"category,omitempty"`
	Type       *string `json:"type,omitempty"`
	Expiration *string `json:"expiration,omitempty"`
}

func (o Subscription) GetCategory() string {
	if o.Category == nil {
		return ""
	}
	return *o.Category
}
func (o Subscription) GetType() string {
	if o.Type == nil {
		return ""
	}
	return *o.Type
}
func (o Subscription) GetExpiration() string {
	if o.Expiration == nil {
		return ""
	}
	return *o.Expiration
}

// Environment models one entry of MasterBGDetail's environments list.
type Environment struct {
	Id             *string `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	OrganizationId *string `json:"organizationId,omitempty"`
	IsProduction   *bool   `json:"isProduction,omitempty"`
	Type           *string `json:"type,omitempty"`
	ClientId       *string `json:"clientId,omitempty"`
}

func (o Environment) GetId() string {
	if o.Id == nil {
		return ""
	}
	return *o.Id
}
func (o Environment) GetName() string {
	if o.Name == nil {
		return ""
	}
	return *o.Name
}
func (o Environment) GetOrganizationId() string {
	if o.OrganizationId == nil {
		return ""
	}
	return *o.OrganizationId
}
func (o Environment) GetIsProduction() bool {
	if o.IsProduction == nil {
		return false
	}
	return *o.IsProduction
}
func (o Environment) GetType() string {
	if o.Type == nil {
		return ""
	}
	return *o.Type
}
func (o Environment) GetClientId() string {
	if o.ClientId == nil {
		return ""
	}
	return *o.ClientId
}

// MasterBGDetail is the full GET /organizations/{organizationId} response.
type MasterBGDetail struct {
	Id                              *string       `json:"id,omitempty"`
	Name                            *string       `json:"name,omitempty"`
	CreatedAt                       *string       `json:"createdAt,omitempty"`
	UpdatedAt                       *string       `json:"updatedAt,omitempty"`
	OwnerId                         *string       `json:"ownerId,omitempty"`
	ClientId                        *string       `json:"clientId,omitempty"`
	IdproviderId                    *string       `json:"idprovider_id,omitempty"`
	IsFederated                     *bool         `json:"isFederated,omitempty"`
	ParentOrganizationIds           []string      `json:"parentOrganizationIds,omitempty"`
	SubOrganizationIds              []string      `json:"subOrganizationIds,omitempty"`
	TenantOrganizationIds           []string      `json:"tenantOrganizationIds,omitempty"`
	MfaRequired                     *string       `json:"mfaRequired,omitempty"`
	IsAutomaticAdminPromotionExempt *bool         `json:"isAutomaticAdminPromotionExempt,omitempty"`
	Domain                          *string       `json:"domain,omitempty"`
	IsMaster                        *bool         `json:"isMaster,omitempty"`
	Properties                      any           `json:"properties,omitempty"`
	Subscription                    *Subscription `json:"subscription,omitempty"`
	Environments                    []Environment `json:"environments,omitempty"`
	Entitlements                    *Entitlements `json:"entitlements,omitempty"`
	Owner                           *Owner        `json:"owner,omitempty"`
	SessionTimeout                  *int32        `json:"sessionTimeout,omitempty"`
}

func (o *MasterBGDetail) GetId() string {
	if o == nil || o.Id == nil {
		return ""
	}
	return *o.Id
}
func (o *MasterBGDetail) GetName() string {
	if o == nil || o.Name == nil {
		return ""
	}
	return *o.Name
}
func (o *MasterBGDetail) GetCreatedAt() string {
	if o == nil || o.CreatedAt == nil {
		return ""
	}
	return *o.CreatedAt
}
func (o *MasterBGDetail) GetUpdatedAt() string {
	if o == nil || o.UpdatedAt == nil {
		return ""
	}
	return *o.UpdatedAt
}
func (o *MasterBGDetail) GetOwnerId() string {
	if o == nil || o.OwnerId == nil {
		return ""
	}
	return *o.OwnerId
}
func (o *MasterBGDetail) GetClientId() string {
	if o == nil || o.ClientId == nil {
		return ""
	}
	return *o.ClientId
}
func (o *MasterBGDetail) GetIdproviderId() string {
	if o == nil || o.IdproviderId == nil {
		return ""
	}
	return *o.IdproviderId
}
func (o *MasterBGDetail) GetIsFederated() bool {
	if o == nil || o.IsFederated == nil {
		return false
	}
	return *o.IsFederated
}
func (o *MasterBGDetail) GetParentOrganizationIds() []string {
	if o == nil {
		return nil
	}
	return o.ParentOrganizationIds
}
func (o *MasterBGDetail) GetSubOrganizationIds() []string {
	if o == nil {
		return nil
	}
	return o.SubOrganizationIds
}
func (o *MasterBGDetail) GetTenantOrganizationIds() []string {
	if o == nil {
		return nil
	}
	return o.TenantOrganizationIds
}
func (o *MasterBGDetail) GetMfaRequired() string {
	if o == nil || o.MfaRequired == nil {
		return ""
	}
	return *o.MfaRequired
}
func (o *MasterBGDetail) GetIsAutomaticAdminPromotionExempt() bool {
	if o == nil || o.IsAutomaticAdminPromotionExempt == nil {
		return false
	}
	return *o.IsAutomaticAdminPromotionExempt
}
func (o *MasterBGDetail) GetDomain() string {
	if o == nil || o.Domain == nil {
		return ""
	}
	return *o.Domain
}
func (o *MasterBGDetail) GetIsMaster() bool {
	if o == nil || o.IsMaster == nil {
		return false
	}
	return *o.IsMaster
}
func (o *MasterBGDetail) GetProperties() any {
	if o == nil {
		return nil
	}
	return o.Properties
}
func (o *MasterBGDetail) GetSubscription() Subscription {
	if o == nil || o.Subscription == nil {
		return Subscription{}
	}
	return *o.Subscription
}
func (o *MasterBGDetail) GetEnvironments() []Environment {
	if o == nil {
		return nil
	}
	return o.Environments
}
func (o *MasterBGDetail) GetEntitlements() Entitlements {
	if o == nil || o.Entitlements == nil {
		return Entitlements{}
	}
	return *o.Entitlements
}
func (o *MasterBGDetail) GetOwner() Owner {
	if o == nil || o.Owner == nil {
		return Owner{}
	}
	return *o.Owner
}
func (o *MasterBGDetail) GetSessionTimeout() int32 {
	if o == nil || o.SessionTimeout == nil {
		return 0
	}
	return *o.SessionTimeout
}

// BGPostReqBody is the POST /organizations request body.
type BGPostReqBody struct {
	Name                 *string       `json:"name,omitempty"`
	OwnerId              *string       `json:"ownerId,omitempty"`
	ParentOrganizationId *string       `json:"parentOrganizationId,omitempty"`
	Entitlements         *Entitlements `json:"entitlements,omitempty"`
}

func NewBGPostReqBodyWithDefaults() *BGPostReqBody { return &BGPostReqBody{} }

func (o *BGPostReqBody) SetName(v string)                 { o.Name = &v }
func (o *BGPostReqBody) SetOwnerId(v string)              { o.OwnerId = &v }
func (o *BGPostReqBody) SetParentOrganizationId(v string) { o.ParentOrganizationId = &v }
func (o *BGPostReqBody) SetEntitlements(v Entitlements)   { o.Entitlements = &v }

// BGUpdateReqBody is the PUT /organizations/{organizationId} request body.
type BGUpdateReqBody struct {
	Name           *string       `json:"name,omitempty"`
	OwnerId        *string       `json:"ownerId,omitempty"`
	SessionTimeout *int32        `json:"sessionTimeout,omitempty"`
	Entitlements   *Entitlements `json:"entitlements,omitempty"`
}

func NewBGUpdateReqBodyWithDefaults() *BGUpdateReqBody { return &BGUpdateReqBody{} }

func (o *BGUpdateReqBody) SetName(v string)               { o.Name = &v }
func (o *BGUpdateReqBody) SetOwnerId(v string)            { o.OwnerId = &v }
func (o *BGUpdateReqBody) SetSessionTimeout(v int32)      { o.SessionTimeout = &v }
func (o *BGUpdateReqBody) SetEntitlements(v Entitlements) { o.Entitlements = &v }

// BGCore is the response shape for POST /organizations (create) and PUT (update).
type BGCore struct {
	Id *string `json:"id,omitempty"`
}

func (o *BGCore) GetId() string {
	if o == nil || o.Id == nil {
		return ""
	}
	return *o.Id
}
