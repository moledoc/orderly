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

func (s *OrderSuite) TestPutSitReps() {
	tt := s.T()

	createSitReps := func(count int) []*order.SitRep {
		if count == 0 {
			return nil
		}
		tasks := make([]*order.SitRep, count)
		for i := 1; i <= count; i++ {
			tasks[i-1] = setup.SitrepObj(fmt.Sprintf("%v-%v", count, i))
		}
		return tasks
	}

	o := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj())

	for _, i := range []int{1, 10} {
		tt.Run(fmt.Sprintf("count.%v", i), func(t *testing.T) {
			sitreps := createSitReps(i)

			respPatch, err := s.API.PutSitReps(t, context.Background(), &request.PutSitRepsRequest{
				OrderID: o.GetID(),
				SitReps: sitreps,
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.IgnoreID,
				compare.IgnoreMeta,
				compare.SorterOrder(compare.SortOrderByID),
				compare.SorterTask(compare.SortTaskByID),
				compare.SorterSitRep(compare.SortSitRepByDateTime),
			}

			o.SetSitReps(append(o.GetSitReps(), sitreps...))
			expectedPut := &response.PutSitRepsResponse{
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
