package setup

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/tests/api"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/stretchr/testify/require"
)

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
