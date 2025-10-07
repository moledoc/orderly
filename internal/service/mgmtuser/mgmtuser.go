package mgmtuser

import (
	"context"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/pkg/utils"
)

func (s *serviceMgmtUser) PostUser(ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PostUser")
	defer middleware.SpanStop(ctx, "PostUser")

	if err := validatePostUserRequest(req); err != nil {
		return nil, err
	}

	u := req.GetUser()
	u.ID = meta.ID(utils.RandAlphanum())

	now := time.Now().UTC()
	u.Meta = &meta.Meta{
		Version: 1,
		Created: now,
		Updated: now,
	}

	user, err := s.Repository.Write(ctx, u)
	if err != nil {
		return nil, err
	}

	return &response.PostUserResponse{
		User: user,
	}, nil
}

func (s *serviceMgmtUser) GetUserByID(ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetUserByID")
	defer middleware.SpanStop(ctx, "GetUserByID")

	if err := validateGetUserByIDRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.Repository.ReadByID(ctx, req.GetID())
	if err != nil {
		return nil, err
	}
	return &response.GetUserByIDResponse{
		User: resp,
	}, nil
}

func (s *serviceMgmtUser) GetUsers(ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetUserByID")
	defer middleware.SpanStop(ctx, "GetUserByID")

	if err := validateGetUsersRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.Repository.ReadAll(ctx)
	if err != nil {
		return nil, err
	}
	return &response.GetUsersResponse{
		Users: resp,
	}, nil
}

func (s *serviceMgmtUser) GetUserSubOrdinates(ctx context.Context, req *request.GetUserSubOrdinatesRequest) (*response.GetUserSubOrdinatesResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "GetUserSubOrdinates")
	defer middleware.SpanStop(ctx, "GetUserSubOrdinates")

	if err := validateGetUserSubOrdinatesRequest(req); err != nil {
		return nil, err
	}

	resp, err := s.Repository.ReadSubOrdinates(ctx, req.GetID())
	if err != nil {
		return nil, err
	}
	return &response.GetUserSubOrdinatesResponse{
		SubOrdinates: resp,
	}, nil
}

func (s *serviceMgmtUser) PatchUser(ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "PatchUser")
	defer middleware.SpanStop(ctx, "PatchUser")

	if err := validatePatchUserRequest(req); err != nil {
		return nil, err
	}

	user, err := s.Repository.ReadByID(ctx, req.GetUser().GetID())
	if err != nil {
		return nil, err
	}

	patchedUser := user.Clone()
	reqUser := req.GetUser()
	hasChanges := false

	if reqUser.GetName() != patchedUser.GetName() {
		patchedUser.SetName(reqUser.GetName())
		hasChanges = true
	}
	if reqUser.GetEmail() != patchedUser.GetEmail() {
		patchedUser.SetEmail(reqUser.GetEmail())
		hasChanges = true
	}
	if reqUser.GetSupervisor() != patchedUser.GetSupervisor() {
		patchedUser.SetSupervisor(reqUser.GetSupervisor())
		hasChanges = true
	}

	if !hasChanges { // no changes, return current user
		return &response.PatchUserResponse{
			User: user,
		}, nil
	}

	now := time.Now().UTC()
	patchedUser.GetMeta().VersionIncr()
	patchedUser.GetMeta().SetUpdated(now)

	resp, err := s.Repository.Write(ctx, patchedUser)
	if err != nil {
		return nil, err
	}
	return &response.PatchUserResponse{
		User: resp,
	}, nil
}

func (s *serviceMgmtUser) DeleteUser(ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error) {
	middleware.SpanStart(ctx, "DeleteUser")
	defer middleware.SpanStop(ctx, "DeleteUser")

	if err := validateDeleteUserRequest(req); err != nil {
		return nil, err
	}

	return &response.DeleteUserResponse{}, s.Repository.Delete(ctx, req.GetID())
}
