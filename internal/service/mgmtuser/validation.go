package mgmtuser

import (
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/service/common/validation"
)

func ValidateUser(user *user.User, ignore validation.IgnoreField) errwrap.Error {
	if user == nil {
		return nil
	}

	err := validation.ValidateID(user.GetID())
	if !validation.IsFieldIgnored(validation.IgnoreID, ignore) && err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid user.id: %s", err.GetStatusMessage())
	}

	if len(user.GetName()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid user.name length")
	}

	if err := validation.ValidateEmail(user.GetEmail()); err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid user.email: %s", err.GetStatusMessage())
	}

	if err := validation.ValidateEmail(user.GetSupervisor()); err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid user.supervisor: %s", err.GetStatusMessage())
	}

	return nil
}

func ValidatePostUserRequest(req *request.PostUserRequest) errwrap.Error {
	if req.User == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetUser().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "user.id disallowed")
	}

	if req.GetUser().Meta != nil {
		return errwrap.NewError(http.StatusBadRequest, "user.meta disallowed")
	}

	return ValidateUser(req.GetUser(), validation.IgnoreID)
}

func ValidateGetUserByIDRequest(req *request.GetUserByIDRequest) errwrap.Error {
	if req == nil || len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}

	return nil
}

func ValidateGetUsersRequest(*request.GetUsersRequest) errwrap.Error {
	return nil
}

func ValidateGetUserSubOrdinatesRequest(req *request.GetUserSubOrdinatesRequest) errwrap.Error {
	if req == nil || len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}

	return nil
}

func ValidatePatchUserRequest(req *request.PatchUserRequest) errwrap.Error {

	if req.User == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty user")
	}

	if len(req.GetUser().GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "user.id missing")
	}

	return ValidateUser(req.GetUser(), validation.IgnoreNothing)
}

func ValidateDeleteUserRequest(req *request.DeleteUserRequest) errwrap.Error {
	if req == nil || len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}

	return nil
}
