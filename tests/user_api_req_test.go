package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/stretchr/testify/suite"
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
		BaseURL:    "http://127.0.0.1:8080",
	}
}

var (
	_ UserAPI = (*UserAPIReq)(nil)
)

func TestUserReqSuite(t *testing.T) {
	t.Run("UserAPIReq", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: NewUserAPIReq(),
		})
	})
}

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

	if respHttp.StatusCode == http.StatusOK || respHttp.StatusCode == http.StatusCreated {
		var resp response.PostUserResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
	}

	return nil, &errw
}

func (api *UserAPIReq) GetUserByID(t *testing.T, ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error) {
	t.Helper()

	respHttp, err := api.HttpClient.Get(fmt.Sprintf("%s/v1/mgmt/user/%v", api.BaseURL, req.GetID()))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetUserByIDResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
	}

	return nil, &errw
}

func (api *UserAPIReq) GetUsers(t *testing.T, ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error) {
	t.Helper()

	respHttp, err := api.HttpClient.Get(fmt.Sprintf("%s/v1/mgmt/users", api.BaseURL))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetUsersResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
	}

	return nil, &errw
}

func (api *UserAPIReq) GetUserSubOrdinates(t *testing.T, ctx context.Context, req *request.GetUserSubOrdinatesRequest) (*response.GetUserSubOrdinatesResponse, errwrap.Error) {
	t.Helper()

	respHttp, err := api.HttpClient.Get(fmt.Sprintf("%s/v1/mgmt/user/%v/subordinates", api.BaseURL, req.GetID()))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetUserSubOrdinatesResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
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

	if respHttp.StatusCode == http.StatusOK {
		var resp response.PatchUserResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
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

	if respHttp.StatusCode == http.StatusOK || respHttp.StatusCode == http.StatusNoContent {
		var resp response.DeleteUserResponse
		if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
			return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
		}
		return &resp, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
	}

	return nil, &errw
}
