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

func (s *OrderSuite) TestGetOrders() {
	tt := s.T()

	createOrders := func(t *testing.T, count int) []*order.Order {
		if count == 0 {
			return nil
		}
		orders := make([]*order.Order, count)
		for i := 1; i <= count; i++ {
			orders[i-1] = setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, orderObj())
		}
		return orders
	}

	for _, i := range []int{0, 1, 10} {
		tt.Run(fmt.Sprintf("count.%v", i), func(t *testing.T) {
			orders := createOrders(t, i)

			resp, err := s.API.GetOrders(t, context.Background(), &request.GetOrdersRequest{})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.SorterOrder(compare.SortOrderByID),
			}

			expected := &response.GetOrdersResponse{
				Orders: orders,
			}
			compare.RequireEqual(t, expected, resp, opts...)
			for _, u := range resp.GetOrders() {
				assert.NotEmpty(t, u.GetMeta())
			}
		})
	}
}
