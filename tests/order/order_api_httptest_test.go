package tests

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
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/tests/api"
	"github.com/stretchr/testify/suite"
)

type OrderAPIHTTPTest struct { // NOTE: tests service through HTTP requests
	// TODO: local vs db
	Mux *http.ServeMux
}

func NewOrderAPIHTTPTest(mux *http.ServeMux) *OrderAPIHTTPTest {
	// TODO: local vs db
	return &OrderAPIHTTPTest{
		Mux: mux,
	}
}

var (
	_ api.Order = (*OrderAPIHTTPTest)(nil)
)

func TestOrderHTTPTestSuite(t *testing.T) {
	mux := router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))
	t.Run("OrderAPIHTTPTest", func(t *testing.T) {
		suite.Run(t, &OrderSuite{
			API: NewOrderAPIHTTPTest(mux),
		})
	})
}

func (api *OrderAPIHTTPTest) PostOrder(t *testing.T, ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error) {
	t.Helper()
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp := httptest.NewRequest(http.MethodPost, "/v1/mgmt/order", bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK || respHttp.StatusCode == http.StatusCreated {
		var resp response.PostOrderResponse
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

func (api *OrderAPIHTTPTest) GetOrderByID(t *testing.T, ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error) {
	t.Helper()

	reqHttp := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/mgmt/order/%v", req.GetID()), nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetOrderByIDResponse
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

func (api *OrderAPIHTTPTest) GetOrders(t *testing.T, ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error) {
	t.Helper()

	reqHttp := httptest.NewRequest(http.MethodGet, "/v1/mgmt/orders", nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetOrdersResponse
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

func (api *OrderAPIHTTPTest) GetOrderSubOrders(t *testing.T, ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error) {
	t.Helper()

	path := fmt.Sprintf("/v1/mgmt/order/%v/suborders", req.GetID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodGet, path, nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetOrderSubOrdersResponse
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

func (api *OrderAPIHTTPTest) PatchOrder(t *testing.T, ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error) {
	t.Helper()
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp := httptest.NewRequest(http.MethodPatch, "/v1/mgmt/order", bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.PatchOrderResponse
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

func (api *OrderAPIHTTPTest) DeleteOrder(t *testing.T, ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error) {
	t.Helper()

	reqHttp := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/mgmt/order/%v", req.GetID()), nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusNoContent {
		return &response.DeleteOrderResponse{}, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}

////

func (api *OrderAPIHTTPTest) PutDelegatedTask(t *testing.T, ctx context.Context, req *request.PutDelegatedTaskRequest) (*response.PutDelegatedTaskResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	path := fmt.Sprintf("/v1/mgmt/order/%v/delegated_task", req.GetOrderID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.PutDelegatedTaskResponse
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

func (api *OrderAPIHTTPTest) PatchDelegatedTask(t *testing.T, ctx context.Context, req *request.PatchDelegatedTaskRequest) (*response.PatchDelegatedTaskResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	path := fmt.Sprintf("/v1/mgmt/order/%v/delegated_task", req.GetOrderID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodPatch, path, bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.PatchDelegatedTaskResponse
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

func (api *OrderAPIHTTPTest) DeleteDelegatedTask(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTaskRequest) (*response.DeleteDelegatedTaskResponse, errwrap.Error) {
	t.Helper()

	path := fmt.Sprintf("/v1/mgmt/order/%v/delegated_task/%v", req.GetOrderID(), req.GetDelegatedTaskID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodDelete, path, nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusNoContent {
		return &response.DeleteDelegatedTaskResponse{}, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}

////

func (api *OrderAPIHTTPTest) PutSitRep(t *testing.T, ctx context.Context, req *request.PutSitRepRequest) (*response.PutSitRepResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	path := fmt.Sprintf("/v1/mgmt/order/%v/sitrep", req.GetOrderID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.PutSitRepResponse
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

func (api *OrderAPIHTTPTest) PatchSitRep(t *testing.T, ctx context.Context, req *request.PatchSitRepRequest) (*response.PatchSitRepResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	path := fmt.Sprintf("/v1/mgmt/order/%v/sitrep", req.GetOrderID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodPatch, path, bytes.NewBuffer(reqBytes))
	reqHttp.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.PatchSitRepResponse
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

func (api *OrderAPIHTTPTest) DeleteSitRep(t *testing.T, ctx context.Context, req *request.DeleteSitRepRequest) (*response.DeleteSitRepResponse, errwrap.Error) {
	t.Helper()

	path := fmt.Sprintf("/v1/mgmt/order/%v/sitrep/%v", req.GetOrderID(), req.GetSitRepID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodDelete, path, nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusNoContent {
		return &response.DeleteSitRepResponse{}, nil
	}
	var errw errwrap.Err
	if err := json.NewDecoder(respHttp.Body).Decode(&errw); err != nil {
		rawResponse, _ := httputil.DumpResponse(respHttp, false)
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s\nRaw response: %v", err, string(rawResponse))
	}

	return nil, &errw
}
