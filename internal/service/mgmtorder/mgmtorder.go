package mgmtorder

import (
	"context"
	"net/http"
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

func (s *serviceMgmtOrder) GetRootOrder(context.Context) *order.Order {
	// TODO: utilize ctx
	return s.RootOrder
}

func (s *serviceMgmtOrder) PostOrder(ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "PostOrder")
	defer middleware.SpanStop(ctx, "PostOrder")

	if err := ValidatePostOrderRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	o := req.GetOrder().Clone()
	o.SetID(meta.NewID())

	for _, delegated := range o.GetDelegatedOrders() {
		delegated.SetID(meta.NewID())
		delegated.SetParentOrderID(o.GetID())
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
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.PostOrderResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) GetOrderByID(ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "GetOrderByID")
	defer middleware.SpanStop(ctx, "GetOrderByID")

	if err := ValidateGetOrderByIDRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	resp, err := s.Repository.ReadByID(ctx, req.GetID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.GetOrderByIDResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) GetOrders(ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "GetOrders")
	defer middleware.SpanStop(ctx, "GetOrders")

	if err := ValidateGetOrdersRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	resp, err := s.Repository.ReadBy(ctx, req)
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.GetOrdersResponse{
		Orders: resp,
	}, nil
}

func (s *serviceMgmtOrder) PatchOrder(ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "PatchOrder")
	defer middleware.SpanStop(ctx, "PatchOrder")

	if err := ValidatePatchOrderRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrder().GetID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	now := time.Now().UTC()
	hasChanges := false
	patchedOrder := order.Clone()

	hasChanges = patchOrder(req.GetOrder(), patchedOrder) || hasChanges

	// TODO: optimize
	for _, reqDelegatedOrder := range req.GetOrder().GetDelegatedOrders() {
		for _, patchedDelegatedOrder := range patchedOrder.GetDelegatedOrders() {
			if reqDelegatedOrder.GetID() == patchedDelegatedOrder.GetID() {
				hasChanges = patchOrder(reqDelegatedOrder, patchedDelegatedOrder) || hasChanges
				break
			}
		}
	}

	// TODO: optimize
	for _, reqSitRep := range req.GetOrder().GetSitReps() {
		for _, patchedSitRep := range patchedOrder.GetSitReps() {
			if reqSitRep.GetID() == patchedSitRep.GetID() {
				hasChanges = patchSitReps(reqSitRep, patchedSitRep) || hasChanges
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
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.PatchOrderResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteOrder(ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "DeleteOrder")
	defer middleware.SpanStop(ctx, "DeleteOrder")

	if err := ValidateDeleteOrderRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	patchedOrder, err := s.Repository.ReadByID(ctx, req.GetID())
	if err != nil && err.GetStatusCode() != http.StatusNotFound {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	if err != nil && err.GetStatusCode() == http.StatusNotFound {
		return &response.DeleteOrderResponse{}, nil
	}

	err = s.Repository.DeleteOrder(ctx, req.GetID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	patchedOrder.GetMeta().VersionIncr()
	patchedOrder.GetMeta().SetUpdated(time.Now().UTC())

	return &response.DeleteOrderResponse{}, nil
}

func (s *serviceMgmtOrder) PutDelegatedOrders(ctx context.Context, req *request.PutDelegatedOrdersRequest) (*response.PutDelegatedOrdersResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "PutDelegatedOrders")
	defer middleware.SpanStop(ctx, "PutDelegatedOrders")

	if err := ValidatePutDelegatedOrderRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	patchedOrder := order.Clone()

	tasks := req.GetOrders()
	now := time.Now().UTC()
	for _, task := range tasks {
		task.SetID(meta.ID(utils.RandAlphanum()))
		patchedOrder.GetMeta().SetCreated(now)
		patchedOrder.GetMeta().SetUpdated(now)
		patchedOrder.GetMeta().VersionIncr()
		patchedOrder.SetDelegatedOrders(append(patchedOrder.GetDelegatedOrders(), task))
	}

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.PutDelegatedOrdersResponse{
		Order: resp,
	}, nil
}

func patchOrder(reqOrder *order.Order, patchedOrder *order.Order) bool {
	hasChanges := false

	if !utils.IsZeroValue(reqOrder.GetState()) && reqOrder.GetState() != patchedOrder.GetState() {
		patchedOrder.SetState(reqOrder.GetState())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqOrder.GetAccountableID()) && reqOrder.GetAccountableID() != patchedOrder.GetAccountableID() {
		patchedOrder.SetAccountableID(reqOrder.GetAccountableID())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqOrder.GetObjective()) && reqOrder.GetObjective() != patchedOrder.GetObjective() {
		patchedOrder.SetObjective(reqOrder.GetObjective())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqOrder.GetDeadline()) && reqOrder.GetDeadline() != patchedOrder.GetDeadline() {
		patchedOrder.SetDeadline(reqOrder.GetDeadline())
		hasChanges = true
	}
	return hasChanges
}

func (s *serviceMgmtOrder) PatchDelegatedOrders(ctx context.Context, req *request.PatchDelegatedOrdersRequest) (*response.PatchDelegatedOrdersResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "PatchDelegatedOrders")
	defer middleware.SpanStop(ctx, "PatchDelegatedOrders")

	if err := ValidatePatchDelegatedOrderRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	ordr, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	patchedOrder := ordr.Clone()

	patchedOrderDelegatedOrders := make(map[meta.ID]*order.Order)
	for _, delegated := range patchedOrder.GetDelegatedOrders() {
		patchedOrderDelegatedOrders[delegated.GetID()] = delegated
	}

	now := time.Now().UTC()
	hasChanges := false
	for _, delegated := range req.GetOrders() {
		patchedDelegatedOrder, ok := patchedOrderDelegatedOrders[delegated.GetID()]
		if !ok {
			continue
		}
		hasChanges = patchOrder(delegated, patchedDelegatedOrder) || hasChanges

	}
	if !hasChanges { // NOTE: no changes, return existing order
		return &response.PatchDelegatedOrdersResponse{
			Order: ordr,
		}, nil
	}

	patchedOrder.GetMeta().VersionIncr()
	patchedOrder.GetMeta().SetUpdated(now)

	resp, err := s.Repository.Write(ctx, patchedOrder)
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.PatchDelegatedOrdersResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteDelegatedOrders(ctx context.Context, req *request.DeleteDelegatedOrdersRequest) (*response.DeleteDelegatedOrdersResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "DeleteDelegatedOrders")
	defer middleware.SpanStop(ctx, "DeleteDelegatedOrders")

	if err := ValidateDeleteDelegatedOrderRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	patchedOrder := o.Clone()

	patchedOrder.DelegatedOrders = slices.DeleteFunc(patchedOrder.DelegatedOrders, func(a *order.Order) bool {
		return slices.Contains(req.GetDelegatedOrderIDs(), a.GetID())
	})

	didDelete, err := s.Repository.DeleteOrders(ctx, req.GetDelegatedOrderIDs()) // NOTE: delete tasks in db.tasks
	if err != nil {
		// TODO: log warning
	}

	if didDelete {
		patchedOrder.GetMeta().VersionIncr()
		patchedOrder.GetMeta().SetUpdated(time.Now().UTC())
		patchedOrder, err = s.Repository.Write(ctx, patchedOrder) // NOTE: update order obj in db.order
	}

	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.DeleteDelegatedOrdersResponse{
		Order: patchedOrder,
	}, nil
}

func (s *serviceMgmtOrder) PutSitReps(ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "PutSitReps")
	defer middleware.SpanStop(ctx, "PutSitReps")

	if err := ValidatePutSitRepRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	order, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
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
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.PutSitRepsResponse{
		Order: resp,
	}, nil
}

func patchSitReps(reqSitRep *order.SitRep, patchedSitRep *order.SitRep) bool {
	hasChanges := false

	if !utils.IsZeroValue(reqSitRep.GetDateTime()) && reqSitRep.GetDateTime() != patchedSitRep.GetDateTime() {
		patchedSitRep.SetDateTime(reqSitRep.GetDateTime())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetBy()) && reqSitRep.GetBy() != patchedSitRep.GetBy() {
		patchedSitRep.SetBy(reqSitRep.GetBy())
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
	if !utils.IsZeroValue(reqSitRep.GetTODO()) && reqSitRep.GetTODO() != patchedSitRep.GetTODO() {
		patchedSitRep.SetTODO(reqSitRep.GetTODO())
		hasChanges = true
	}
	if !utils.IsZeroValue(reqSitRep.GetIssues()) && reqSitRep.GetIssues() != patchedSitRep.GetIssues() {
		patchedSitRep.SetIssues(reqSitRep.GetIssues())
		hasChanges = true
	}
	return hasChanges
}
func (s *serviceMgmtOrder) PatchSitReps(ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "PatchSitReps")
	defer middleware.SpanStop(ctx, "PatchSitReps")

	if err := ValidatePatchSitRepRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
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
		hasChanges = patchSitReps(sitrep, patchedSitRep) || hasChanges

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
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.PatchSitRepsResponse{
		Order: resp,
	}, nil
}

func (s *serviceMgmtOrder) DeleteSitReps(ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error) {
	ctx = middleware.AddTraceToCtx(ctx)
	middleware.SpanStart(ctx, "DeleteSitReps")
	defer middleware.SpanStop(ctx, "DeleteSitReps")

	if err := ValidateDeleteSitRepRequest(req); err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	o, err := s.Repository.ReadByID(ctx, req.GetOrderID())
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}

	patchedOrder := o.Clone()

	patchedOrder.SitReps = slices.DeleteFunc(patchedOrder.SitReps, func(a *order.SitRep) bool {
		return slices.Contains(req.GetSitRepIDs(), a.GetID())
	})

	didDelete, err := s.Repository.DeleteSitReps(ctx, req.GetSitRepIDs()) // NOTE: delete sitrep in db.sitrep
	if err != nil {
		// TODO: log warning
	}

	if didDelete {
		patchedOrder.GetMeta().VersionIncr()
		patchedOrder.GetMeta().SetUpdated(time.Now().UTC())
		patchedOrder, err = s.Repository.Write(ctx, patchedOrder) // NOTE: update order obj in db.order
	}

	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return &response.DeleteSitRepsResponse{
		Order: patchedOrder,
	}, nil
}
