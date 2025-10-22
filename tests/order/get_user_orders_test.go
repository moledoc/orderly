package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestGetUserOrders() {
	tt := s.T()

	createOrdersWithAccountable := func(t *testing.T, count int, accountable *user.User) []*order.Order {
		if count == 0 {
			return nil
		}
		orders := make([]*order.Order, count)
		for i := 1; i <= count; i++ {
			o := setup.OrderObj(fmt.Sprintf("%v", i))
			o.GetTask().SetAccountable(accountable)
			order := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, o)
			orders[i-1] = order
		}
		return orders
	}

	for _, i := range []int{0, 1, 10} {
		tt.Run(fmt.Sprintf("count.%d", i), func(t *testing.T) {

			user := setup.MustCreateUserWithCleanup(t, context.Background(), s.UserAPI, setup.UserObj())
			orders := createOrdersWithAccountable(t, i, user)

			user2 := setup.MustCreateUserWithCleanup(t, context.Background(), s.UserAPI, setup.UserObj())
			_ = createOrdersWithAccountable(t, 5, user2)

			resp, err := s.API.GetUserOrders(t, context.Background(), &request.GetUserOrdersRequest{
				UserID: user.GetID(),
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.SorterOrder(compare.SortOrderByID),
			}

			expected := &response.GetUserOrdersResponse{
				Orders: orders,
			}
			compare.RequireEqual(t, expected, resp, opts...)
		})
	}
}
