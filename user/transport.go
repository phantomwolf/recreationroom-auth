package user

import (
	"github.com/PhantomWolf/recreationroom-auth/response"
	_ "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null"

	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

var (
	errHTTPStatusMap = map[error]int{
		ErrInvalidRequest:           http.StatusBadRequest,
		ErrUserInvalidName:          http.StatusBadRequest,
		ErrUserInvalidPassword:      http.StatusBadRequest,
		ErrUserInvalidEmail:         http.StatusBadRequest,
		ErrUserWrongLoginOrPassword: http.StatusBadRequest,

		ErrUserNotFound: http.StatusNotFound,

		ErrUserAlreadyExists: http.StatusConflict,
	}
)

// MakeHandler returns a handler for the user service.
func MakeHandler(serv Service, r *mux.Router) *mux.Router {
	opts := []kithttp.ServerOption{}
	createUserHandler := kithttp.NewServer(
		makeCreateUserEndpoint(serv),
		decodeCreateUserRequest,
		encodeCreateUserResponse,
		opts...,
	)
	updateUserHandler := kithttp.NewServer(
		makeUpdateUserEndpoint(serv),
		decodeUpdateUserRequest,
		encodeUpdateUserResponse,
		opts...,
	)
	patchUserHandler := kithttp.NewServer(
		makePatchUserEndpoint(serv),
		decodePatchUserRequest,
		encodePatchUserResponse,
		opts...,
	)
	deleteUserHandler := kithttp.NewServer(
		makeDeleteUserEndpoint(serv),
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
		opts...,
	)
	getUserHandler := kithttp.NewServer(
		makeGetUserEndpoint(serv),
		decodeGetUserRequest,
		encodeGetUserResponse,
		opts...,
	)
	resetPasswordHandler := kithttp.NewServer(
		makeResetPasswordEndpoint(serv),
		decodeResetPasswordRequest,
		encodeResetPasswordResponse,
		opts...,
	)
	createPasswordHandler := kithttp.NewServer(
		makeCreatePasswordEndpoint(serv),
		decodeCreatePasswordRequest,
		encodeCreatePasswordResponse,
		opts...,
	)
	updatePasswordHandler := kithttp.NewServer(
		makeUpdatePasswordEndpoint(serv),
		decodeUpdatePasswordRequest,
		encodeUpdatePasswordResponse,
		opts...,
	)

	// POST /users
	r.Handle("/users", createUserHandler).Methods("POST")
	// PUT /users/{id}
	r.Handle("/users/{id:[0-9]+}", updateUserHandler).Methods("PUT")
	// PATCH /users/{id}
	r.Handle("/users/{id:[0-9]+}", patchUserHandler).Methods("PATCH")
	// DELETE /users/{id}
	r.Handle("/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")
	// GET /users/{id}
	r.Handle("/users/{id:[0-9]+}", getUserHandler).Methods("GET")
	// GET /password/reset
	r.Handle("/password/reset", resetPasswordHandler).Methods("GET")
	// POST /users/{id}/password
	r.Handle("/users/{id}/password", createPasswordHandler).Methods("POST")
	// PUT /users/{id}/password
	r.Handle("/users/{id}/password", updatePasswordHandler).Methods("PUT", "PATCH")
	return r
}

func errToStatus(err error) int {
	status, ok := errHTTPStatusMap[err]
	if ok {
		return status
	}
	return http.StatusBadRequest
}

func encodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Name     null.String `json:"name"`
		Password null.String `json:"password"`
		Email    null.String `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Debugf("[user/transport.go] Error decoding CreateUserRequest: %s\n", err.Error())
		return nil, nil
	}
	return &createUserRequest{
		Name:     body.Name,
		Password: body.Password,
		Email:    body.Email,
	}, nil
}

func encodeUpdateUserResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Debugf("[user/transport.go] Invalid id %s: %s\n", vars["id"], err.Error())
		return nil, ErrInvalidRequest
	}

	var body struct {
		Name     null.String `json:"name"`
		Password null.String `json:"password"`
		Email    null.String `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Debugf("[user/transport.go] Error decoding UpdateUserRequest: %s\n", err.Error())
		return nil, nil
	}
	return &updateUserRequest{
		ID:       id,
		Name:     body.Name,
		Password: body.Password,
		Email:    body.Email,
	}, nil
}

func encodePatchUserResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

func decodePatchUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		log.Debugf("[user/transport.go] Invalid id %s: %s\n", vars["id"], err.Error())
		return nil, ErrInvalidRequest
	}

	var body struct {
		Name     null.String `json:"name"`
		Password null.String `json:"password"`
		Email    null.String `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Debugf("[user/transport.go] Error decoding PatchUserRequest: %s\n", err.Error())
		return nil, nil
	}

	req := patchUserRequest{"id": id}
	if body.Name.Valid {
		req["name"] = body.Name.String
	}
	if body.Password.Valid {
		req["password"] = body.Password.String
	}
	if body.Email.Valid {
		req["email"] = body.Email.String
	}
	if len(req) < 2 {
		log.Debugf("[user/transport.go] Invalid request: %v\n", body)
		return nil, ErrInvalidRequest
	}
	return req, nil
}

func encodeDeleteUserResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

func decodeDeleteUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		log.Debugf("[user/transport.go] Invalid id %s: %s\n", vars["id"], err.Error())
		return nil, nil
	}
	return &deleteUserRequest{ID: id}, nil
}

func encodeGetUserResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Debugf("[user/transport.go] Invalid id %s: %s\n", vars["id"], err.Error())
		return nil, nil
	}
	return &getUserRequest{ID: id}, nil
}

func encodeResetPasswordResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

// GET /password/reset
func decodeResetPasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		NameOrEmail null.String `json:"name_or_email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Debugf("[user/transport.go] Error decoding ResetPasswordRequest: %s\n", err.Error())
		return nil, nil
	}
	return &resetPasswordRequest{NameOrEmail: body.NameOrEmail}, nil
}

func encodeCreatePasswordResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

func decodeCreatePasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return nil, nil
	}

	var body struct {
		Token       null.String `json:"token"`
		NewPassword null.String `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, nil
	}
	return &createPasswordRequest{ID: id, Token: body.Token, NewPassword: body.NewPassword}, nil
}

func encodeUpdatePasswordResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	r := res.(*response.Response)
	if r.Status == StatusOK {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(errToStatus(r.Err))
	}
	return json.NewEncoder(w).Encode(r)
}

// PUT /users/<id>/password
func decodeUpdatePasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return nil, nil
	}

	var body struct {
		Password    null.String `json:"password"`
		NewPassword null.String `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, nil
	}
	return &updatePasswordRequest{ID: id, Password: body.Password, NewPassword: body.NewPassword}, nil
}
