package tests

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestGetOrderByID() {
	tt := s.T()

	o := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj())
	resp, err := s.API.GetOrderByID(tt, context.Background(), &request.GetOrderByIDRequest{
		ID: o.GetID(),
	})
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

func (s *OrderSuite) TestGetOrderByID_Failed() {

	tt := s.T()

	tt.Run("NotFound", func(t *testing.T) {
		resp, err := s.API.GetOrderByID(tt, context.Background(), &request.GetOrderByIDRequest{
			ID: meta.NewID(),
		})
		defer cleanup.Order(tt, s.API, resp.GetOrder())
		require.Error(t, err)
		require.Empty(t, resp)
	})
}
