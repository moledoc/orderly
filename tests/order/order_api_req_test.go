package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))
		wg.Done()
		http.ListenAndServe(":8080", nil)
	}()
	wg.Wait()

	t.Run("OrderAPIReq", func(t *testing.T) {
		suite.Run(t, &OrderSuite{
			API: NewOrderAPIReq(),
		})
	})
}

func TestOrderReqPerformanceSuite(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))
		wg.Done()
		http.ListenAndServe(":8080", nil)
	}()
	wg.Wait()

	t.Run("OrderAPIReqPerformance", func(t *testing.T) {
		suite.Run(t, &OrderPerformanceSuite{
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

func (api *OrderAPIReq) PutDelegatedTasks(t *testing.T, ctx context.Context, req *request.PutDelegatedTasksRequest) (*response.PutDelegatedTasksResponse, errwrap.Error) {
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

func (api *OrderAPIReq) PatchDelegatedTasks(t *testing.T, ctx context.Context, req *request.PatchDelegatedTasksRequest) (*response.PatchDelegatedTasksResponse, errwrap.Error) {
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

func (api *OrderAPIReq) DeleteDelegatedTasks(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTasksRequest) (*response.DeleteDelegatedTasksResponse, errwrap.Error) {
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

func (api *OrderAPIReq) PutSitReps(t *testing.T, ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error) {
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

func (api *OrderAPIReq) PatchSitReps(t *testing.T, ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error) {
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

func (api *OrderAPIReq) DeleteSitReps(t *testing.T, ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error) {
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
