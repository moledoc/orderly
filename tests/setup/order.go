package setup

import (
	"context"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/api"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/stretchr/testify/require"
)

func TaskObj(extra ...string) *order.Task {
	return &order.Task{
		State:       utils.Ptr(order.NotStarted),
		Accountable: UserObjWithID(append(extra, "accountable")...),
		Objective:   "objective description",
		Deadline:    time.Now().UTC(),
	}
}
func TaskObjWithID(extra ...string) *order.Task {
	tt := TaskObj(extra...)
	tt.SetID(meta.NewID())
	return tt
}

func SitrepObj(extra ...string) *order.SitRep {
	return &order.SitRep{
		DateTime:  time.Now().UTC(),
		By:        UserObjWithID(append(extra, "by")...),
		Situation: "situation description",
		Actions:   "list of actions taken",
		TODO:      "list of things to do still",
		Issues:    "list of encountered issues",
	}
}

func SitrepObjWithID(extra ...string) *order.SitRep {
	sr := SitrepObj(extra...)
	sr.SetID(meta.NewID())
	return sr
}

func OrderObj(extra ...string) *order.Order {
	return &order.Order{
		Task: TaskObj(extra...),
		DelegatedTasks: []*order.Task{
			TaskObj(),
			TaskObj(),
			TaskObj(),
		},
		ParentOrderID: meta.NewID(),
		SitReps: []*order.SitRep{
			SitrepObj(),
			SitrepObj(),
			SitrepObj(),
		},
		Meta: &meta.Meta{
			Version: 1,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		},
	}
}

func OrderObjWithIDs(extra ...string) *order.Order {
	return &order.Order{
		Task: TaskObjWithID(extra...),
		DelegatedTasks: []*order.Task{
			TaskObjWithID(),
			TaskObjWithID(),
			TaskObjWithID(),
		},
		ParentOrderID: meta.NewID(),
		SitReps: []*order.SitRep{
			SitrepObjWithID(),
			SitrepObjWithID(),
			SitrepObjWithID(),
		},
		Meta: &meta.Meta{
			Version: 1,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		},
	}
}

func ZeroOrderIDs(o *order.Order) {
	o.GetTask().SetID("")
	for _, delegated := range o.GetDelegatedTasks() {
		delegated.SetID("")
	}
	for _, sitrep := range o.GetSitReps() {
		sitrep.SetID("")
	}
}

func CreateOrderWithCleanup(t *testing.T, ctx context.Context, api api.Order, orderObj *order.Order) (*order.Order, errwrap.Error) {
	resp, err := api.PostOrder(t, ctx, &request.PostOrderRequest{
		Order: orderObj,
	})
	if err != nil {
		return nil, err
	}
	cleanup.Order(t, api, resp.GetOrder())
	return resp.GetOrder(), nil
}

func MustCreateOrderWithCleanup(t *testing.T, ctx context.Context, api api.Order, orderObj *order.Order) *order.Order {
	order, err := CreateOrderWithCleanup(t, ctx, api, orderObj)
	require.NoError(t, err)
	return order
}
