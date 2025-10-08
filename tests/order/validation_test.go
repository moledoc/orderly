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
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/stretchr/testify/require"
)

func taskObj(extra ...string) *order.Task {
	ee := strings.Join(extra, ".")
	return &order.Task{
		ID:          meta.NewID(),
		State:       order.NotStarted,
		Accountable: user.Email(fmt.Sprintf("example%v@example.com", ee)),
		Objective:   "objective description",
		Deadline:    time.Now().UTC(),
	}
}

func sitrepObj(extra ...string) *order.SitRep {
	ee := strings.Join(extra, ".")
	return &order.SitRep{
		ID:            meta.NewID(),
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
		to := taskObj()
		to.SetID("")
		err := mgmtorder.ValidateTask(to)
		require.NoError(t, err) // NOTE: only being validated if len(task.id) > 0; it's to enable to use common task validation func across endpoints. task.ID checks are done in request validation
	})
	tt.Run("task.state.low", func(t *testing.T) {
		to := taskObj()
		to.SetState(order.NotStarted - 1)
		err := mgmtorder.ValidateTask(to)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.state", err.GetStatusMessage())
	})
	tt.Run("task.state.high", func(t *testing.T) {
		to := taskObj()
		to.SetState(order.Completed + 1)
		err := mgmtorder.ValidateTask(to)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.state", err.GetStatusMessage())
	})
	tt.Run("task.accountable", func(t *testing.T) {
		to := taskObj()
		to.SetAccountable("")
		err := mgmtorder.ValidateTask(to)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid email length", err.GetStatusMessage())
	})
	tt.Run("task.objective", func(t *testing.T) {
		to := taskObj()
		to.SetObjective("")
		err := mgmtorder.ValidateTask(to)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.objective", err.GetStatusMessage())
	})
	tt.Run("task.deadline", func(t *testing.T) {
		to := taskObj()
		to.SetDeadline(time.Time{})
		err := mgmtorder.ValidateTask(to)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.deadline", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_SitRep() {
	tt := s.T()

	tt.Run("sitrep.id", func(t *testing.T) {
		sp := sitrepObj()
		sp.SetID("")
		err := mgmtorder.ValidateSitRep(sp)
		require.NoError(t, err) // NOTE: only being validated if len(sitrep.id) > 0; it's to enable to use common sitrep validation func across endpoints. sitrep.ID checks are done in request validation
	})
	tt.Run("sitrep.datetime", func(t *testing.T) {
		sp := sitrepObj()
		sp.SetDateTime(time.Time{})
		err := mgmtorder.ValidateSitRep(sp)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.datetime", err.GetStatusMessage())
	})
	tt.Run("sitrep.by", func(t *testing.T) {
		sp := sitrepObj()
		sp.SetBy("")
		err := mgmtorder.ValidateSitRep(sp)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.by: invalid email length", err.GetStatusMessage())
	})
	tt.Run("sitrep.ping", func(t *testing.T) {
		sp := sitrepObj()
		sp.SetPing(append([]user.Email{""}, sitrepObj().GetPing()...))
		err := mgmtorder.ValidateSitRep(sp)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid sitrep.0.ping: invalid email length", err.GetStatusMessage())
	})
	tt.Run("sitrep.no_content", func(t *testing.T) {
		sp := sitrepObj()
		sp.SetSituation("")
		sp.SetActions("")
		sp.SetTBD("")
		sp.SetIssues("")
		err := mgmtorder.ValidateSitRep(sp)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty sitrep", err.GetStatusMessage())
	})
}

func (s *OrderSuite) TestValidation_Order() {
	tt := s.T()

	tt.Run("order.task.nil", func(t *testing.T) {
		o := orderObj()
		o.SetTask(nil)
		err := mgmtorder.ValidateOrder(o)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.task", err.GetStatusMessage())
	})
	tt.Run("order.task.empty", func(t *testing.T) {
		o := orderObj()
		o.SetTask(&order.Task{})
		err := mgmtorder.ValidateOrder(o)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid task.state", err.GetStatusMessage())
	})
	tt.Run("order.parent_order_id", func(t *testing.T) {
		o := orderObj()
		o.SetParentOrderID("")
		err := mgmtorder.ValidateOrder(o)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid order.parent_order_id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("order.meta", func(t *testing.T) {
		o := orderObj()
		o.SetMeta(&meta.Meta{
			Version: 2,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		})
		err := mgmtorder.ValidateOrder(o)
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
		o := orderObj()
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
		o := orderObj()
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
		o := orderObj()
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
