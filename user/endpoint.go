package user

import (
	"context"
	"time"

	"github.com/PhantomWolf/recreationroom-auth/response"
	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type userPayload struct {
	ID        int64     `json:"id,string"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type createUserRequest struct {
	Name     null.String
	Password null.String
	Email    null.String
}

func makeCreateUserEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createUserRequest)
		res := response.New()
		if req.Name.IsZero() || req.Password.IsZero() || req.Email.IsZero() {
			res.SetStatus(statusError, codeInvalidRequest)
			res.AddError(ErrInvalidRequest)
			return res, nil
		}

		user, err := serv.Create(ctx, req.Name.String, req.Password.String, req.Email.String)
		if err != nil {
			res.SetStatus(statusError, codeUserCreateFailure)
			res.AddError(err)
			return res, nil
		}
		res.SetResult("user", userPayload{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
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
		res := response.New()
		if req.Name.IsZero() || req.Password.IsZero() || req.Email.IsZero() {
			res.SetStatus(statusError, codeInvalidRequest)
			res.AddError(ErrInvalidRequest)
			return res, nil
		}

		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			res.SetStatus(statusError, codeUserGetFailure)
			res.AddError(err)
			return res, nil
		}
		if err := user.SetName(req.Name.String); err != nil {
			res.SetStatus(statusError, code)
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
		if req.ID <= 0 {
			return response.New(statusError, codeInvalidRequest, nil, "Invalid user id"), nil
		}
		user, err := serv.Get(ctx, req.ID)
		if err != nil {
			log.Debugf("[user/endpoint.go:GetUserEndpoint] User get failed: %s\n", err.Error())
			return response.New(statusError, codeNotFound, nil, err.Error()), nil
		}
		payload := map[string]interface{}{
			"user": userPayload{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			},
		}
		return response.New(statusOK, codeSuccess, payload, "User found"), nil
	}
}

type resetPasswordRequest struct {
	NameOrEmail null.String
}

func makeResetPasswordEndpoint(serv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*resetPasswordRequest)
		if req.NameOrEmail.IsZero() {
			return response.New(statusError, codeInvalidRequest, nil, "Can't be empty: name_or_email"), nil
		}
		if err := serv.ResetPassword(ctx, req.NameOrEmail.String); err != nil {
			return response.New(statusError, codeFailure, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil, "Password reset link sent"), nil
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
			return response.New(statusError, codeInvalidRequest, nil, "Invalid request"), nil
		}
		if err := serv.CreatePassword(ctx, req.ID, req.Token.String, req.NewPassword.String); err != nil {
			return response.New(statusError, codeFailure, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil, "Password created successfully"), nil
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
			return response.New(statusError, codeInvalidRequest, nil, "Missing password or new_password"), nil
		}
		if err := serv.UpdatePassword(ctx, req.ID, req.Password.String, req.NewPassword.String); err != nil {
			return response.New(statusError, codeFailure, nil, err.Error()), nil
		}
		return response.New(statusOK, codeSuccess, nil, "Password updated successfully"), nil
	}
}
