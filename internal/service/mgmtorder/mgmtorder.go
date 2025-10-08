package mgmtorder

import (
	"context"
	"net/http"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/pkg/utils"
)

func (s *serviceMgmtOrder) PostOrder(ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PostOrder")
	defer middleware.SpanStop(ctx, "PostOrder")

	if err := ValidatePostOrderRequest(req); err != nil {
		return nil, err
	}

	o := req.GetOrder()
	o.Task.ID = meta.ID(utils.RandAlphanum())

	now := time.Now().UTC()
	o.SetMeta(&meta.Meta{
		Version: 1,
		Created: now,
		Updated: now,
	})

	resp, err := s.Repository.Write(ctx, o)
	if err != nil {
		return nil, err
	}
	return &response.PostOrderResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) GetOrderByID(ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetOrderByID")
	defer middleware.SpanStop(ctx, "GetOrderByID")

	if err := ValidateGetOrderByIDRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.Repository.ReadByID(ctx, req.GetID())
	if err != nil {
		return nil, err
	}
	return &response.GetOrderByIDResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) GetOrders(ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetOrders")
	defer middleware.SpanStop(ctx, "GetOrders")

	if err := ValidateGetOrdersRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.Repository.ReadAll(ctx)
	if err != nil {
		return nil, err
	}
	return &response.GetOrdersResponse{
		Orders: resp,
	}, nil
}

func (s *serviceMgmtOrder) GetOrderSubOrders(ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetOrderSubOrders")
	defer middleware.SpanStop(ctx, "GetOrderSubOrders")

	if err := ValidateGetOrderSubOrdersRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.Repository.ReadSubOrders(ctx, req.GetID())
	if err != nil {
		return nil, err
	}
	return &response.GetOrderSubOrdersResponse{
		SubOrders: resp,
	}, nil
}

func (s *serviceMgmtOrder) PatchOrder(ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PatchOrder")
	defer middleware.SpanStop(ctx, "PatchOrder")

	if err := ValidatePatchOrderRequest(req); err != nil {
		return nil, err
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrder().GetTask().GetID())
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	hasChanges := false
	patchedOrder := order.Clone()

	hasChanges = hasChanges || patchTask(req.GetOrder().GetTask(), patchedOrder.GetTask())

	// TODO: optimize
	for _, reqDelegatedTask := range req.GetOrder().GetDelegatedTasks() {
		for _, patchedDelegatedTask := range patchedOrder.GetDelegatedTasks() {
			if reqDelegatedTask.GetID() == patchedDelegatedTask.GetID() {
				hasChanges = hasChanges || patchTask(patchedDelegatedTask, patchedDelegatedTask)
				break
			}
		}
	}

	// TODO: optimize
	for _, reqSitRep := range req.GetOrder().GetSitReps() {
		for _, patchedSitRep := range patchedOrder.GetSitReps() {
			if reqSitRep.GetID() == patchedSitRep.GetID() {
				hasChanges = hasChanges || patchSitRep(patchedSitRep, patchedSitRep)
				break
			}
		}
	}
	if req.GetOrder().GetParentOrderID() != patchedOrder.GetParentOrderID() {
		patchedOrder.SetParentOrderID(req.GetOrder().GetParentOrderID())
		hasChanges = true
	}

	if !hasChanges { // NOTE: no changes, return current
		return &response.PatchOrderResponse{
			Order: order,
		}, nil
	}

	patchedOrder.GetMeta().SetUpdated(now)
	patchedOrder.GetMeta().VersionIncr()

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PatchOrderResponse{
		Order: resp,
	}, err
}

func (s *serviceMgmtOrder) DeleteOrder(ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "DeleteOrder")
	defer middleware.SpanStop(ctx, "DeleteOrder")

	if err := ValidateDeleteOrderRequest(req); err != nil {
		return nil, err
	}

	return &response.DeleteOrderResponse{}, s.Repository.Delete(ctx, req.GetID())
}

func (s *serviceMgmtOrder) PutDelegatedTask(ctx context.Context, req *request.PutDelegatedTaskRequest) (*response.PutDelegatedTaskResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PutDelegatedTask")
	defer middleware.SpanStop(ctx, "PutDelegatedTask")

	if err := ValidatePutDelegatedTaskRequest(req); err != nil {
		return nil, err
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := order.Clone()

	task := req.GetTask()
	now := time.Now().UTC()
	task.SetID(meta.ID(utils.RandAlphanum()))

	patchedOrder.GetMeta().SetUpdated(now)
	patchedOrder.GetMeta().VersionIncr()

	patchedOrder.SetDelegatedTasks(append(patchedOrder.GetDelegatedTasks(), task))

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PutDelegatedTaskResponse{
		Order: resp,
	}, nil
}

func patchTask(reqTask *order.Task, patchedTask *order.Task) bool {
	hasChanges := false

	if reqTask.GetState() != patchedTask.GetState() {
		patchedTask.SetState(reqTask.GetState())
		hasChanges = true
	}
	if reqTask.GetAccountable() != patchedTask.GetAccountable() {
		patchedTask.SetAccountable(reqTask.GetAccountable())
		hasChanges = true
	}
	if reqTask.GetObjective() != patchedTask.GetObjective() {
		patchedTask.SetObjective(reqTask.GetObjective())
		hasChanges = true
	}
	if reqTask.GetDeadline() != patchedTask.GetDeadline() {
		patchedTask.SetDeadline(reqTask.GetDeadline())
		hasChanges = true
	}
	return hasChanges
}

func (s *serviceMgmtOrder) PatchDelegatedTask(ctx context.Context, req *request.PatchDelegatedTaskRequest) (*response.PatchDelegatedTaskResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PatchDelegatedTask")
	defer middleware.SpanStop(ctx, "PatchDelegatedTask")

	if err := ValidatePatchDelegatedTaskRequest(req); err != nil {
		return nil, err
	}

	ordr, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := ordr.Clone()

	var patchedDelegatedTask *order.Task
	for _, delegatedTask := range patchedOrder.GetDelegatedTasks() {
		if delegatedTask.GetID() == req.GetTask().GetID() {
			patchedDelegatedTask = delegatedTask
			break
		}
	}
	if patchedDelegatedTask == nil {
		return nil, errwrap.NewError(http.StatusNotFound, "delegated task not found")
	}

	now := time.Now().UTC()

	hasChanges := patchTask(req.GetTask(), patchedDelegatedTask)
	if !hasChanges { // NOTE: no changes, return existing order
		return &response.PatchDelegatedTaskResponse{
			Order: ordr,
		}, nil
	}

	patchedOrder.GetMeta().VersionIncr()
	patchedOrder.GetMeta().SetUpdated(now)

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PatchDelegatedTaskResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteDelegatedTask(ctx context.Context, req *request.DeleteDelegatedTaskRequest) (*response.DeleteDelegatedTaskResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "DeleteDelegatedTask")
	defer middleware.SpanStop(ctx, "DeleteDelegatedTask")

	if err := ValidateDeleteDelegatedTaskRequest(req); err != nil {
		return nil, err
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := o.Clone()

	var patchedDelegatedTask *order.Task
	var idx int
	for i, delegatedTask := range patchedOrder.GetDelegatedTasks() {
		if delegatedTask.GetID() == req.GetDelegatedTaskID() {
			patchedDelegatedTask = delegatedTask
			idx = i
			break
		}
	}
	if patchedDelegatedTask == nil {
		return &response.DeleteDelegatedTaskResponse{}, nil

	}

	patchedOrder.SetDelegatedTasks(append(patchedOrder.GetDelegatedTasks()[:idx], patchedOrder.GetDelegatedTasks()[idx+1:]...))

	_, err = s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.DeleteDelegatedTaskResponse{}, nil
}

func (s *serviceMgmtOrder) PutSitRep(ctx context.Context, req *request.PutSitRepRequest) (*response.PutSitRepResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PutSitRep")
	defer middleware.SpanStop(ctx, "PutSitRep")

	if err := ValidatePutSitRepRequest(req); err != nil {
		return nil, err
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := order.Clone()

	sitrep := req.GetSitRep()
	now := time.Now().UTC()
	sitrep.ID = meta.ID(utils.RandAlphanum())

	patchedOrder.GetMeta().SetUpdated(now)
	patchedOrder.GetMeta().VersionIncr()

	patchedOrder.SetSitReps(append(patchedOrder.GetSitReps(), sitrep))

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PutSitRepResponse{
		Order: resp,
	}, nil
}

func patchSitRep(reqSitRep *order.SitRep, patchedSitRep *order.SitRep) bool {
	hasChanges := false

	if reqSitRep.GetState() != patchedSitRep.GetState() {
		patchedSitRep.SetState(reqSitRep.GetState())
		hasChanges = true
	}
	if reqSitRep.GetWorkCompleted() != patchedSitRep.GetWorkCompleted() {
		patchedSitRep.SetWorkCompleted(reqSitRep.GetWorkCompleted())
		hasChanges = true
	}
	if reqSitRep.GetSummary() != patchedSitRep.GetSummary() {
		patchedSitRep.SetSummary(reqSitRep.GetSummary())
		hasChanges = true
	}
	return hasChanges
}
func (s *serviceMgmtOrder) PatchSitRep(ctx context.Context, req *request.PatchSitRepRequest) (*response.PatchSitRepResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PatchSitRep")
	defer middleware.SpanStop(ctx, "PatchSitRep")

	if err := ValidatePatchSitRepRequest(req); err != nil {
		return nil, err
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := o.Clone()

	var patchedSitRep *order.SitRep
	for _, sitrep := range patchedOrder.GetSitReps() {
		if sitrep.GetID() == req.GetSitRep().GetID() {
			patchedSitRep = sitrep
			break
		}
	}
	if patchedSitRep == nil {
		return nil, errwrap.NewError(http.StatusNotFound, "sitrep not found")
	}

	now := time.Now().UTC()
	hasChanges := patchSitRep(req.GetSitRep(), patchedSitRep)
	if !hasChanges { // NOTE: no changes, return current order
		return &response.PatchSitRepResponse{
			Order: o,
		}, nil
	}

	patchedOrder.GetMeta().SetUpdated(now)
	patchedOrder.GetMeta().VersionIncr()

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PatchSitRepResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteSitRep(ctx context.Context, req *request.DeleteSitRepRequest) (*response.DeleteSitRepResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "DeleteSitRep")
	defer middleware.SpanStop(ctx, "DeleteSitRep")

	if err := ValidateDeleteSitRepRequest(req); err != nil {
		return nil, err
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := o.Clone()

	var patchedSitrep *order.Task
	var idx int
	for i, sitrep := range patchedOrder.GetDelegatedTasks() {
		if sitrep.GetID() == req.GetSitRepID() {
			patchedSitrep = sitrep
			idx = i
			break
		}
	}
	if patchedSitrep == nil {
		return &response.DeleteSitRepResponse{}, nil
	}

	patchedOrder.SetDelegatedTasks(append(patchedOrder.GetDelegatedTasks()[:idx], patchedOrder.GetDelegatedTasks()[idx+1:]...))

	_, err = s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.DeleteSitRepResponse{}, nil
}
