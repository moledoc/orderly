package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestDeleteOrder() {
	tt := s.T()

	tt.Run("exists.yes", func(t *testing.T) {
		o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())

		respDelete, err := s.API.DeleteOrder(t, context.Background(), &request.DeleteOrderRequest{
			ID: o.GetID(),
		})
		require.NoError(t, err)
		require.Empty(t, respDelete)

		respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
			ID: o.GetID(),
		})
		require.Error(t, err)
		require.Empty(t, respGet)
		require.Equal(t, http.StatusNotFound, err.GetStatusCode(), err)
		require.Equal(t, "not found", err.GetStatusMessage())
	})

	tt.Run("exists.no", func(t *testing.T) {
		id := meta.NewID()
		respDelete, err := s.API.DeleteOrder(t, context.Background(), &request.DeleteOrderRequest{
			ID: id,
		})
		require.NoError(t, err)
		require.Empty(t, respDelete)

		respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
			ID: id,
		})
		require.Error(t, err)
		require.Empty(t, respGet)
		require.Equal(t, http.StatusNotFound, err.GetStatusCode(), err)
		require.Equal(t, "not found", err.GetStatusMessage())
	})

}
