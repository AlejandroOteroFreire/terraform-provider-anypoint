package apim

import "encoding/json"

type SlaTier struct {
	Id          *int64        `json:"id,omitempty"`
	Name        *string       `json:"name,omitempty"`
	Description *string       `json:"description,omitempty"`
	AutoApprove *bool         `json:"autoApprove,omitempty"`
	Status      *string       `json:"status,omitempty"`
	Limits      []SlaTierLimit `json:"limits,omitempty"`
}

type SlaTierLimit struct {
	TimePeriodInMilliseconds *int64 `json:"timePeriodInMilliseconds,omitempty"`
	MaximumRequests          *int64 `json:"maximumRequests,omitempty"`
	Visible                  *bool  `json:"visible,omitempty"`
}

type SlaTierListResponse struct {
	Total int64     `json:"total,omitempty"`
	Tiers []SlaTier `json:"tiers,omitempty"`
}

func NewSlaTier() *SlaTier { return &SlaTier{} }

func (o *SlaTier) GetId() int64 {
	if o == nil || o.Id == nil {
		return 0
	}
	return *o.Id
}

func (o *SlaTier) GetName() string {
	if o == nil || o.Name == nil {
		return ""
	}
	return *o.Name
}

func (o *SlaTier) GetDescription() string {
	if o == nil || o.Description == nil {
		return ""
	}
	return *o.Description
}

func (o *SlaTier) GetAutoApprove() bool {
	if o == nil || o.AutoApprove == nil {
		return false
	}
	return *o.AutoApprove
}

func (o *SlaTier) GetStatus() string {
	if o == nil || o.Status == nil {
		return ""
	}
	return *o.Status
}

func (o *SlaTier) GetLimits() []SlaTierLimit {
	if o == nil {
		return []SlaTierLimit{}
	}
	return o.Limits
}

func (o *SlaTierLimit) GetTimePeriodInMilliseconds() int64 {
	if o == nil || o.TimePeriodInMilliseconds == nil {
		return 0
	}
	return *o.TimePeriodInMilliseconds
}

func (o *SlaTierLimit) GetMaximumRequests() int64 {
	if o == nil || o.MaximumRequests == nil {
		return 0
	}
	return *o.MaximumRequests
}

func (o *SlaTierLimit) GetVisible() bool {
	if o == nil || o.Visible == nil {
		return true
	}
	return *o.Visible
}

func (o SlaTier) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.AutoApprove != nil {
		toSerialize["autoApprove"] = o.AutoApprove
	}
	if o.Status != nil {
		toSerialize["status"] = o.Status
	}
	if len(o.Limits) > 0 {
		toSerialize["limits"] = o.Limits
	}
	return json.Marshal(toSerialize)
}
