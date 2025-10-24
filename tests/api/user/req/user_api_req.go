package req

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/api"
)

type UserAPIReq struct { // NOTE: tests service through HTTP requests
	// TODO: local vs db
	HttpClient *http.Client
	BaseURL    string
}

func NewUserAPIReq() *UserAPIReq {
	// TODO: local vs db
	return &UserAPIReq{
		HttpClient: &http.Client{},
		BaseURL:    "http://localhost:8080",
	}
}

var (
	_ api.User = (*UserAPIReq)(nil)
)

func (api *UserAPIReq) PostUser(t *testing.T, ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	respHttp, err := api.HttpClient.Post(fmt.Sprintf("%s/v1/mgmt/user", api.BaseURL), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK || respHttp.StatusCode == http.StatusCreated {
		var resp response.PostUserResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}

func (api *UserAPIReq) GetUserByID(t *testing.T, ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error) {
	t.Helper()

	respHttp, err := api.HttpClient.Get(fmt.Sprintf("%s/v1/mgmt/user/%v", api.BaseURL, req.GetID()))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetUserByIDResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}

func (api *UserAPIReq) GetUsers(t *testing.T, ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error) {
	t.Helper()

	baseURL, _ := url.Parse(fmt.Sprintf("%s/v1/mgmt/users", api.BaseURL))
	params := url.Values{}
	if len(req.GetEmails()) > 0 {
		for _, em := range req.GetEmails() {
			params.Add("emails", string(em))
		}
	}
	if len(req.GetSupervisor()) > 0 {
		params.Add("supervisor", string(req.GetSupervisor()))
	}
	baseURL.RawQuery = params.Encode()

	respHttp, err := api.HttpClient.Get(baseURL.String())
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetUsersResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}

func (api *UserAPIReq) PatchUser(t *testing.T, ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error) {
	t.Helper()
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/v1/mgmt/user", api.BaseURL), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "new request failed: %s", err)
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	respHttp, err := api.HttpClient.Do(reqHttp)
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.PatchUserResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}

func (api *UserAPIReq) DeleteUser(t *testing.T, ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error) {
	t.Helper()

	reqHttp, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/mgmt/user/%v", api.BaseURL, req.GetID()), nil)
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "new request failed: %s", err)
	}

	respHttp, err := api.HttpClient.Do(reqHttp)
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusNoContent {
		return &response.DeleteUserResponse{}, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}
