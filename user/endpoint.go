package user

import (
	"context"

	"github.com/PhantomWolf/recreationroom-auth/response"
	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type createUserRequest struct {
	Name     null.String
	Password null.String
	Email    null.String
}

func makeCreateUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createUserRequest)

		if req.Name.IsZero() || req.Password.IsZero() || req.Email.IsZero() {
			return response.New(statusError, codeInvalidRequest, nil, "name/password/email can't be empty"), nil
		}

		user, err := serv.Create(ctx, req.Name.String, req.Password.String, req.Email.String)
		if err != nil {
			log.Debugf("[user/endpoint.go:CreateUserEndpoint] User creation failed: %s\n", err.Error())
			return response.New(statusError, codeFailure, nil, err.Error()), nil
		}

		payload := map[string]interface{}{
			"user": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		}
		return response.New(statusOK, codeSuccess, payload, "User created successfully"), nil
	}
}

type updateUserRequest struct {
	ID       int64
	Name     null.String
	Password null.String
	Email    null.String
}

func makeUpdateUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*updateUserRequest)
		if req.Name.IsZero() || req.Password.IsZero() || req.Email.IsZero() {
			return response.New(statusError, codeInvalidRequest, nil, "name/password/email can't be empty"), nil
		}

		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			log.Debugf("[user/endpoint.go:UpdateUserEndpoint] Failed to get user %d: %s\n", req.ID, err.Error())
			return response.New(statusError, codeNotFound, nil, err.Error()), nil
		}
		if err := user.SetName(req.Name.String); err != nil {
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		if err := user.SetPassword(req.Password.String); err != nil {
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		if err := user.SetEmail(req.Email.String); err != nil {
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		if err := serv.Update(ctx, user); err != nil {
			log.Debugf("[user/endpoint.go:UpdateUserEndpoint] User updating failed: %s\n", err.Error())
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil, "User updated successfully"), nil
	}
}

type patchUserRequest map[string]interface{}

func makePatchUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(patchUserRequest)
		if err := serv.Patch(ctx, req); err != nil {
			log.Debugf("[user/endpoint.go:PatchUserEndpoint] User patching failed: %s\n", err.Error())
			return response.New(statusError, codeFailure, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil, "User modified successfully"), nil
	}
}

type deleteUserRequest struct {
	ID int64
}

func makeDeleteUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*deleteUserRequest)
		if err := serv.Delete(ctx, req.ID); err != nil {
			log.Debugf("[user/endpoint.go:DeleteUserEndpoint] User deletion failed: %s\n", err.Error())
			return response.New(statusError, codeFailure, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil, "User deleted successfully"), nil
	}
}

type getUserRequest struct {
	ID int64
}

func makeGetUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*getUserRequest)
		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			log.Debugf("[user/endpoint.go:GetUserEndpoint] User get failed: %s\n", err.Error())
			return response.New(statusError, codeNotFound, nil, err.Error()), nil
		}
		payload := map[string]interface{}{
			"user": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		}
		return response.New(statusOK, codeSuccess, payload, "User found"), nil
	}
}
