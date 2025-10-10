package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/api"
	"github.com/stretchr/testify/suite"
)

type OrderAPIReq struct { // NOTE: tests service through HTTP requests
	// TODO: local vs db
	HttpClient *http.Client
	BaseURL    string
}

func NewOrderAPIReq() *OrderAPIReq {
	// TODO: local vs db
	return &OrderAPIReq{
		HttpClient: &http.Client{},
		BaseURL:    "http://localhost:8080",
	}
}

var (
	_ api.Order = (*OrderAPIReq)(nil)
)

func TestOrderReqSuite(t *testing.T) {
	t.Run("OrderAPIReq", func(t *testing.T) {
		suite.Run(t, &OrderSuite{
			API: NewOrderAPIReq(),
		})
	})
}

func (api *OrderAPIReq) PostOrder(t *testing.T, ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error) {
	t.Helper()
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}
	respHttp, err := api.HttpClient.Post(fmt.Sprintf("%s/v1/mgmt/order", api.BaseURL), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
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

func (api *OrderAPIReq) GetOrderByID(t *testing.T, ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error) {
	t.Helper()

	respHttp, err := api.HttpClient.Get(fmt.Sprintf("%s/v1/mgmt/order/%v", api.BaseURL, req.GetID()))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
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

func (api *OrderAPIReq) GetOrders(t *testing.T, ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error) {
	t.Helper()

	respHttp, err := api.HttpClient.Get(fmt.Sprintf("%s/v1/mgmt/orders", api.BaseURL))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
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

func (api *OrderAPIReq) GetOrderSubOrders(t *testing.T, ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error) {
	t.Helper()

	respHttp, err := api.HttpClient.Get(fmt.Sprintf("%s/v1/mgmt/order/%v/suborders", api.BaseURL, req.GetID()))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
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

func (api *OrderAPIReq) PatchOrder(t *testing.T, ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error) {
	t.Helper()
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/v1/mgmt/order", api.BaseURL), bytes.NewBuffer(reqBytes))
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

func (api *OrderAPIReq) DeleteOrder(t *testing.T, ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error) {
	t.Helper()

	reqHttp, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/mgmt/order/%v", api.BaseURL, req.GetID()), nil)
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "new request failed: %s", err)
	}

	respHttp, err := api.HttpClient.Do(reqHttp)
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
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

func (api *OrderAPIReq) PutDelegatedTask(t *testing.T, ctx context.Context, req *request.PutDelegatedTaskRequest) (*response.PutDelegatedTaskResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/mgmt/order/%v/delegated_task", api.BaseURL, req.GetOrderID()), bytes.NewBuffer(reqBytes))
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

func (api *OrderAPIReq) PatchDelegatedTask(t *testing.T, ctx context.Context, req *request.PatchDelegatedTaskRequest) (*response.PatchDelegatedTaskResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/v1/mgmt/order/%v/delegated_task", api.BaseURL, req.GetOrderID()), bytes.NewBuffer(reqBytes))
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

func (api *OrderAPIReq) DeleteDelegatedTask(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTaskRequest) (*response.DeleteDelegatedTaskResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/mgmt/order/%v/delegated_task", api.BaseURL, req.GetOrderID()), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "new request failed: %s", err)
	}

	respHttp, err := api.HttpClient.Do(reqHttp)
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
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

func (api *OrderAPIReq) PutSitRep(t *testing.T, ctx context.Context, req *request.PutSitRepRequest) (*response.PutSitRepResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/mgmt/order/%v/sitrep", api.BaseURL, req.GetOrderID()), bytes.NewBuffer(reqBytes))
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

func (api *OrderAPIReq) PatchSitRep(t *testing.T, ctx context.Context, req *request.PatchSitRepRequest) (*response.PatchSitRepResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/v1/mgmt/order/%v/sitrep", api.BaseURL, req.GetOrderID()), bytes.NewBuffer(reqBytes))
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

func (api *OrderAPIReq) DeleteSitRep(t *testing.T, ctx context.Context, req *request.DeleteSitRepRequest) (*response.DeleteSitRepResponse, errwrap.Error) {
	t.Helper()

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}

	reqHttp, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/mgmt/order/%v/sitrep", api.BaseURL, req.GetOrderID()), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "new request failed: %s", err)
	}

	respHttp, err := api.HttpClient.Do(reqHttp)
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
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
