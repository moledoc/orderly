package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestGetOrderSubOrders() {
	tt := s.T()

	parentOrder := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj("parent-order"))

	createSubOrders := func(t *testing.T, count int) []*order.Order {
		if count == 0 {
			return nil
		}
		orders := make([]*order.Order, count)
		for i := 1; i <= count; i++ {

			o := setup.OrderObj(fmt.Sprintf("%v", i))
			o.SetParentOrderID(parentOrder.GetID())
			order := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, o)

			orders[i-1] = order
		}
		return orders
	}

	for _, i := range []int{0, 1, 10} {
		tt.Run(fmt.Sprintf("count.%d", i), func(t *testing.T) {

			orders := createSubOrders(t, i)

			resp, err := s.API.GetOrderSubOrders(t, context.Background(), &request.GetOrderSubOrdersRequest{
				ID: parentOrder.GetID(),
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.SorterOrder(compare.SortOrderByID),
			}

			expected := &response.GetOrderSubOrdersResponse{
				SubOrders: orders,
			}
			compare.RequireEqual(t, expected, resp, opts...)
			for _, u := range resp.GetSubOrders() {
				assert.NotEmpty(tt, u.GetMeta())
			}
		})
	}
}
