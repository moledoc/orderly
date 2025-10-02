package mgmtorder

import (
	"context"

	"github.com/moledoc/orderly/actions"
	"github.com/moledoc/orderly/middleware"
	"github.com/moledoc/orderly/models"
)

func HandlePostOrder(ctx context.Context, order *models.Order) (*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "handlePostOrder")
	defer middleware.SpanStop(ctx, "handlePostOrder")

	// TODO: validation

	return strg.Write(ctx, actions.CREATE, order)
}

func HandleGetOrderByID(ctx context.Context, orderID uint) (*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "handleGetOrderByID")
	defer middleware.SpanStop(ctx, "handleGetOrderByID")

	// TODO: validation

	resp, err := strg.Read(ctx, actions.READ, orderID)
	if len(resp) >= 1 {
		return resp[0], err
	}
	return nil, err
}

func HandleGetOrders(ctx context.Context) ([]*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleGetOrders")
	defer middleware.SpanStop(ctx, "HandleGetOrders")

	// TODO: validation

	return strg.Read(ctx, actions.READALL, 0)
}

func HandleGetOrderVersions(ctx context.Context, orderID uint) ([]*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleGetOrderVersions")
	defer middleware.SpanStop(ctx, "HandleGetOrderVersions")

	// TODO: validation

	return strg.Read(ctx, actions.READVERSIONS, orderID)
}

func HandleGetOrderSubOrders(ctx context.Context, orderID uint) ([]*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleGetOrderSubOrders")
	defer middleware.SpanStop(ctx, "HandleGetOrderSubOrders")

	// TODO: validation
	return strg.Read(ctx, actions.READSUBORDERS, orderID)
}

func HandlePatchOrder(ctx context.Context, order *models.Order) (*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandlePatchOrder")
	defer middleware.SpanStop(ctx, "HandlePatchOrder")

	// TODO: validation

	return strg.Write(ctx, actions.UPDATE, order)
}

func HandleDeleteOrderHard(ctx context.Context, order *models.Order) (*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleDeleteOrder:Hard")
	defer middleware.SpanStop(ctx, "HandleDeleteOrder:Hard")

	// TODO: validation

	return strg.Write(ctx, actions.DELETEHARD, order)

}

func HandleDeleteOrderSoft(ctx context.Context, order *models.Order) (*models.Order, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleDeleteOrder:Soft")
	defer middleware.SpanStop(ctx, "HandleDeleteOrder:Soft")

	// TODO: validation

	return strg.Write(ctx, actions.DELETESOFT, order)

}
