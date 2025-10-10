package tests

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestPostOrder_Minimal() {
	tt := s.T()

	o := setup.OrderObj()
	o.SetDelegatedTasks([]*order.Task{})
	o.SetSitReps([]*order.SitRep{})
	resp, err := s.API.PostOrder(tt, context.Background(), &request.PostOrderRequest{
		Order: o,
	})
	defer cleanup.Order(tt, s.API, resp.GetOrder())
	require.NoError(tt, err)

	opts := []cmp.Option{
		compare.IgnoreID,
		compare.IgnoreMeta,
	}

	expected := &response.PostOrderResponse{
		Order: o,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
}

func (s *OrderSuite) TestPostOrder_WithDelegatedTasks() {
	tt := s.T()

	o := setup.OrderObj()
	o.SetSitReps([]*order.SitRep{})
	resp, err := s.API.PostOrder(tt, context.Background(), &request.PostOrderRequest{
		Order: o,
	})
	defer cleanup.Order(tt, s.API, resp.GetOrder())
	require.NoError(tt, err)

	opts := []cmp.Option{
		compare.IgnoreID,
		compare.IgnoreMeta,
	}

	expected := &response.PostOrderResponse{
		Order: o,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
}

func (s *OrderSuite) TestPostOrder_WithSitReps() {
	tt := s.T()

	o := setup.OrderObj()
	o.SetSitReps([]*order.SitRep{})
	resp, err := s.API.PostOrder(tt, context.Background(), &request.PostOrderRequest{
		Order: o,
	})
	defer cleanup.Order(tt, s.API, resp.GetOrder())
	require.NoError(tt, err)

	opts := []cmp.Option{
		compare.IgnoreID,
		compare.IgnoreMeta,
	}

	expected := &response.PostOrderResponse{
		Order: o,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
}

func (s *OrderSuite) TestPostOrder_Full() {
	tt := s.T()

	o := setup.OrderObj()
	resp, err := s.API.PostOrder(tt, context.Background(), &request.PostOrderRequest{
		Order: o,
	})
	defer cleanup.Order(tt, s.API, resp.GetOrder())
	require.NoError(tt, err)

	opts := []cmp.Option{
		compare.IgnoreID,
		compare.IgnoreMeta,
	}

	expected := &response.PostOrderResponse{
		Order: o,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
}
