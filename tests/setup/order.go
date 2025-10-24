package setup

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/api"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/stretchr/testify/require"
)

func OrderObj(extra ...string) *order.Order {
	ee := strings.Join(append([]string{""}, extra...), ".")
	return &order.Order{
		State:       utils.Ptr(order.NotStarted),
		Accountable: user.Email(fmt.Sprintf("example%v@example.com", ee)),
		Objective:   "objective description",
		Deadline:    time.Now().UTC(),
	}
}
func OrderObjWithID(extra ...string) *order.Order {
	tt := OrderObj(extra...)
	tt.SetID(meta.NewID())
	return tt
}

func SitrepObj(extra ...string) *order.SitRep {
	ee := strings.Join(append([]string{""}, extra...), ".")
	return &order.SitRep{
		DateTime:  time.Now().UTC(),
		By:        user.Email(fmt.Sprintf("example%v@example.com", ee)),
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
		Order: OrderObj(extra...),
		DelegatedOrders: []*order.Order{
			OrderObj(),
			OrderObj(),
			OrderObj(),
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
		Order: OrderObjWithID(extra...),
		DelegatedOrders: []*order.Order{
			OrderObjWithID(),
			OrderObjWithID(),
			OrderObjWithID(),
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
	o.GetOrder().SetID("")
	for _, delegated := range o.GetDelegatedOrders() {
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
