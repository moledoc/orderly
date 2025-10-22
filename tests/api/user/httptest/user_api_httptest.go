package httptest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/api"
)

type UserAPIHTTPTest struct { // NOTE: tests service through HTTP requests using httptest
	Mux *http.ServeMux
}

func NewUserAPIHTTPTest(mux *http.ServeMux) *UserAPIHTTPTest {
	// TODO: local vs db
	return &UserAPIHTTPTest{
		Mux: mux,
	}
}

var (
	_ api.User = (*UserAPIHTTPTest)(nil)
)

func (api *UserAPIHTTPTest) PostUser(t *testing.T, ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp := httptest.NewRequest(http.MethodPost, "/v1/mgmt/user", bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
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

func (api *UserAPIHTTPTest) GetUserByID(t *testing.T, ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error) {
	t.Helper()

	reqHttp := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/mgmt/user/%v", req.GetID()), nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
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

func (api *UserAPIHTTPTest) GetUsers(t *testing.T, ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error) {
	t.Helper()

	reqHttp := httptest.NewRequest(http.MethodGet, "/v1/mgmt/users", nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
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

func (api *UserAPIHTTPTest) GetUserSubOrdinates(t *testing.T, ctx context.Context, req *request.GetUserSubOrdinatesRequest) (*response.GetUserSubOrdinatesResponse, errwrap.Error) {
	t.Helper()

	path := fmt.Sprintf("/v1/mgmt/user/%v/subordinates", req.GetID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodGet, path, nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetUserSubOrdinatesResponse
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

func (api *UserAPIHTTPTest) PatchUser(t *testing.T, ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp := httptest.NewRequest(http.MethodPatch, "/v1/mgmt/user", bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
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

func (api *UserAPIHTTPTest) DeleteUser(t *testing.T, ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error) {
	t.Helper()

	reqHttp := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/mgmt/user/%v", req.GetID()), nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
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
