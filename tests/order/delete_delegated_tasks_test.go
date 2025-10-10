package tests

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestDeleteDelegatedTasks() {
	tt := s.T()

	tt.Run("delegated_task_ids.no_match", func(t *testing.T) {
		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())

		respDelete, err := s.API.DeleteDelegatedTasks(t, context.Background(), &request.DeleteDelegatedTasksRequest{
			OrderID:          o.GetID(),
			DelegatedTaskIDs: []meta.ID{meta.NewID(), meta.NewID(), meta.NewID()},
		})
		require.NoError(t, err)
		expectedDelete := &response.DeleteDelegatedTasksResponse{
			Order: o,
		}

		opts := []cmp.Option{
			compare.SorterTask(compare.SortTaskByID),
		}
		compare.RequireEqual(t, expectedDelete, respDelete, opts...)

		respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
			ID: o.GetID(),
		})
		require.NoError(t, err)

		expectedGet := &response.GetOrderByIDResponse{
			Order: o,
		}
		compare.RequireEqual(t, expectedGet, respGet, opts...)
	})

	tt.Run("delegated_task_ids.1_match", func(t *testing.T) {
		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
		if len(o.GetDelegatedTasks()) == 0 {
			t.Logf("unexpected delegated_tasks length")
			t.FailNow()
		}

		respDelete, err := s.API.DeleteDelegatedTasks(t, context.Background(), &request.DeleteDelegatedTasksRequest{
			OrderID:          o.GetID(),
			DelegatedTaskIDs: []meta.ID{o.GetDelegatedTasks()[0].GetID()},
		})
		require.NoError(t, err)

		o.SetDelegatedTasks(o.GetDelegatedTasks()[1:])
		o.GetMeta().VersionIncr()
		expectedDelete := &response.DeleteDelegatedTasksResponse{
			Order: o,
		}
		opts := []cmp.Option{
			compare.SorterTask(compare.SortTaskByID),
			compare.IgnorePaths("Order.Meta.Updated"),
		}
		compare.RequireEqual(t, expectedDelete, respDelete, opts...)

		respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
			ID: o.GetID(),
		})
		require.NoError(t, err)

		expectedGet := &response.GetOrderByIDResponse{
			Order: o,
		}
		compare.RequireEqual(t, expectedGet, respGet, opts...)
	})

	tt.Run("delegated_task_ids.all_match", func(t *testing.T) {
		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())

		var ids []meta.ID
		for _, delegated := range o.GetDelegatedTasks() {
			ids = append(ids, delegated.GetID())
		}
		respDelete, err := s.API.DeleteDelegatedTasks(t, context.Background(), &request.DeleteDelegatedTasksRequest{
			OrderID:          o.GetID(),
			DelegatedTaskIDs: ids,
		})
		require.NoError(t, err)

		o.SetDelegatedTasks(nil)
		o.GetMeta().VersionIncr()
		expectedDelete := &response.DeleteDelegatedTasksResponse{
			Order: o,
		}
		opts := []cmp.Option{
			compare.IgnorePaths("Order.Meta.Updated"),
		}
		compare.RequireEqual(t, expectedDelete, respDelete, opts...)

		respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
			ID: o.GetID(),
		})
		require.NoError(t, err)

		expectedGet := &response.GetOrderByIDResponse{
			Order: o,
		}
		compare.RequireEqual(t, expectedGet, respGet, opts...)
	})

}
