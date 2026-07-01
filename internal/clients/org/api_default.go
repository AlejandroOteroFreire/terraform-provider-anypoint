package org

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type DefaultApiService service

// --- POST /organizations (CreateBG) ---

type DefaultApiCreateBGRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	body       *BGPostReqBody
}

func (r DefaultApiCreateBGRequest) BGPostReqBody(body BGPostReqBody) DefaultApiCreateBGRequest {
	r.body = &body
	return r
}

func (r DefaultApiCreateBGRequest) Execute() (*BGCore, *http.Response, error) {
	return r.ApiService.CreateBGExecute(r)
}

func (a *DefaultApiService) CreateBG(ctx context.Context) DefaultApiCreateBGRequest {
	return DefaultApiCreateBGRequest{ApiService: a, ctx: ctx}
}

func (a *DefaultApiService) CreateBGExecute(r DefaultApiCreateBGRequest) (*BGCore, *http.Response, error) {
	var result *BGCore
	baseURL, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.CreateBG")
	if err != nil {
		return nil, nil, &GenericOpenAPIError{error: err.Error()}
	}
	jsonBody, err := json.Marshal(r.body)
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequestWithContext(r.ctx, http.MethodPost, baseURL+"/organizations", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	setAuthHeader(r.ctx, req)

	resp, err := a.doAndDecode(req, &result)
	return result, resp, err
}

// --- GET /organizations/{orgId} (GetBG) ---

type DefaultApiGetBGRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
}

func (r DefaultApiGetBGRequest) Execute() (*MasterBGDetail, *http.Response, error) {
	return r.ApiService.GetBGExecute(r)
}

func (a *DefaultApiService) GetBG(ctx context.Context, orgId string) DefaultApiGetBGRequest {
	return DefaultApiGetBGRequest{ApiService: a, ctx: ctx, orgId: orgId}
}

func (a *DefaultApiService) GetBGExecute(r DefaultApiGetBGRequest) (*MasterBGDetail, *http.Response, error) {
	var result *MasterBGDetail
	baseURL, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.GetBG")
	if err != nil {
		return nil, nil, &GenericOpenAPIError{error: err.Error()}
	}
	path := baseURL + "/organizations/" + strings.TrimSpace(r.orgId)
	req, err := http.NewRequestWithContext(r.ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	setAuthHeader(r.ctx, req)

	resp, err := a.doAndDecode(req, &result)
	return result, resp, err
}

// --- PUT /organizations/{orgId} (UpdateBG) ---

type DefaultApiUpdateBGRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
	body       *BGUpdateReqBody
}

func (r DefaultApiUpdateBGRequest) BGUpdateReqBody(body BGUpdateReqBody) DefaultApiUpdateBGRequest {
	r.body = &body
	return r
}

func (r DefaultApiUpdateBGRequest) Execute() (*BGCore, *http.Response, error) {
	return r.ApiService.UpdateBGExecute(r)
}

func (a *DefaultApiService) UpdateBG(ctx context.Context, orgId string) DefaultApiUpdateBGRequest {
	return DefaultApiUpdateBGRequest{ApiService: a, ctx: ctx, orgId: orgId}
}

func (a *DefaultApiService) UpdateBGExecute(r DefaultApiUpdateBGRequest) (*BGCore, *http.Response, error) {
	var result *BGCore
	baseURL, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.UpdateBG")
	if err != nil {
		return nil, nil, &GenericOpenAPIError{error: err.Error()}
	}
	jsonBody, err := json.Marshal(r.body)
	if err != nil {
		return nil, nil, err
	}
	path := baseURL + "/organizations/" + strings.TrimSpace(r.orgId)
	req, err := http.NewRequestWithContext(r.ctx, http.MethodPut, path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	setAuthHeader(r.ctx, req)

	resp, err := a.doAndDecode(req, &result)
	return result, resp, err
}

// --- DELETE /organizations/{orgId} (DeleteBG) ---

type DefaultApiDeleteBGRequest struct {
	ctx        context.Context
	ApiService *DefaultApiService
	orgId      string
}

func (r DefaultApiDeleteBGRequest) Execute() (map[string]interface{}, *http.Response, error) {
	return r.ApiService.DeleteBGExecute(r)
}

func (a *DefaultApiService) DeleteBG(ctx context.Context, orgId string) DefaultApiDeleteBGRequest {
	return DefaultApiDeleteBGRequest{ApiService: a, ctx: ctx, orgId: orgId}
}

func (a *DefaultApiService) DeleteBGExecute(r DefaultApiDeleteBGRequest) (map[string]interface{}, *http.Response, error) {
	baseURL, err := a.client.cfg.ServerURLWithContext(r.ctx, "DefaultApiService.DeleteBG")
	if err != nil {
		return nil, nil, &GenericOpenAPIError{error: err.Error()}
	}
	path := baseURL + "/organizations/" + strings.TrimSpace(r.orgId)
	req, err := http.NewRequestWithContext(r.ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}
	setAuthHeader(r.ctx, req)

	httpClient := a.client.cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, resp, err
	}
	// Keep the body open/re-readable for callers that inspect it on error, matching
	// the convention every other client in this repo already follows.
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(b))
	if resp.StatusCode >= 300 {
		return nil, resp, &GenericOpenAPIError{body: b, error: resp.Status}
	}
	return nil, resp, nil
}

// --- shared helpers ---

func setAuthHeader(ctx context.Context, req *http.Request) {
	if ctx == nil {
		return
	}
	if token, ok := ctx.Value(ContextAccessToken).(string); ok && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// doAndDecode executes req, leaves the response body re-readable (so callers can
// still inspect it on error the same way the rest of this repo's clients do), and
// JSON-decodes it into out on success.
func (a *DefaultApiService) doAndDecode(req *http.Request, out interface{}) (*http.Response, error) {
	httpClient := a.client.cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(b))

	if resp.StatusCode >= 300 {
		return resp, &GenericOpenAPIError{body: b, error: resp.Status}
	}
	if len(b) > 0 {
		if err := json.Unmarshal(b, out); err != nil {
			return resp, &GenericOpenAPIError{body: b, error: err.Error()}
		}
	}
	return resp, nil
}
