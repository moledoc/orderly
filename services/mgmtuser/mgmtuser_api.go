package mgmtuser

import (
	"context"

	"github.com/moledoc/orderly/actions"
	"github.com/moledoc/orderly/middleware"
	"github.com/moledoc/orderly/models"
)

func HandlePostUser(ctx context.Context, user *models.User) (*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandlePostUser")
	defer middleware.SpanStop(ctx, "HandlePostUser")

	// TODO: validation

	return strg.Write(ctx, actions.CREATE, user)
}

func HandleGetUserByID(ctx context.Context, userID uint) ([]*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleGetUserByID")
	defer middleware.SpanStop(ctx, "HandleGetUserByID")

	// TODO: validation

	return strg.Read(ctx, actions.READ, userID)
}

func HandleGetUsers(ctx context.Context) ([]*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "handleGetUsers")
	defer middleware.SpanStop(ctx, "handleGetUsers")

	// TODO: validation

	return strg.Read(ctx, actions.READALL, 0)
}

func HandleGetUserVersions(ctx context.Context, userID uint) ([]*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleGetUserVersions")
	defer middleware.SpanStop(ctx, "HandleGetUserVersions")

	// TODO: validation

	return strg.Read(ctx, actions.READVERSIONS, userID)
}

func HandleGetUserSubOrdinates(ctx context.Context, userID uint) ([]*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleGetUserSubOrdinates")
	defer middleware.SpanStop(ctx, "HandleGetUserSubOrdinates")

	// TODO: validation

	return strg.Read(ctx, actions.READSUBORDINATES, userID)
}

func HandlePatchUser(ctx context.Context, user *models.User) (*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandlePatchUser")
	defer middleware.SpanStop(ctx, "HandlePatchUser")

	// TODO: validation

	return strg.Write(ctx, actions.UPDATE, user)
}

func HandleDeleteUserHard(ctx context.Context, user *models.User) (*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleDeleteUser:Hard")
	defer middleware.SpanStop(ctx, "HandleDeleteUser:Hard")

	// TODO: valiation

	return strg.Write(ctx, actions.DELETEHARD, user)
}

func HandleDeleteUserSoft(ctx context.Context, user *models.User) (*models.User, models.IError) {
	ctx = middleware.AddTrace(ctx, nil)

	middleware.SpanStart(ctx, "HandleDeleteUser:Soft")
	defer middleware.SpanStop(ctx, "HandleDeleteUser:Soft")

	// TODO: valiation

	return strg.Write(ctx, actions.DELETESOFT, user)
}
