package auth

import (
	"context"
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func decodeRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeRegisterResponse(ctx context.Context, w http.ResponseWriter, res registerResponse) error {
	switch res.Err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	return json.NewEncoder(w).Encode(res)
}

func decodeUnregisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req unregisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeUnregisterResponse(ctx context.Context, w http.ResponseWriter, res unregisterResponse) error {
	switch res.Err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, res loginResponse) error {
	switch res.Err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}

func decodeLogoutRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeLogoutResponse(ctx context.Context, w http.ResponseWriter, res logoutResponse) error {
	switch res.Err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}

func MakeHandler(bs service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	registerHandler := kithttp.NewServer(
		makeRegisterEndpoint(bs),
		decodeRegisterRequest,
		encodeRegisterResponse,
		opts...,
	)

	unregisterHandler := kithttp.NewServer(
		makeUnregisterEndpoint(bs),
		decodeUnregisterRequest,
		encodeUnregisterResponse,
		opts...,
	)

	loginHandler := kithttp.NewServer(
		makeLoginEndpoint(bs),
		decodeLoginRequest,
		encodeLoginResponse,
		opts...,
	)

	logoutHandler := kithttp.NewServer(
		makeLogoutEndpoint(bs),
		decodeLogoutRequest,
		encodeLogoutResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/auth/v1/register", registerHandler).Methods("POST")
	r.Handle("/auth/v1/unregister", unregisterHandler).Methods("POST")
	r.Handle("/auth/v1/login", loginHandler).Methods("POST")
	r.Handle("/auth/v1/logout", logoutHandler).Methods("POST")

	return r
}
