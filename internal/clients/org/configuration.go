/*
Access Management API

Description of the Anypoint Access Management API (organizations/business groups)

API version: 1.0.0
*/

package org

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

func (c contextKey) String() string {
	return "auth " + string(c)
}

var (
	// ContextAccessToken takes a string oauth2 access token as authentication for the request.
	ContextAccessToken = contextKey("accesstoken")

	// ContextServerIndex uses a server configuration from the index.
	ContextServerIndex = contextKey("serverIndex")

	ContextOperationServerIndices = contextKey("serverOperationIndices")

	ContextServerVariables = contextKey("serverVariables")

	ContextOperationServerVariables = contextKey("serverOperationVariables")
)

type ServerVariable struct {
	Description  string
	DefaultValue string
	EnumValues   []string
}

type ServerConfiguration struct {
	URL         string
	Description string
	Variables   map[string]ServerVariable
}

type ServerConfigurations []ServerConfiguration

type Configuration struct {
	Host             string            `json:"host,omitempty"`
	Scheme           string            `json:"scheme,omitempty"`
	DefaultHeader    map[string]string `json:"defaultHeader,omitempty"`
	UserAgent        string            `json:"userAgent,omitempty"`
	Debug            bool              `json:"debug,omitempty"`
	Servers          ServerConfigurations
	OperationServers map[string]ServerConfigurations
	HTTPClient       *http.Client
}

// NewConfiguration returns a new Configuration object.
// Server list matches the region variants used by the org module this replaces
// (github.com/mulesoft-anypoint/anypoint-client-go/org), and by the real
// Access Management API (mulesoft/mulesoft-dx apis/access-management/api.yaml).
func NewConfiguration() *Configuration {
	cfg := &Configuration{
		DefaultHeader: make(map[string]string),
		UserAgent:     "terraform-provider-anypoint/internal-org-client",
		Debug:         false,
		Servers: ServerConfigurations{
			{
				URL:         "https://anypoint.mulesoft.com/accounts/api",
				Description: "Anypoint Cloudhub",
			},
			{
				URL:         "https://eu1.anypoint.mulesoft.com/accounts/api",
				Description: "Anypoint Cloudhub EU",
			},
			{
				URL:         "https://gov.anypoint.mulesoft.com/accounts/api",
				Description: "Anypoint Cloudhub GOV",
			},
		},
		OperationServers: map[string]ServerConfigurations{},
	}
	return cfg
}

func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}

func (sc ServerConfigurations) URL(index int, variables map[string]string) (string, error) {
	if index < 0 || len(sc) <= index {
		return "", fmt.Errorf("index %v out of range %v", index, len(sc)-1)
	}
	server := sc[index]
	url := server.URL
	for name, variable := range server.Variables {
		if value, ok := variables[name]; ok {
			found := len(variable.EnumValues) == 0
			for _, enumValue := range variable.EnumValues {
				if value == enumValue {
					found = true
				}
			}
			if !found {
				return "", fmt.Errorf("the variable %s in the server URL has invalid value %v. Must be %v", name, value, variable.EnumValues)
			}
			url = strings.Replace(url, "{"+name+"}", value, -1)
		} else {
			url = strings.Replace(url, "{"+name+"}", variable.DefaultValue, -1)
		}
	}
	return url, nil
}

func (c *Configuration) ServerURL(index int, variables map[string]string) (string, error) {
	return c.Servers.URL(index, variables)
}

func getServerIndex(ctx context.Context) (int, error) {
	si := ctx.Value(ContextServerIndex)
	if si != nil {
		if index, ok := si.(int); ok {
			return index, nil
		}
		return 0, reportError("Invalid type %T should be int", si)
	}
	return 0, nil
}

func getServerOperationIndex(ctx context.Context, endpoint string) (int, error) {
	osi := ctx.Value(ContextOperationServerIndices)
	if osi != nil {
		if operationIndices, ok := osi.(map[string]int); !ok {
			return 0, reportError("Invalid type %T should be map[string]int", osi)
		} else {
			index, ok := operationIndices[endpoint]
			if ok {
				return index, nil
			}
		}
	}
	return getServerIndex(ctx)
}

func getServerVariables(ctx context.Context) (map[string]string, error) {
	sv := ctx.Value(ContextServerVariables)
	if sv != nil {
		if variables, ok := sv.(map[string]string); ok {
			return variables, nil
		}
		return nil, reportError("ctx value of ContextServerVariables has invalid type %T should be map[string]string", sv)
	}
	return nil, nil
}

func getServerOperationVariables(ctx context.Context, endpoint string) (map[string]string, error) {
	osv := ctx.Value(ContextOperationServerVariables)
	if osv != nil {
		if operationVariables, ok := osv.(map[string]map[string]string); !ok {
			return nil, reportError("ctx value of ContextOperationServerVariables has invalid type %T should be map[string]map[string]string", osv)
		} else {
			variables, ok := operationVariables[endpoint]
			if ok {
				return variables, nil
			}
		}
	}
	return getServerVariables(ctx)
}

func (c *Configuration) ServerURLWithContext(ctx context.Context, endpoint string) (string, error) {
	sc, ok := c.OperationServers[endpoint]
	if !ok {
		sc = c.Servers
	}
	if ctx == nil {
		return sc.URL(0, nil)
	}
	index, err := getServerOperationIndex(ctx, endpoint)
	if err != nil {
		return "", err
	}
	variables, err := getServerOperationVariables(ctx, endpoint)
	if err != nil {
		return "", err
	}
	return sc.URL(index, variables)
}
