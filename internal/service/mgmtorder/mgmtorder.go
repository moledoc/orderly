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

	if err := validatePostOrderRequest(req); err != nil {
		return nil, err
	}

	o := req.GetOrder()
	o.Task.ID = utils.Ptr(meta.ID(utils.RandAlphanum()))

	now := time.Now().UTC()
	o.Task.Meta = &meta.Meta{
		Version: 1,
		Created: now,
		Updated: now,
		Deleted: false,
	}
	for _, subtask := range o.GetSubTasks() {
		subtask.Meta = &meta.Meta{
			Version: 1,
			Created: now,
			Updated: now,
			Deleted: false,
		}
	}
	for _, sitrep := range o.GetSitReps() {
		sitrep.Meta = &meta.Meta{
			Version: 1,
			Created: now,
			Updated: now,
			Deleted: false,
		}
	}

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

	if err := validateGetOrderByIDRequest(req); err != nil {
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

	if err := validateGetOrdersRequest(req); err != nil {
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

func (s *serviceMgmtOrder) GetOrderVersions(ctx context.Context, req *request.GetOrderVersionsRequest) (*response.GetOrderVersionsResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetOrderVersions")
	defer middleware.SpanStop(ctx, "GetOrderVersions")

	if err := validateGetOrderVersionsRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.Repository.ReadVersions(ctx, req.GetID())
	if err != nil {
		return nil, err
	}
	return &response.GetOrderVersionsResponse{
		OrderVersions: resp,
	}, err
}

func (s *serviceMgmtOrder) GetOrderSubOrders(ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetOrderSubOrders")
	defer middleware.SpanStop(ctx, "GetOrderSubOrders")

	if err := validateGetOrderSubOrdersRequest(req); err != nil {
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

	if err := validatePatchOrderRequest(req); err != nil {
		return nil, err
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrder().GetTask().GetID())
	if err != nil {
		return nil, err
	}

	patchedOrder := order.Clone()
	// TODO: apply req.user diff to patchedUser

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

	if err := validateDeleteOrderRequest(req); err != nil {
		return nil, err
	}

	var err errwrap.Error

	if req.GetHard() {
		return &response.DeleteOrderResponse{}, s.Repository.Delete(ctx, req.GetID())
	}

	order, err := s.Repository.ReadByID(ctx, req.GetID())
	if err != nil {
		return nil, err
	}

	softDeletedOrder := order.Clone()
	softDeletedOrder.GetTask().GetMeta().SetDeleted(true)
	_, err = s.Repository.Write(ctx, softDeletedOrder)

	return &response.DeleteOrderResponse{}, err

}

func (s *serviceMgmtOrder) PutSubTask(ctx context.Context, req *request.PutSubTaskRequest) (*response.PutSubTaskResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PutSubTask")
	defer middleware.SpanStop(ctx, "PutSubTask")

	if err := validatePutSubTaskRequest(req); err != nil {
		return nil, err
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := order.Clone()

	task := req.GetTask()
	now := time.Now().UTC()
	task.ID = utils.Ptr(meta.ID(utils.RandAlphanum()))
	task.Meta = &meta.Meta{
		Version: 1,
		Created: now,
		Updated: now,
		Deleted: false,
	}

	patchedOrder.SetSubTasks(append(patchedOrder.GetSubTasks(), task))

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PutSubTaskResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) PatchSubTask(ctx context.Context, req *request.PatchSubTaskRequest) (*response.PatchSubTaskResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PatchSubTask")
	defer middleware.SpanStop(ctx, "PatchSubTask")

	if err := validatePatchSubTaskRequest(req); err != nil {
		return nil, err
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := o.Clone()

	var patchedSubtask *order.Task
	for _, subtask := range patchedOrder.GetSubTasks() {
		if subtask.GetID() == req.GetTask().GetID() {
			patchedSubtask = subtask
			break
		}
	}
	if patchedSubtask == nil {
		return nil, errwrap.NewError(http.StatusNotFound, "subtask not found")
	}

	now := time.Now().UTC()
	tsk := req.GetTask()

	patchedSubtask.SetState(tsk.GetState())
	patchedSubtask.SetAccountable(tsk.GetAccountable())
	patchedSubtask.SetObjective(tsk.GetObjective())

	patchedSubtask.GetMeta().VersionIncr()
	patchedSubtask.GetMeta().SetUpdated(now)
	patchedSubtask.GetMeta().SetDeleted(tsk.GetMeta().GetDeleted())

	patchedOrder.SetSubTasks(append(patchedOrder.GetSubTasks(), tsk))

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PatchSubTaskResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteSubTask(ctx context.Context, req *request.DeleteSubTaskRequest) (*response.DeleteSubTaskResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "DeleteSubTask")
	defer middleware.SpanStop(ctx, "DeleteSubTask")

	if err := validateDeleteSubTaskRequest(req); err != nil {
		return nil, err
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := o.Clone()

	var patchedSubtask *order.Task
	var idx int
	for i, subtask := range patchedOrder.GetSubTasks() {
		if subtask.GetID() == req.GetSubTaskID() {
			patchedSubtask = subtask
			idx = i
			break
		}
	}
	if patchedSubtask == nil {
		return nil, errwrap.NewError(http.StatusNotFound, "subtask not found")
	}

	if req.GetHard() {
		patchedOrder.SetSubTasks(append(patchedOrder.GetSubTasks()[:idx], patchedOrder.GetSubTasks()[idx+1:]...))
	} else {
		now := time.Now().UTC()

		patchedSubtask.GetMeta().VersionIncr()
		patchedSubtask.GetMeta().SetUpdated(now)
		patchedSubtask.GetMeta().SetDeleted(true)
	}

	_, err = s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.DeleteSubTaskResponse{}, nil
}

func (s *serviceMgmtOrder) PutSitRep(ctx context.Context, req *request.PutSitRepRequest) (*response.PutSitRepResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PutSitRep")
	defer middleware.SpanStop(ctx, "PutSitRep")

	if err := validatePutSitRepRequest(req); err != nil {
		return nil, err
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := order.Clone()

	sitrep := req.GetSitRep()
	now := time.Now().UTC()
	sitrep.ID = utils.Ptr(meta.ID(utils.RandAlphanum()))
	sitrep.Meta = &meta.Meta{
		Version: 1,
		Created: now,
		Updated: now,
		Deleted: false,
	}

	patchedOrder.SetSitReps(append(patchedOrder.GetSitReps(), sitrep))

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PutSitRepResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) PatchSitRep(ctx context.Context, req *request.PatchSitRepRequest) (*response.PatchSitRepResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PatchSitRep")
	defer middleware.SpanStop(ctx, "PatchSitRep")

	if err := validatePatchSitRepRequest(req); err != nil {
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
	strp := req.GetSitRep()

	patchedSitRep.SetState(strp.GetState())
	patchedSitRep.SetWorkCompleted(strp.GetWorkCompleted())
	patchedSitRep.SetSummary(strp.GetSummary())

	patchedSitRep.GetMeta().VersionIncr()
	patchedSitRep.GetMeta().SetUpdated(now)
	patchedSitRep.GetMeta().SetDeleted(strp.GetMeta().GetDeleted())

	patchedOrder.SetSitReps(append(patchedOrder.GetSitReps(), strp))

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

	if err := validateDeleteSitRepRequest(req); err != nil {
		return nil, err
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, err
	}

	patchedOrder := o.Clone()

	var patchedSitrep *order.Task
	var idx int
	for i, sitrep := range patchedOrder.GetSubTasks() {
		if sitrep.GetID() == req.GetSitRepID() {
			patchedSitrep = sitrep
			idx = i
			break
		}
	}
	if patchedSitrep == nil {
		return nil, errwrap.NewError(http.StatusNotFound, "sitrep not found")
	}

	if req.GetHard() {
		patchedOrder.SetSubTasks(append(patchedOrder.GetSubTasks()[:idx], patchedOrder.GetSubTasks()[idx+1:]...))
	} else {
		now := time.Now().UTC()

		patchedSitrep.GetMeta().VersionIncr()
		patchedSitrep.GetMeta().SetUpdated(now)
		patchedSitrep.GetMeta().SetDeleted(true)
	}

	_, err = s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.DeleteSitRepResponse{}, nil
}
