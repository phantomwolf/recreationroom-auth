package auth

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(bs service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	registerHandler := kithttp.NewServer(
		makeRegisterEndpoint(bs),
		decodeRegisterRequest,
		encodeResponse,
		opts...,
	)
}

type errorer interface {
    error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    if e, ok :=
}
