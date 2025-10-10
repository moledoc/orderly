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

func (s *OrderSuite) TestDeleteSitReps() {
	tt := s.T()

	tt.Run("sitrep_ids.no_match", func(t *testing.T) {
		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())

		respDelete, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{
			OrderID:   o.GetID(),
			SitRepIDs: []meta.ID{meta.NewID(), meta.NewID(), meta.NewID()},
		})
		require.NoError(t, err)
		expectedDelete := &response.DeleteSitRepsResponse{
			Order: o,
		}

		opts := []cmp.Option{
			compare.SorterSitRep(compare.SortSitRepByID),
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

	tt.Run("sitrep_ids.1_match", func(t *testing.T) {
		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
		if len(o.GetSitReps()) == 0 {
			t.Logf("unexpected sitreps length")
			t.FailNow()
		}

		respDelete, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{
			OrderID:   o.GetID(),
			SitRepIDs: []meta.ID{o.GetSitReps()[0].GetID()},
		})
		require.NoError(t, err)

		o.SetSitReps(o.GetSitReps()[1:])
		o.GetMeta().VersionIncr()
		expectedDelete := &response.DeleteSitRepsResponse{
			Order: o,
		}
		opts := []cmp.Option{
			compare.SorterSitRep(compare.SortSitRepByID),
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

	tt.Run("sitrep_ids.all_match", func(t *testing.T) {
		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())

		var ids []meta.ID
		for _, sitrep := range o.GetSitReps() {
			ids = append(ids, sitrep.GetID())
		}
		respDelete, err := s.API.DeleteSitReps(t, context.Background(), &request.DeleteSitRepsRequest{
			OrderID:   o.GetID(),
			SitRepIDs: ids,
		})
		require.NoError(t, err)

		o.SetSitReps(nil)
		o.GetMeta().VersionIncr()
		expectedDelete := &response.DeleteSitRepsResponse{
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
