package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"gokit_example/pkg/endpoint"
	"gokit_example/pkg/entity"
	"gokit_example/pkg/service"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(s service.Service, logger log.Logger) http.Handler {
	pr := mux.NewRouter()
	e := endpoint.MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r := pr.PathPrefix("/prescription").Subrouter()

	// device-check
	r.Methods("POST").Path("/product").Handler(httptransport.NewServer(
		e.Check,
		decodeCheck,
		encodeResponseHTTP,
		options...,
	))

	return pr
}

//Stock In Item Detaill
func decodeCheck(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req entity.JSONRequestProduct
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponseHTTP(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case service.ErrNotFound:
		return http.StatusNotFound
	case service.ErrAlreadyExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
