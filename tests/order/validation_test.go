package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/service/common/validation"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestValidation_Task() {
	tt := s.T()

	tt.Run("task.id", func(t *testing.T) {
		to := setup.TaskObjWithID()
		to.SetID("")
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("task.state.low", func(t *testing.T) {
		to := setup.TaskObjWithID()
		to.SetState(order.NotStarted - 1)
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.state", err.GetStatusMessage())
	})
	tt.Run("task.state.high", func(t *testing.T) {
		to := setup.TaskObjWithID()
		to.SetState(order.Completed + 1)
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.state", err.GetStatusMessage())
	})
	tt.Run("task.accountable", func(t *testing.T) {
		to := setup.TaskObjWithID()
		to.SetAccountable("")
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid email length", err.GetStatusMessage())
	})
	tt.Run("task.objective", func(t *testing.T) {
		to := setup.TaskObjWithID()
		to.SetObjective("")
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.objective", err.GetStatusMessage())
	})
	tt.Run("task.deadline", func(t *testing.T) {
		to := setup.TaskObjWithID()
		to.SetDeadline(time.Time{})
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.deadline", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_SitReps() {
	tt := s.T()

	tt.Run("sitrep.id", func(t *testing.T) {
		sp := setup.SitrepObjWithID()
		sp.SetID("")
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("sitrep.datetime", func(t *testing.T) {
		sp := setup.SitrepObjWithID()
		sp.SetDateTime(time.Time{})
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.datetime", err.GetStatusMessage())
	})
	tt.Run("sitrep.by", func(t *testing.T) {
		sp := setup.SitrepObjWithID()
		sp.SetBy("")
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.by: invalid email length", err.GetStatusMessage())
	})
	tt.Run("sitrep.ping", func(t *testing.T) {
		sp := setup.SitrepObjWithID()
		sp.SetPing(append([]user.Email{""}, setup.SitrepObjWithID().GetPing()...))
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.ping.0: invalid email length", err.GetStatusMessage())
	})
	tt.Run("sitrep.no_content", func(t *testing.T) {
		sp := setup.SitrepObjWithID()
		sp.SetSituation("")
		sp.SetActions("")
		sp.SetTBD("")
		sp.SetIssues("")
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitrep", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_Order() {
	tt := s.T()

	tt.Run("order.task.nil", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		o.SetTask(nil)
		err := mgmtorder.ValidateOrder(o, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.task", err.GetStatusMessage())
	})
	tt.Run("order.task.empty", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		o.SetTask(&order.Task{})
		err := mgmtorder.ValidateOrder(o, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.task: invalid task.id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("order.parent_order_id", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		o.SetParentOrderID("")
		err := mgmtorder.ValidateOrder(o, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.parent_order_id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("order.meta", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		o.SetMeta(&meta.Meta{
			Version: 2,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		})
		err := mgmtorder.ValidateOrder(o, validation.IgnoreNothing)
		require.NoError(t, err) // NOTE: request.meta is ignored
	})
}

func (s *OrderSuite) TestValidation_PostOrderRequest() {
	tt := s.T()

	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PostOrder(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})

	tt.Run("order.task.id_provided", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		o.GetTask().SetID(meta.NewID())
		resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "order.task.id disallowed", err.GetStatusMessage())
	})
	tt.Run("order.delegated_task.id_provided", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		setup.ZeroOrderIDs(o)
		o.GetDelegatedTasks()[0].SetID(meta.NewID())
		resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "order.delegated_tasks.0.id disallowed", err.GetStatusMessage())
	})
	tt.Run("order.sitrep.id_provided", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		setup.ZeroOrderIDs(o)
		o.GetSitReps()[0].SetID(meta.NewID())
		resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "order.sitreps.0.id disallowed", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_GetOrderByIDRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.GetOrderByID(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_GetOrdersRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.GetOrders(t, context.Background(), nil)
		require.NoError(t, err)
		require.Empty(t, resp)
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.GetOrders(t, context.Background(), &request.GetOrdersRequest{})
		require.NoError(t, err)
		require.Empty(t, resp)
	})
}

func (s *OrderSuite) TestValidation_GetOrderSubOrdersRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.GetOrderSubOrders(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.GetOrderSubOrders(t, context.Background(), &request.GetOrderSubOrdersRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PatchOrderRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PatchOrder(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PatchOrder(t, context.Background(), &request.PatchOrderRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})
	tt.Run("order.task.id.empty", func(t *testing.T) {
		o := setup.OrderObj()
		resp, err := s.API.PatchOrder(t, context.Background(), &request.PatchOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.task.id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("order.delegated_task.id.empty", func(t *testing.T) {

		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())

		o.GetDelegatedTasks()[0].SetID(meta.EmptyID())
		o.GetDelegatedTasks()[0].SetObjective("patched objective")
		o.SetSitReps(nil)
		resp, err := s.API.PatchOrder(t, context.Background(), &request.PatchOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.delegated_tasks.0: invalid task.id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("order.sitrep.id.empty", func(t *testing.T) {

		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())

		o.GetSitReps()[0].SetID(meta.EmptyID())
		o.GetSitReps()[0].SetActions("patched actions")
		resp, err := s.API.PatchOrder(t, context.Background(), &request.PatchOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.sitreps.0: invalid sitrep.id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_DeleteOrderRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.DeleteOrder(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.DeleteOrder(t, context.Background(), &request.DeleteOrderRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PutDelegatedTaskRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PutDelegatedTasks(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PutDelegatedTasks(t, context.Background(), &request.PutDelegatedTasksRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		setup.ZeroOrderIDs(o)
		resp, err := s.API.PutDelegatedTasks(t, context.Background(), &request.PutDelegatedTasksRequest{
			Tasks: []*order.Task{o.GetTask()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.nil", func(t *testing.T) {
		resp, err := s.API.PutDelegatedTasks(t, context.Background(), &request.PutDelegatedTasksRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty tasks", err.GetStatusMessage())
	})
	tt.Run("empty.task", func(t *testing.T) {
		resp, err := s.API.PutDelegatedTasks(t, context.Background(), &request.PutDelegatedTasksRequest{
			OrderID: meta.NewID(),
			Tasks:   []*order.Task{},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty tasks", err.GetStatusMessage())
	})
	tt.Run("delegated_task.id.provided", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		oid := o.GetID()
		resp, err := s.API.PutDelegatedTasks(t, context.Background(), &request.PutDelegatedTasksRequest{
			OrderID: oid,
			Tasks:   []*order.Task{o.GetTask()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "tasks.0.id disallowed", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PatchDelegatedTaskRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PatchDelegatedTasks(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PatchDelegatedTasks(t, context.Background(), &request.PatchDelegatedTasksRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		setup.ZeroOrderIDs(o)
		resp, err := s.API.PatchDelegatedTasks(t, context.Background(), &request.PatchDelegatedTasksRequest{
			Tasks: []*order.Task{o.GetTask()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("nil.task", func(t *testing.T) {
		resp, err := s.API.PatchDelegatedTasks(t, context.Background(), &request.PatchDelegatedTasksRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty tasks", err.GetStatusMessage())
	})
	tt.Run("empty.task", func(t *testing.T) {
		resp, err := s.API.PatchDelegatedTasks(t, context.Background(), &request.PatchDelegatedTasksRequest{
			OrderID: meta.NewID(),
			Tasks:   []*order.Task{},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty tasks", err.GetStatusMessage())
	})
	tt.Run("delegated_task.id.empty", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		oid := o.GetID()
		setup.ZeroOrderIDs(o)

		resp, err := s.API.PatchDelegatedTasks(t, context.Background(), &request.PatchDelegatedTasksRequest{
			OrderID: oid,
			Tasks:   []*order.Task{o.GetTask()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid tasks.0.id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_DeleteDelegatedTaskRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTasks(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTasks(t, context.Background(), &request.DeleteDelegatedTasksRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTasks(t, context.Background(), &request.DeleteDelegatedTasksRequest{
			DelegatedTaskIDs: []meta.ID{meta.NewID()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("nil.delegated_task_id", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTasks(t, context.Background(), &request.DeleteDelegatedTasksRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty delegated_task_ids", err.GetStatusMessage())
	})
	tt.Run("invalid.delegated_task_id", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTasks(t, context.Background(), &request.DeleteDelegatedTasksRequest{
			OrderID:          meta.NewID(),
			DelegatedTaskIDs: []meta.ID{meta.NewID(), meta.EmptyID()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid delegated_task_ids.1: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PutSitRepRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PutSitReps(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PutSitReps(t, context.Background(), &request.PutSitRepsRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		setup.ZeroOrderIDs(o)
		resp, err := s.API.PutSitReps(t, context.Background(), &request.PutSitRepsRequest{
			SitReps: []*order.SitRep{o.GetSitReps()[0]},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("nil.sitrep", func(t *testing.T) {
		resp, err := s.API.PutSitReps(t, context.Background(), &request.PutSitRepsRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitreps", err.GetStatusMessage())
	})
	tt.Run("empty.sitrep", func(t *testing.T) {
		resp, err := s.API.PutSitReps(t, context.Background(), &request.PutSitRepsRequest{
			OrderID: meta.NewID(),
			SitReps: []*order.SitRep{},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitreps", err.GetStatusMessage())
	})
	tt.Run("sitrep.id.provided", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		oid := o.GetID()
		resp, err := s.API.PutSitReps(t, context.Background(), &request.PutSitRepsRequest{
			OrderID: oid,
			SitReps: []*order.SitRep{o.GetSitReps()[0]},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "sitreps.0.id disallowed", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PatchSitRepRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PatchSitReps(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PatchSitReps(t, context.Background(), &request.PatchSitRepsRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		setup.ZeroOrderIDs(o)
		resp, err := s.API.PatchSitReps(t, context.Background(), &request.PatchSitRepsRequest{
			SitReps: []*order.SitRep{o.GetSitReps()[0]},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("nil.sitrep", func(t *testing.T) {
		resp, err := s.API.PatchSitReps(t, context.Background(), &request.PatchSitRepsRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitreps", err.GetStatusMessage())
	})
	tt.Run("empty.sitrep", func(t *testing.T) {
		resp, err := s.API.PatchSitReps(t, context.Background(), &request.PatchSitRepsRequest{
			OrderID: meta.NewID(),
			SitReps: []*order.SitRep{},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitreps", err.GetStatusMessage())
	})
	tt.Run("sitrep.id.empty", func(t *testing.T) {
		o := setup.OrderObjWithIDs()
		oid := o.GetID()
		setup.ZeroOrderIDs(o)

		resp, err := s.API.PatchSitReps(t, context.Background(), &request.PatchSitRepsRequest{
			OrderID: oid,
			SitReps: []*order.SitRep{o.GetSitReps()[0]},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitreps.0.id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_DeleteSitRepRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.DeleteSitReps(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		resp, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{
			SitRepIDs: []meta.ID{meta.NewID()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("nil.sitrep_ids", func(t *testing.T) {
		resp, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitrep_ids", err.GetStatusMessage())
	})
	tt.Run("empty.sitrep_ids", func(t *testing.T) {
		resp, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{
			OrderID:   meta.NewID(),
			SitRepIDs: []meta.ID{},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitrep_ids", err.GetStatusMessage())
	})
	tt.Run("invalid.sitrep_ids", func(t *testing.T) {
		resp, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{
			OrderID:   meta.NewID(),
			SitRepIDs: []meta.ID{meta.NewID(), meta.EmptyID()},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep_ids.1: invalid id length", err.GetStatusMessage())
	})
}
