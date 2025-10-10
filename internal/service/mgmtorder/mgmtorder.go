package mgmtorder

import (
	"context"
	"slices"
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

	o := req.GetOrder().Clone()
	o.SetID(meta.NewID())

	for _, delegated := range o.GetDelegatedTasks() {
		delegated.SetID(meta.NewID())
	}
	for _, sitrep := range o.GetSitReps() {
		sitrep.SetID(meta.NewID())
	}

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

	order, err := s.Repository.ReadByID(ctx, req.GetOrder().GetID())
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	hasChanges := false
	patchedOrder := order.Clone()

	hasChanges = patchTask(req.GetOrder().GetTask(), patchedOrder.GetTask()) || hasChanges

	// TODO: optimize
	for _, reqDelegatedTask := range req.GetOrder().GetDelegatedTasks() {
		for _, patchedDelegatedTask := range patchedOrder.GetDelegatedTasks() {
			if reqDelegatedTask.GetID() == patchedDelegatedTask.GetID() {
				hasChanges = patchTask(reqDelegatedTask, patchedDelegatedTask) || hasChanges
				break
			}
		}
	}

	// TODO: optimize
	for _, reqSitRep := range req.GetOrder().GetSitReps() {
		for _, patchedSitRep := range patchedOrder.GetSitReps() {
			if reqSitRep.GetID() == patchedSitRep.GetID() {
				hasChanges = patchSitRep(reqSitRep, patchedSitRep) || hasChanges
				break
			}
		}
	}
	if !utils.IsZeroValue(req.GetOrder().GetParentOrderID()) && req.GetOrder().GetParentOrderID() != patchedOrder.GetParentOrderID() {
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
	}, nil
}

func (s *serviceMgmtOrder) DeleteOrder(ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "DeleteOrder")
	defer middleware.SpanStop(ctx, "DeleteOrder")

	if err := ValidateDeleteOrderRequest(req); err != nil {
		return nil, err
	}

	return &response.DeleteOrderResponse{}, s.Repository.DeleteOrder(ctx, req.GetID())
}

func (s *serviceMgmtOrder) PutDelegatedTask(ctx context.Context, req *request.PutDelegatedTasksRequest) (*response.PutDelegatedTasksResponse, errwrap.Error) {
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

	tasks := req.GetTasks()
	now := time.Now().UTC()
	for _, task := range tasks {
		task.SetID(meta.ID(utils.RandAlphanum()))
		patchedOrder.GetMeta().SetUpdated(now)
		patchedOrder.GetMeta().VersionIncr()
		patchedOrder.SetDelegatedTasks(append(patchedOrder.GetDelegatedTasks(), task))
	}

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PutDelegatedTasksResponse{
		Order: resp,
	}, nil
}

func patchTask(reqTask *order.Task, patchedTask *order.Task) bool {
	hasChanges := false

	if !utils.IsZeroValue(reqTask.GetState()) && reqTask.GetState() != patchedTask.GetState() {
		patchedTask.SetState(reqTask.GetState())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqTask.GetAccountable()) && reqTask.GetAccountable() != patchedTask.GetAccountable() {
		patchedTask.SetAccountable(reqTask.GetAccountable())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqTask.GetObjective()) && reqTask.GetObjective() != patchedTask.GetObjective() {
		patchedTask.SetObjective(reqTask.GetObjective())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqTask.GetDeadline()) && reqTask.GetDeadline() != patchedTask.GetDeadline() {
		patchedTask.SetDeadline(reqTask.GetDeadline())
		hasChanges = true
	}
	return hasChanges
}

func (s *serviceMgmtOrder) PatchDelegatedTask(ctx context.Context, req *request.PatchDelegatedTasksRequest) (*response.PatchDelegatedTasksResponse, errwrap.Error) {
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

	patchedOrderDelegatedTasks := make(map[meta.ID]*order.Task)
	for _, delegated := range patchedOrder.GetDelegatedTasks() {
		patchedOrderDelegatedTasks[delegated.GetID()] = delegated
	}

	now := time.Now().UTC()
	hasChanges := false
	for _, delegated := range req.GetTasks() {
		patchedDelegatedTask, ok := patchedOrderDelegatedTasks[delegated.GetID()]
		if !ok {
			continue
		}
		hasChanges = patchTask(delegated, patchedDelegatedTask) || hasChanges

	}
	if !hasChanges { // NOTE: no changes, return existing order
		return &response.PatchDelegatedTasksResponse{
			Order: ordr,
		}, nil
	}

	patchedOrder.GetMeta().VersionIncr()
	patchedOrder.GetMeta().SetUpdated(now)

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PatchDelegatedTasksResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteDelegatedTask(ctx context.Context, req *request.DeleteDelegatedTasksRequest) (*response.DeleteDelegatedTasksResponse, errwrap.Error) {
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

	patchedOrder.DelegatedTasks = slices.DeleteFunc(patchedOrder.DelegatedTasks, func(a *order.Task) bool {
		return slices.Contains(req.GetDelegatedTaskIDs(), a.GetID())
	})

	_, err = s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	for _, delegatedID := range req.GetDelegatedTaskIDs() {
		err = s.Repository.DeleteTask(ctx, delegatedID)
		if err != nil {
			// TODO: log warning
		}
	}
	return &response.DeleteDelegatedTasksResponse{}, nil
}

func (s *serviceMgmtOrder) PutSitRep(ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error) {
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

	sitreps := req.GetSitReps()
	now := time.Now().UTC()
	for _, sitrep := range sitreps {
		sitrep.SetID(meta.ID(utils.RandAlphanum()))
		patchedOrder.GetMeta().SetUpdated(now)
		patchedOrder.GetMeta().VersionIncr()
		patchedOrder.SetSitReps(append(patchedOrder.GetSitReps(), sitrep))
	}

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PutSitRepsResponse{
		Order: resp,
	}, nil
}

func patchSitRep(reqSitRep *order.SitRep, patchedSitRep *order.SitRep) bool {
	hasChanges := false

	if !utils.IsZeroValue(reqSitRep.GetDateTime()) && reqSitRep.GetDateTime() != patchedSitRep.GetDateTime() {
		patchedSitRep.SetDateTime(reqSitRep.GetDateTime())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetBy()) && reqSitRep.GetBy() != patchedSitRep.GetBy() {
		patchedSitRep.SetBy(reqSitRep.GetBy())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetPing()) && !slices.Equal(reqSitRep.GetPing(), patchedSitRep.GetPing()) {
		patchedSitRep.SetPing(reqSitRep.GetPing())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetSituation()) && reqSitRep.GetSituation() != patchedSitRep.GetSituation() {
		patchedSitRep.SetSituation(reqSitRep.GetSituation())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetActions()) && reqSitRep.GetActions() != patchedSitRep.GetActions() {
		patchedSitRep.SetActions(reqSitRep.GetActions())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetTBD()) && reqSitRep.GetTBD() != patchedSitRep.GetTBD() {
		patchedSitRep.SetTBD(reqSitRep.GetTBD())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetIssues()) && reqSitRep.GetIssues() != patchedSitRep.GetIssues() {
		patchedSitRep.SetIssues(reqSitRep.GetIssues())
		hasChanges = true
	}
	return hasChanges
}
func (s *serviceMgmtOrder) PatchSitRep(ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error) {
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

	patchedOrderSitReps := make(map[meta.ID]*order.SitRep)
	for _, sitrep := range patchedOrder.GetSitReps() {
		patchedOrderSitReps[sitrep.GetID()] = sitrep
	}

	now := time.Now().UTC()
	hasChanges := false
	for _, sitrep := range req.GetSitReps() {
		patchedSitRep, ok := patchedOrderSitReps[sitrep.GetID()]
		if !ok {
			continue
		}
		hasChanges = patchSitRep(sitrep, patchedSitRep) || hasChanges

	}

	if !hasChanges { // NOTE: no changes, return current order
		return &response.PatchSitRepsResponse{
			Order: o,
		}, nil
	}

	patchedOrder.GetMeta().SetUpdated(now)
	patchedOrder.GetMeta().VersionIncr()

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	return &response.PatchSitRepsResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteSitRep(ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error) {
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

	patchedOrder.SitReps = slices.DeleteFunc(patchedOrder.SitReps, func(a *order.SitRep) bool {
		return slices.Contains(req.GetSitRepIDs(), a.GetID())
	})

	_, err = s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, err
	}
	for _, delegatedID := range req.GetSitRepIDs() {
		err = s.Repository.DeleteTask(ctx, delegatedID)
		if err != nil {
			// TODO: log warning
		}
	}
	return &response.DeleteSitRepsResponse{}, nil
}
