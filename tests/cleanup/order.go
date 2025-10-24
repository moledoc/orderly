package cleanup

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/tests/api"
)

func Order(t *testing.T, api api.Order, o *order.Order) {
	t.Helper()
	t.Cleanup(func() {
		id := meta.ID(o.GetOrder().GetID())
		if len(id) == 0 {
			return
		}
		_, err := api.DeleteOrder(t, context.Background(), &request.DeleteOrderRequest{
			ID: id,
		})
		if err != nil {
			t.Logf("[WARNING]: order '%v' wasn't cleaned up: %s\n", id, err)
		}
	})
}
