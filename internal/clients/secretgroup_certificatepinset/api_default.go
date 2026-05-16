package secretgroup_certificatepinset

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type DefaultApiService service

type DefaultApiCreateCertificatePinsetRequest struct {
	ctx           context.Context
	ApiService    *DefaultApiService
	orgId         string
	envId         string
	sgId          string
	name          string
	pinsetBase64  string
}

func (r DefaultApiCreateCertificatePinsetRequest) Name(name string) DefaultApiCreateCertificatePinsetRequest {
	r.name = name
	return r
}

func (r DefaultApiCreateCertificatePinsetRequest) PinsetBase64(b64 string) DefaultApiCreateCertificatePinsetRequest {
	r.pinsetBase64 = b64
	return r
}

func (r DefaultApiCreateCertificatePinsetRequest) Execute() (*CertificatePinset, *http.Response, error) {
	return r.ApiService.CreateCertificatePinsetExecute(r)
}

func (a *DefaultApiService) CreateCertificatePinset(ctx context.Context, orgId, envId, sgId string) DefaultApiCreateCertificatePinsetRequest {
	return DefaultApiCreateCertificatePinsetRequest{ApiService: a, ctx: ctx, orgId: orgId, envId: envId, sgId: sgId}
}

func (a *DefaultApiService) CreateCertificatePinsetExecute(r DefaultApiCreateCertificatePinsetRequest) (*CertificatePinset, *http.Response, error) {
	var localVarReturnValue *CertificatePinset
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.CreateCertificatePinset")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}
	localVarPath := localBasePath + "/organizations/{orgId}/environments/{envId}/secretGroups/{sgId}/certificatePinsets"
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"envId"+"}", url.PathEscape(r.envId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"sgId"+"}", url.PathEscape(r.sgId), -1)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", r.name)
	_ = writer.WriteField("certificatePinset", r.pinsetBase64)
	writer.Close()

	req, err := http.NewRequestWithContext(r.ctx, http.MethodPost, localVarPath, body)
	if err != nil {
		return localVarReturnValue, nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	// Add auth headers from context
	if authHeader, ok := r.ctx.Value("Authorization").(string); ok {
		req.Header.Set("Authorization", authHeader)
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

type DefaultApiGetCertificatePinsetRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
	envId      string
	sgId       string
	pinId      string
}

func (r DefaultApiGetCertificatePinsetRequest) Execute() (*CertificatePinset, *http.Response, error) {
	return r.ApiService.GetCertificatePinsetExecute(r)
}

func (a *DefaultApiService) GetCertificatePinset(ctx context.Context, orgId, envId, sgId, pinId string) DefaultApiGetCertificatePinsetRequest {
	return DefaultApiGetCertificatePinsetRequest{ApiService: a, ctx: ctx, orgId: orgId, envId: envId, sgId: sgId, pinId: pinId}
}

func (a *DefaultApiService) GetCertificatePinsetExecute(r DefaultApiGetCertificatePinsetRequest) (*CertificatePinset, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *CertificatePinset
	)
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.GetCertificatePinset")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}
	localVarPath := localBasePath + "/organizations/{orgId}/environments/{envId}/secretGroups/{sgId}/certificatePinsets/{pinId}"
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"envId"+"}", url.PathEscape(r.envId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"sgId"+"}", url.PathEscape(r.sgId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"pinId"+"}", url.PathEscape(r.pinId), -1)
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

type DefaultApiDeleteCertificatePinsetRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
	envId      string
	sgId       string
	pinId      string
}

func (r DefaultApiDeleteCertificatePinsetRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteCertificatePinsetExecute(r)
}

func (a *DefaultApiService) DeleteCertificatePinset(ctx context.Context, orgId, envId, sgId, pinId string) DefaultApiDeleteCertificatePinsetRequest {
	return DefaultApiDeleteCertificatePinsetRequest{ApiService: a, ctx: ctx, orgId: orgId, envId: envId, sgId: sgId, pinId: pinId}
}

func (a *DefaultApiService) DeleteCertificatePinsetExecute(r DefaultApiDeleteCertificatePinsetRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   interface{}
		formFiles          []formFile
	)
	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.DeleteCertificatePinset")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}
	localVarPath := localBasePath + "/organizations/{orgId}/environments/{envId}/secretGroups/{sgId}/certificatePinsets/{pinId}"
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"envId"+"}", url.PathEscape(r.envId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"sgId"+"}", url.PathEscape(r.sgId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"pinId"+"}", url.PathEscape(r.pinId), -1)
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
