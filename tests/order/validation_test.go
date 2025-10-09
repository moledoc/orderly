package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

func taskObj(extra ...string) *order.Task {
	ee := strings.Join(extra, ".")
	return &order.Task{
		State:       order.NotStarted,
		Accountable: user.Email(fmt.Sprintf("example%v@example.com", ee)),
		Objective:   "objective description",
		Deadline:    time.Now().UTC(),
	}
}
func taskObjWithID(extra ...string) *order.Task {
	tt := taskObj(extra...)
	tt.SetID(meta.NewID())
	return tt
}

func sitrepObj(extra ...string) *order.SitRep {
	ee := strings.Join(extra, ".")
	return &order.SitRep{
		State:         order.NotStarted,
		WorkCompleted: 50,
		Summary:       "summary",

		DateTime: time.Now().UTC(),
		By:       user.Email(fmt.Sprintf("by%v@example.com", ee)),
		Ping: []user.Email{
			user.Email(fmt.Sprintf("ping1%v@example.com", ee)),
			user.Email(fmt.Sprintf("ping2%v@example.com", ee)),
		},
		Situation: "situation description",
		Actions:   "list of actions taken",
		TBD:       "list of things to do still",
		Issues:    "list of encountered issues",
	}
}

func sitrepObjWithID(extra ...string) *order.SitRep {
	sr := sitrepObj(extra...)
	sr.SetID(meta.NewID())
	return sr
}

func orderObj(extra ...string) *order.Order {
	return &order.Order{
		Task: taskObj(extra...),
		DelegatedTasks: []*order.Task{
			taskObj(),
			taskObj(),
			taskObj(),
		},
		ParentOrderID: meta.NewID(),
		SitReps: []*order.SitRep{
			sitrepObj(),
			sitrepObj(),
			sitrepObj(),
		},
		Meta: &meta.Meta{
			Version: 1,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		},
	}
}

func orderObjWithIDs(extra ...string) *order.Order {
	return &order.Order{
		Task: taskObjWithID(extra...),
		DelegatedTasks: []*order.Task{
			taskObjWithID(),
			taskObjWithID(),
			taskObjWithID(),
		},
		ParentOrderID: meta.NewID(),
		SitReps: []*order.SitRep{
			sitrepObjWithID(),
			sitrepObjWithID(),
			sitrepObjWithID(),
		},
		Meta: &meta.Meta{
			Version: 1,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		},
	}
}

func zeroOrderIDs(o *order.Order) {
	o.GetTask().SetID("")
	for _, delegated := range o.GetDelegatedTasks() {
		delegated.SetID("")
	}
	for _, sitrep := range o.GetSitReps() {
		sitrep.SetID("")
	}
}

func (s *OrderSuite) TestValidation_Task() {
	tt := s.T()

	tt.Run("task.id", func(t *testing.T) {
		to := taskObjWithID()
		to.SetID("")
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("task.state.low", func(t *testing.T) {
		to := taskObjWithID()
		to.SetState(order.NotStarted - 1)
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.state", err.GetStatusMessage())
	})
	tt.Run("task.state.high", func(t *testing.T) {
		to := taskObjWithID()
		to.SetState(order.Completed + 1)
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.state", err.GetStatusMessage())
	})
	tt.Run("task.accountable", func(t *testing.T) {
		to := taskObjWithID()
		to.SetAccountable("")
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid email length", err.GetStatusMessage())
	})
	tt.Run("task.objective", func(t *testing.T) {
		to := taskObjWithID()
		to.SetObjective("")
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.objective", err.GetStatusMessage())
	})
	tt.Run("task.deadline", func(t *testing.T) {
		to := taskObjWithID()
		to.SetDeadline(time.Time{})
		err := mgmtorder.ValidateTask(to, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.deadline", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_SitRep() {
	tt := s.T()

	tt.Run("sitrep.id", func(t *testing.T) {
		sp := sitrepObjWithID()
		sp.SetID("")
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("sitrep.datetime", func(t *testing.T) {
		sp := sitrepObjWithID()
		sp.SetDateTime(time.Time{})
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.datetime", err.GetStatusMessage())
	})
	tt.Run("sitrep.by", func(t *testing.T) {
		sp := sitrepObjWithID()
		sp.SetBy("")
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.by: invalid email length", err.GetStatusMessage())
	})
	tt.Run("sitrep.ping", func(t *testing.T) {
		sp := sitrepObjWithID()
		sp.SetPing(append([]user.Email{""}, sitrepObjWithID().GetPing()...))
		err := mgmtorder.ValidateSitRep(sp, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.0.ping: invalid email length", err.GetStatusMessage())
	})
	tt.Run("sitrep.no_content", func(t *testing.T) {
		sp := sitrepObjWithID()
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
		o := orderObjWithIDs()
		o.SetTask(nil)
		err := mgmtorder.ValidateOrder(o, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.task", err.GetStatusMessage())
	})
	tt.Run("order.task.empty", func(t *testing.T) {
		o := orderObjWithIDs()
		o.SetTask(&order.Task{})
		err := mgmtorder.ValidateOrder(o, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.task: invalid task.id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("order.parent_order_id", func(t *testing.T) {
		o := orderObjWithIDs()
		o.SetParentOrderID("")
		err := mgmtorder.ValidateOrder(o, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.parent_order_id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("order.meta", func(t *testing.T) {
		o := orderObjWithIDs()
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
		o := orderObjWithIDs()
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
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		o.GetDelegatedTasks()[0].SetID(meta.NewID())
		resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "order.delegated.0.id disallowed", err.GetStatusMessage())
	})
	tt.Run("order.sitrep.id_provided", func(t *testing.T) {
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		o.GetSitReps()[0].SetID(meta.NewID())
		resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "order.sitrep.0.id disallowed", err.GetStatusMessage())
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
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		resp, err := s.API.PatchOrder(t, context.Background(), &request.PatchOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.task: invalid task.id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("order.delegated_task.id.empty", func(t *testing.T) {
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		oo := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, o)

		o.GetTask().SetID(oo.GetTask().GetID())
		o.GetDelegatedTasks()[0].SetObjective("patched objective")
		o.SetSitReps(nil)
		resp, err := s.API.PatchOrder(t, context.Background(), &request.PatchOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.delegated_task.0: invalid task.id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("order.sitrep.id.empty", func(t *testing.T) {
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		oo := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, o)

		o.GetTask().SetID(oo.GetTask().GetID())
		o.SetDelegatedTasks(nil)
		o.GetSitReps()[0].SetActions("patched actions")
		resp, err := s.API.PatchOrder(t, context.Background(), &request.PatchOrderRequest{
			Order: o,
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.sitrep.0: invalid sitrep.id: invalid id length", err.GetStatusMessage())
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
		resp, err := s.API.PutDelegatedTask(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PutDelegatedTask(t, context.Background(), &request.PutDelegatedTaskRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		resp, err := s.API.PutDelegatedTask(t, context.Background(), &request.PutDelegatedTaskRequest{
			Task: o.GetTask(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.task", func(t *testing.T) {
		resp, err := s.API.PutDelegatedTask(t, context.Background(), &request.PutDelegatedTaskRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty task", err.GetStatusMessage())
	})
	tt.Run("delegated_task.id.provided", func(t *testing.T) {
		o := orderObjWithIDs()
		oid := o.GetID()
		resp, err := s.API.PutDelegatedTask(t, context.Background(), &request.PutDelegatedTaskRequest{
			OrderID: oid,
			Task:    o.GetTask(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "task.id disallowed", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PatchDelegatedTaskRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PatchDelegatedTask(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PatchDelegatedTask(t, context.Background(), &request.PatchDelegatedTaskRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		resp, err := s.API.PatchDelegatedTask(t, context.Background(), &request.PatchDelegatedTaskRequest{
			Task: o.GetTask(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.task", func(t *testing.T) {
		resp, err := s.API.PatchDelegatedTask(t, context.Background(), &request.PatchDelegatedTaskRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty task", err.GetStatusMessage())
	})
	tt.Run("delegated_task.id.empty", func(t *testing.T) {
		o := orderObjWithIDs()
		oid := o.GetID()
		zeroOrderIDs(o)

		resp, err := s.API.PatchDelegatedTask(t, context.Background(), &request.PatchDelegatedTaskRequest{
			OrderID: oid,
			Task:    o.GetTask(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_DeleteDelegatedTaskRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTask(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTask(t, context.Background(), &request.DeleteDelegatedTaskRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTask(t, context.Background(), &request.DeleteDelegatedTaskRequest{
			DelegatedTaskID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.delegated_task_id", func(t *testing.T) {
		resp, err := s.API.DeleteDelegatedTask(t, context.Background(), &request.DeleteDelegatedTaskRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid delegated_task_id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PutSitRepRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PutSitRep(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PutSitRep(t, context.Background(), &request.PutSitRepRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.sitrep_id", func(t *testing.T) {
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		resp, err := s.API.PutSitRep(t, context.Background(), &request.PutSitRepRequest{
			SitRep: o.GetSitReps()[0],
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.sitrep", func(t *testing.T) {
		resp, err := s.API.PutSitRep(t, context.Background(), &request.PutSitRepRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitrep", err.GetStatusMessage())
	})
	tt.Run("sitrep.id.provided", func(t *testing.T) {
		o := orderObjWithIDs()
		oid := o.GetID()
		resp, err := s.API.PutSitRep(t, context.Background(), &request.PutSitRepRequest{
			OrderID: oid,
			SitRep:  o.GetSitReps()[0],
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "sitrep.id disallowed", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_PatchSitRepRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.PatchSitRep(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PatchSitRep(t, context.Background(), &request.PatchSitRepRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		o := orderObjWithIDs()
		zeroOrderIDs(o)
		resp, err := s.API.PatchSitRep(t, context.Background(), &request.PatchSitRepRequest{
			SitRep: o.GetSitReps()[0],
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.sitrep", func(t *testing.T) {
		resp, err := s.API.PatchSitRep(t, context.Background(), &request.PatchSitRepRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitrep", err.GetStatusMessage())
	})
	tt.Run("sitrep.id.empty", func(t *testing.T) {
		o := orderObjWithIDs()
		oid := o.GetID()
		zeroOrderIDs(o)

		resp, err := s.API.PatchSitRep(t, context.Background(), &request.PatchSitRepRequest{
			OrderID: oid,
			SitRep:  o.GetSitReps()[0],
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.id: invalid id length", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_DeleteSitRepRequest() {
	tt := s.T()
	tt.Run("nil.request", func(t *testing.T) {
		resp, err := s.API.DeleteSitRep(t, context.Background(), nil)
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.DeleteSitRep(t, context.Background(), &request.DeleteSitRepRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.order_id", func(t *testing.T) {
		resp, err := s.API.DeleteSitRep(t, context.Background(), &request.DeleteSitRepRequest{
			SitRepID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order_id: invalid id length", err.GetStatusMessage())
	})
	tt.Run("empty.sitrep_id", func(t *testing.T) {
		resp, err := s.API.DeleteSitRep(t, context.Background(), &request.DeleteSitRepRequest{
			OrderID: meta.NewID(),
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep_id: invalid id length", err.GetStatusMessage())
	})
}
