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
			orders[i-1] = setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
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
		})
	}
}

func (s *OrderSuite) TestGetOrders_ByAccountable() {
	tt := s.T()

	createOrders := func(t *testing.T, count int) []*order.Order {
		if count == 0 {
			return nil
		}
		orders := make([]*order.Order, count)
		for i := 1; i <= count; i++ {
			orders[i-1] = setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
			setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj("should-not-get-these-orders"))
		}
		return orders
	}

	for _, i := range []int{0, 1, 10} {
		tt.Run(fmt.Sprintf("count.%v", i), func(t *testing.T) {
			orders := createOrders(t, i)

			resp, err := s.API.GetOrders(t, context.Background(), &request.GetOrdersRequest{
				Accountable: setup.OrderObj().GetTask().GetAccountable(),
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.SorterOrder(compare.SortOrderByID),
			}

			expected := &response.GetOrdersResponse{
				Orders: orders,
			}
			compare.RequireEqual(t, expected, resp, opts...)
		})
	}
}

func (s *OrderSuite) TestGetOrders_ByParentOrderID() {
	tt := s.T()

	parent := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj())
	createOrders := func(t *testing.T, count int) []*order.Order {
		if count == 0 {
			return nil
		}
		orders := make([]*order.Order, count)
		for i := 1; i <= count; i++ {
			obj := setup.OrderObj()
			obj.SetParentOrderID(parent.GetID())
			orders[i-1] = setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, obj)
			setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj("should-not-get-these-orders"))
		}
		return orders
	}

	for _, i := range []int{0, 1, 10} {
		tt.Run(fmt.Sprintf("count.%v", i), func(t *testing.T) {
			orders := createOrders(t, i)

			resp, err := s.API.GetOrders(t, context.Background(), &request.GetOrdersRequest{
				ParentOrderID: parent.GetID(),
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.SorterOrder(compare.SortOrderByID),
			}

			expected := &response.GetOrdersResponse{
				Orders: orders,
			}
			compare.RequireEqual(t, expected, resp, opts...)
		})
	}
}
