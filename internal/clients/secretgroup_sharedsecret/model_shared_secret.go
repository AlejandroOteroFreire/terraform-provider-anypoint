package secretgroup_sharedsecret

import "encoding/json"

type SharedSecret struct {
	Name            *string          `json:"name,omitempty"`
	Type            *string          `json:"type,omitempty"`
	ExpirationDate  *string          `json:"expirationDate,omitempty"`
	Username        *string          `json:"username,omitempty"`
	Password        *string          `json:"password,omitempty"`
	AccessKeyId     *string          `json:"accessKeyId,omitempty"`
	SecretAccessKey *string          `json:"secretAccessKey,omitempty"`
	Key             *string          `json:"key,omitempty"`
	Content         *string          `json:"content,omitempty"`
	Meta            *SharedSecretMeta `json:"meta,omitempty"`
}

type SharedSecretMeta struct {
	Id *string `json:"id,omitempty"`
}

func NewSharedSecret() *SharedSecret { return &SharedSecret{} }

func (o *SharedSecret) GetId() string {
	if o == nil || o.Meta == nil || o.Meta.Id == nil {
		return ""
	}
	return *o.Meta.Id
}

func (o *SharedSecret) GetName() string {
	if o == nil || o.Name == nil {
		return ""
	}
	return *o.Name
}

func (o *SharedSecret) GetType() string {
	if o == nil || o.Type == nil {
		return ""
	}
	return *o.Type
}

func (o *SharedSecret) GetExpirationDate() string {
	if o == nil || o.ExpirationDate == nil {
		return ""
	}
	return *o.ExpirationDate
}

func (o *SharedSecret) GetUsername() string {
	if o == nil || o.Username == nil {
		return ""
	}
	return *o.Username
}

func (o *SharedSecret) GetAccessKeyId() string {
	if o == nil || o.AccessKeyId == nil {
		return ""
	}
	return *o.AccessKeyId
}

func (o SharedSecret) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	if o.ExpirationDate != nil {
		toSerialize["expirationDate"] = o.ExpirationDate
	}
	if o.Username != nil {
		toSerialize["username"] = o.Username
	}
	if o.Password != nil {
		toSerialize["password"] = o.Password
	}
	if o.AccessKeyId != nil {
		toSerialize["accessKeyId"] = o.AccessKeyId
	}
	if o.SecretAccessKey != nil {
		toSerialize["secretAccessKey"] = o.SecretAccessKey
	}
	if o.Key != nil {
		toSerialize["key"] = o.Key
	}
	if o.Content != nil {
		toSerialize["content"] = o.Content
	}
	return json.Marshal(toSerialize)
}

type SharedSecretListResponse struct {
	SharedSecrets []SharedSecret `json:"sharedSecrets,omitempty"`
}

func (o *SharedSecretListResponse) GetSharedSecrets() []SharedSecret {
	if o == nil {
		return []SharedSecret{}
	}
	return o.SharedSecrets
}
