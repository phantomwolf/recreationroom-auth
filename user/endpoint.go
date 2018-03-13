package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type Response interface {
	Err() error
}

type createUserRequest struct {
	Name     string
	Password string
	Email    string
}

type createUserResponse struct {
	ID        null.Int64  `json:"id,string"`
	Name      null.String `json:"name"`
	Email     null.String `json:"email"`
	CreatedAt null.Time   `json:"created_at"`
	Error     error       `json:"error"`
}

func (res *createUserResponse) Err() error {
	return res.Error
}

func makeCreateUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createUserRequest)
		user, err := serv.Create(ctx, req.Name, req.Password, req.Email)
		if err != nil {
			log.Debugf("[user/endpoint.go:CreateUserEndpoint] User creation failed: %s\n", err.Error())
			return &createUserResponse{Error: err}, nil
		}
		res := &createUserResponse{
			ID:        null.Int64From(user.ID),
			Name:      null.StringFrom(user.Name),
			Email:     null.StringFrom(user.Email),
			CreatedAt: null.TimeFrom(user.CreatedAt),
		}
		return res, nil
	}
}

type updateUserRequest struct {
	ID       int64
	Name     string
	Password string
	Email    string
}

type updateUserResponse struct {
	Error error `json:"error"`
}

func (res *updateUserResponse) Err() error {
	return res.Error
}

func makeUpdateUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*updateUserRequest)
		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			log.Debugf("[user/endpoint.go:UpdateUserEndpoint] Failed to get user %d: %s\n", req.ID, err.Error())
			return &updateUserResponse{Error: err}, nil
		}

		if err := user.SetName(req.Name); err != nil {
			return &updateUserResponse{Error: err}, nil
		}
		if err := user.SetPassword(req.Password); err != nil {
			return &updateUserResponse{Error: err}, nil
		}
		if err := user.SetEmail(req.Email); err != nil {
			return &updateUserResponse{Error: err}, nil
		}

		if err := serv.Update(ctx, user); err != nil {
			log.Debugf("[user/endpoint.go:UpdateUserEndpoint] User updating failed: %s\n", err.Error())
			return &updateUserResponse{Error: err}, nil
		}
		return &updateUserResponse{}, nil
	}
}

type patchUserRequest map[string]interface{}

type patchUserResponse struct {
	Error error `json:"error"`
}

func (res *patchUserResponse) Err() error {
	return res.Error
}

func makePatchUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(patchUserRequest)
		if err := serv.Patch(ctx, req); err != nil {
			log.Debugf("[user/endpoint.go:PatchUserEndpoint] User patching failed: %s\n", err.Error())
			return &patchUserResponse{Error: err}, nil
		}
		return &patchUserResponse{}, nil
	}
}

type deleteUserRequest struct {
	ID int64
}

type deleteUserResponse struct {
	Error error `json:"error"`
}

func (res *deleteUserResponse) Err() error {
	return res.Error
}

func makeDeleteUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*deleteUserRequest)
		if err := serv.Delete(ctx, req.ID); err != nil {
			log.Debugf("[user/endpoint.go:DeleteUserEndpoint] User deletion failed: %s\n", err.Error())
			return &deleteUserResponse{Error: err}, nil
		}
		return &deleteUserResponse{}, nil
	}
}

type getUserRequest struct {
	ID int64
}

type getUserResponse struct {
	ID        null.Int64  `json:"id"`
	Name      null.String `json:"name"`
	Email     null.String `json:"email"`
	CreatedAt null.Time   `json:"created_at"`
	Error     error       `json:"error"`
}

func (res *getUserResponse) Err() error {
	return res.Error
}

func makeGetUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*getUserRequest)
		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			log.Debugf("[user/endpoint.go:GetUserEndpoint] User get failed: %s\n", err.Error())
			return &getUserResponse{Error: err}, nil
		}
		res := &getUserResponse{
			ID:        null.Int64From(user.ID),
			Name:      null.StringFrom(user.Name),
			Email:     null.StringFrom(user.Email),
			CreatedAt: null.TimeFrom(user.CreatedAt),
			Error:     nil,
		}
		return res, nil
	}
}
