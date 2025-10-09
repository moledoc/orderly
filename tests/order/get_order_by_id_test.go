package tests

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestGetOrderByID() {
	tt := s.T()
	tt.SkipNow()

	o := orderObj()
	zeroOrderIDs(o)
	o.SetDelegatedTasks(nil)
	o.SetSitReps(nil)
	resp, err := s.API.GetOrderByID(tt, context.Background(), &request.GetOrderByIDRequest{
		ID: meta.NewID(),
	})
	defer cleanup.Order(tt, s.API, resp.GetOrder())
	require.NoError(tt, err)

	opts := []cmp.Option{
		compare.IgnoreID,
	}

	expected := &response.GetOrderByIDResponse{
		Order: o,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
	require.NotEmpty(tt, resp.GetOrder().GetMeta())
}
