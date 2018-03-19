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
			return response.New(StatusError, CodeInvalidRequest, ErrInvalidRequest), nil
		}
		user, err := serv.Create(ctx, req.Name.String, req.Password.String, req.Email.String)
		if err != nil {
			return response.New(StatusError, CodeUserCreateFailure, err), nil
		}
		res := response.New(StatusOK, CodeSuccess, nil)
		res.SetResult("user", user)
		return res, nil
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
			return response.New(StatusError, CodeInvalidRequest, ErrInvalidRequest), nil
		}
		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			return response.New(StatusError, CodeUserGetFailure, err), nil
		}
		if err := user.SetName(req.Name.String); err != nil {
			return response.New(StatusError, CodeUserUpdateFailure, err), nil
		}
		if err := user.SetPassword(req.Password.String); err != nil {
			return response.New(StatusError, CodeUserUpdateFailure, err), nil
		}
		if err := user.SetEmail(req.Email.String); err != nil {
			return response.New(StatusError, CodeUserUpdateFailure, err), nil
		}
		if err := serv.Update(ctx, user); err != nil {
			return response.New(StatusError, CodeUserUpdateFailure, err), nil
		}
		return response.New(StatusOK, CodeSuccess, nil, "User updated successfully"), nil
	}
}

type patchUserRequest map[string]interface{}

func makePatchUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(patchUserRequest)
		if err := serv.Patch(ctx, req); err != nil {
			log.Debugf("[user/endpoint.go:PatchUserEndpoint] User patching failed: %s\n", err.Error())
			return response.New(StatusError, CodeUserUpdateFailure, err), nil
		}
		return response.New(StatusOK, CodeSuccess, nil, "User updated successfully"), nil
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
			return response.New(StatusError, CodeUserDeleteFailure, err), nil
		}
		return response.New(StatusOK, CodeSuccess, nil, "User deleted successfully"), nil
	}
}

type getUserRequest struct {
	ID int64
}

func makeGetUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*getUserRequest)
		if req.ID <= 0 {
			return response.New(StatusError, CodeInvalidRequest, ErrInvalidRequest), nil
		}
		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			log.Debugf("[user/endpoint.go:GetUserEndpoint] User get failed: %s\n", err.Error())
			return response.New(StatusError, CodeUserGetFailure, err), nil
		}
		res := response.New(StatusOK, CodeSuccess, nil)
		res.SetResult("user", user)
		return res, nil
	}
}

type resetPasswordRequest struct {
	NameOrEmail null.String
}

func makeResetPasswordEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*resetPasswordRequest)
		if req.NameOrEmail.IsZero() {
			return response.New(StatusError, CodeInvalidRequest, ErrInvalidRequest), nil
		}
		if err := serv.ResetPassword(ctx, req.NameOrEmail.String); err != nil {
			return response.New(StatusError, CodePasswordResetFailure, err), nil
		}
		return response.New(StatusOK, CodeSuccess, nil, "Password reset link sent"), nil
	}
}

type createPasswordRequest struct {
	ID          int64
	Token       null.String
	NewPassword null.String
}

func makeCreatePasswordEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createPasswordRequest)
		if req.Token.IsZero() || req.NewPassword.IsZero() {
			return response.New(StatusError, CodeInvalidRequest, ErrInvalidRequest), nil
		}
		if err := serv.CreatePassword(ctx, req.ID, req.Token.String, req.NewPassword.String); err != nil {
			return response.New(StatusError, CodePasswordCreateFailure, err), nil
		}
		return response.New(StatusOK, CodeSuccess, nil, "Password created successfully"), nil
	}
}

type updatePasswordRequest struct {
	ID          int64
	Password    null.String
	NewPassword null.String
}

func makeUpdatePasswordEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*updatePasswordRequest)
		if req.Password.IsZero() || req.NewPassword.IsZero() {
			return response.New(StatusError, CodeInvalidRequest, ErrInvalidRequest), nil
		}
		if err := serv.UpdatePassword(ctx, req.ID, req.Password.String, req.NewPassword.String); err != nil {
			return response.New(StatusError, CodePasswordUpdateFailure, err), nil
		}
		return response.New(StatusOK, CodeSuccess, nil, "Password updated successfully"), nil
	}
}
