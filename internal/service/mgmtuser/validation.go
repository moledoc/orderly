package mgmtuser

import (
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/service/common/validation"
)

func validateUser(user *user.User) errwrap.Error {
	if user == nil {
		return nil
	}

	if len(user.GetID()) > 0 { // NOTE: ID is required, but when creating we don't allow ID; relevant ID check is done one level up in validation
		err := validation.ValidateID(user.GetID())
		if err != nil {
			return err
		}
	}

	if len(user.GetName()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid user.name length")
	}

	err := validation.ValidateEmail(user.GetEmail())
	if err != nil {
		return err
	}

	err = validation.ValidateEmail(user.GetSupervisor())
	if err != nil {
		return err
	}

	err = validation.ValidateMeta(user.GetMeta())
	if err != nil {
		return err
	}

	return nil
}

func validatePostUserRequest(req *request.PostUserRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}
	if req.User == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty user")
	}

	if len(req.GetUser().GetID()) > 0 {
		return errwrap.NewError(http.StatusBadRequest, "user.id disallowed")
	}

	if req.GetUser().Meta != nil {
		return errwrap.NewError(http.StatusBadRequest, "user.meta disallowed")
	}

	return validateUser(req.GetUser())
}

func validateGetUserByIDRequest(req *request.GetUserByIDRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}

	return nil
}

func validateGetUsersRequest(*request.GetUsersRequest) errwrap.Error {
	return nil
}

func validateGetUserSubOrdinatesRequest(req *request.GetUserSubOrdinatesRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}

	return nil
}

func validatePatchUserRequest(req *request.PatchUserRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}
	if req.User == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty user")
	}

	if len(req.GetUser().GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "user.id missing")
	}

	if err := validateUser(req.GetUser()); err != nil {
		return err
	}

	return nil
}

func validateDeleteUserRequest(req *request.DeleteUserRequest) errwrap.Error {
	if req == nil {
		return errwrap.NewError(http.StatusBadRequest, "empty request")
	}

	if len(req.GetID()) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "empty id")
	}

	err := validation.ValidateID(req.GetID())
	if err != nil {
		return err
	}

	return nil
}
