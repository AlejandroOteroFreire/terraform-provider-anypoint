package secretgroup_sharedsecret

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DefaultApiService service

type DefaultApiCreateSharedSecretRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
	envId      string
	sgId       string
	body       *SharedSecret
}

func (r DefaultApiCreateSharedSecretRequest) SharedSecretBody(body SharedSecret) DefaultApiCreateSharedSecretRequest {
	r.body = &body
	return r
}

func (r DefaultApiCreateSharedSecretRequest) Execute() (*SharedSecret, *http.Response, error) {
	return r.ApiService.CreateSharedSecretExecute(r)
}

func (a *DefaultApiService) CreateSharedSecret(ctx context.Context, orgId, envId, sgId string) DefaultApiCreateSharedSecretRequest {
	return DefaultApiCreateSharedSecretRequest{ApiService: a, ctx: ctx, orgId: orgId, envId: envId, sgId: sgId}
}

func (a *DefaultApiService) CreateSharedSecretExecute(r DefaultApiCreateSharedSecretRequest) (*SharedSecret, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SharedSecret
	)
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.CreateSharedSecret")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}
	localVarPath := localBasePath + "/organizations/{orgId}/environments/{envId}/secretGroups/{sgId}/sharedSecrets"
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(parameterValueToString(r.orgId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"envId"+"}", url.PathEscape(parameterValueToString(r.envId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"sgId"+"}", url.PathEscape(parameterValueToString(r.sgId, "")), -1)
	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	localVarHTTPContentTypes := []string{"application/json"}
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}
	localVarHeaderParams["Accept"] = "application/json"
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}
	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}
	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}
	if localVarHTTPResponse.StatusCode >= 300 {
		return localVarReturnValue, localVarHTTPResponse, &GenericOpenAPIError{body: localVarBody, error: localVarHTTPResponse.Status}
	}
	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	return localVarReturnValue, localVarHTTPResponse, err
}

type DefaultApiGetSharedSecretRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
	envId      string
	sgId       string
	ssId       string
}

func (r DefaultApiGetSharedSecretRequest) Execute() (*SharedSecret, *http.Response, error) {
	return r.ApiService.GetSharedSecretExecute(r)
}

func (a *DefaultApiService) GetSharedSecret(ctx context.Context, orgId, envId, sgId, ssId string) DefaultApiGetSharedSecretRequest {
	return DefaultApiGetSharedSecretRequest{ApiService: a, ctx: ctx, orgId: orgId, envId: envId, sgId: sgId, ssId: ssId}
}

func (a *DefaultApiService) GetSharedSecretExecute(r DefaultApiGetSharedSecretRequest) (*SharedSecret, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SharedSecret
	)
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.GetSharedSecret")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}
	localVarPath := localBasePath + "/organizations/{orgId}/environments/{envId}/secretGroups/{sgId}/sharedSecrets/{ssId}"
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(parameterValueToString(r.orgId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"envId"+"}", url.PathEscape(parameterValueToString(r.envId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"sgId"+"}", url.PathEscape(parameterValueToString(r.sgId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"ssId"+"}", url.PathEscape(parameterValueToString(r.ssId, "")), -1)
	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	localVarHeaderParams["Accept"] = "application/json"
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}
	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}
	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}
	if localVarHTTPResponse.StatusCode >= 300 {
		return localVarReturnValue, localVarHTTPResponse, &GenericOpenAPIError{body: localVarBody, error: localVarHTTPResponse.Status}
	}
	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	return localVarReturnValue, localVarHTTPResponse, err
}

type DefaultApiUpdateSharedSecretRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
	envId      string
	sgId       string
	ssId       string
	body       *SharedSecret
}

func (r DefaultApiUpdateSharedSecretRequest) SharedSecretBody(body SharedSecret) DefaultApiUpdateSharedSecretRequest {
	r.body = &body
	return r
}

func (r DefaultApiUpdateSharedSecretRequest) Execute() (*SharedSecret, *http.Response, error) {
	return r.ApiService.UpdateSharedSecretExecute(r)
}

func (a *DefaultApiService) UpdateSharedSecret(ctx context.Context, orgId, envId, sgId, ssId string) DefaultApiUpdateSharedSecretRequest {
	return DefaultApiUpdateSharedSecretRequest{ApiService: a, ctx: ctx, orgId: orgId, envId: envId, sgId: sgId, ssId: ssId}
}

func (a *DefaultApiService) UpdateSharedSecretExecute(r DefaultApiUpdateSharedSecretRequest) (*SharedSecret, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SharedSecret
	)
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.UpdateSharedSecret")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}
	localVarPath := localBasePath + "/organizations/{orgId}/environments/{envId}/secretGroups/{sgId}/sharedSecrets/{ssId}"
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(parameterValueToString(r.orgId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"envId"+"}", url.PathEscape(parameterValueToString(r.envId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"sgId"+"}", url.PathEscape(parameterValueToString(r.sgId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"ssId"+"}", url.PathEscape(parameterValueToString(r.ssId, "")), -1)
	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	localVarHTTPContentTypes := []string{"application/json"}
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}
	localVarHeaderParams["Accept"] = "application/json"
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}
	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}
	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}
	if localVarHTTPResponse.StatusCode >= 300 {
		return localVarReturnValue, localVarHTTPResponse, &GenericOpenAPIError{body: localVarBody, error: localVarHTTPResponse.Status}
	}
	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	return localVarReturnValue, localVarHTTPResponse, err
}

type DefaultApiDeleteSharedSecretRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
	envId      string
	sgId       string
	ssId       string
}

func (r DefaultApiDeleteSharedSecretRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteSharedSecretExecute(r)
}

func (a *DefaultApiService) DeleteSharedSecret(ctx context.Context, orgId, envId, sgId, ssId string) DefaultApiDeleteSharedSecretRequest {
	return DefaultApiDeleteSharedSecretRequest{ApiService: a, ctx: ctx, orgId: orgId, envId: envId, sgId: sgId, ssId: ssId}
}

func (a *DefaultApiService) DeleteSharedSecretExecute(r DefaultApiDeleteSharedSecretRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   interface{}
		formFiles          []formFile
	)
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.DeleteSharedSecret")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}
	localVarPath := localBasePath + "/organizations/{orgId}/environments/{envId}/secretGroups/{sgId}/sharedSecrets/{ssId}"
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(parameterValueToString(r.orgId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"envId"+"}", url.PathEscape(parameterValueToString(r.envId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"sgId"+"}", url.PathEscape(parameterValueToString(r.sgId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"ssId"+"}", url.PathEscape(parameterValueToString(r.ssId, "")), -1)
	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}
	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}
	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}
	if localVarHTTPResponse.StatusCode >= 300 {
		return localVarHTTPResponse, &GenericOpenAPIError{body: localVarBody, error: localVarHTTPResponse.Status}
	}
	return localVarHTTPResponse, nil
}
