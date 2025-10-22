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

func (api *OrderAPIHTTPTest) PutDelegatedTasks(t *testing.T, ctx context.Context, req *request.PutDelegatedTasksRequest) (*response.PutDelegatedTasksResponse, errwrap.Error) {
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
		var resp response.PutDelegatedTasksResponse
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

func (api *OrderAPIHTTPTest) PatchDelegatedTasks(t *testing.T, ctx context.Context, req *request.PatchDelegatedTasksRequest) (*response.PatchDelegatedTasksResponse, errwrap.Error) {
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
		var resp response.PatchDelegatedTasksResponse
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

func (api *OrderAPIHTTPTest) DeleteDelegatedTasks(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTasksRequest) (*response.DeleteDelegatedTasksResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}
	path := fmt.Sprintf("/v1/mgmt/order/%v/delegated_task", req.GetOrderID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodDelete, path, bytes.NewBuffer(reqBytes))

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.DeleteDelegatedTasksResponse
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

////

func (api *OrderAPIHTTPTest) PutSitReps(t *testing.T, ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error) {
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
		var resp response.PutSitRepsResponse
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

func (api *OrderAPIHTTPTest) PatchSitReps(t *testing.T, ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error) {
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
		var resp response.PatchSitRepsResponse
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

func (api *OrderAPIHTTPTest) DeleteSitReps(t *testing.T, ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	path := fmt.Sprintf("/v1/mgmt/order/%v/sitrep", req.GetOrderID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodDelete, path, bytes.NewBuffer(reqBytes))

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.DeleteSitRepsResponse
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

func (api *OrderAPIHTTPTest) GetUserOrders(t *testing.T, ctx context.Context, req *request.GetUserOrdersRequest) (*response.GetUserOrdersResponse, errwrap.Error) {
	t.Helper()

	path := fmt.Sprintf("/v1/mgmt/user/%v/orders", req.GetUserID())
	path = strings.ReplaceAll(path, "//", "/")
	reqHttp := httptest.NewRequest(http.MethodGet, path, nil)

	rr := httptest.NewRecorder()
	api.Mux.ServeHTTP(rr, reqHttp)
	respHttp := rr.Result()
	defer respHttp.Body.Close()

	if respHttp.StatusCode == http.StatusOK {
		var resp response.GetUserOrdersResponse
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
