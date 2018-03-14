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
		user, err := serv.Create(ctx, req.Name, req.Password, req.Email)
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
		return response.New(statusOK, codeSuccess, payload), nil
	}
}

type updateUserRequest struct {
	ID       null.Int64
	Name     null.String
	Password null.String
	Email    null.String
}

func makeUpdateUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*updateUserRequest)
		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			log.Debugf("[user/endpoint.go:UpdateUserEndpoint] Failed to get user %d: %s\n", req.ID, err.Error())
			return response.New(statusError, codeNotFound, nil, err.Error()), nil
		}
		if err := user.SetName(req.Name); err != nil {
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		if err := user.SetPassword(req.Password); err != nil {
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		if err := user.SetEmail(req.Email); err != nil {
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		if err := serv.Update(ctx, user); err != nil {
			log.Debugf("[user/endpoint.go:UpdateUserEndpoint] User updating failed: %s\n", err.Error())
			return response.New(statusError, codeInvalidRequest, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil), nil
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
		return response.New(statusOK, codeSuccess, nil), nil
	}
}

type deleteUserRequest struct {
	ID null.Int64
}

func makeDeleteUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*deleteUserRequest)
		if err := serv.Delete(ctx, req.ID); err != nil {
			log.Debugf("[user/endpoint.go:DeleteUserEndpoint] User deletion failed: %s\n", err.Error())
			return response.New(statusError, codeFailure, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil), nil
	}
}

type getUserRequest struct {
	ID null.Int64
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
		return response.New(statusOK, codeSuccess, payload), nil
	}
}
