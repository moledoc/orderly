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

func (s *OrderSuite) TestPutDelegatedTask() {
	tt := s.T()

	createTasks := func(count int) []*order.Task {
		if count == 0 {
			return nil
		}
		tasks := make([]*order.Task, count)
		for i := 1; i <= count; i++ {
			tasks[i-1] = setup.TaskObj(fmt.Sprintf("%v-%v", count, i))
		}
		return tasks
	}

	o := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj())

	for _, i := range []int{1, 10} {
		tt.Run(fmt.Sprintf("count.%v", i), func(t *testing.T) {
			tasks := createTasks(i)

			respPatch, err := s.API.PutDelegatedTask(t, context.Background(), &request.PutDelegatedTasksRequest{
				OrderID: o.GetID(),
				Tasks:   tasks,
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.IgnoreID,
				compare.IgnoreMeta,
				compare.SorterOrder(compare.SortOrderByID),
				compare.SorterTask(compare.SortTaskByAccountable),
				compare.SorterSitRep(compare.SortSitRepByID),
			}

			o.SetDelegatedTasks(append(o.GetDelegatedTasks(), tasks...))
			expectedPut := &response.PutDelegatedTasksResponse{
				Order: o,
			}
			compare.RequireEqual(t, expectedPut, respPatch, opts...)

			respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
				ID: o.GetID(),
			})
			expectedGet := &response.GetOrderByIDResponse{
				Order: o,
			}
			require.NoError(t, err)
			compare.RequireEqual(t, expectedGet, respGet, opts...)
		})
	}
}
