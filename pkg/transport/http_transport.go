package transport

import (
	"context"
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
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

	pr.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	pr.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "doc"}
	sh1 := middleware.Redoc(opts1, nil)
	pr.Handle("/doc", sh1)

	r := pr.PathPrefix("/prescription").Subrouter()

	r.Methods("POST").Path("/product").Handler(httptransport.NewServer(
		e.Save,
		decodeSave,
		encodeResponseHTTP,
		options...,
	))

	r.Methods("GET").Path("/products").Handler(httptransport.NewServer(
		e.Show,
		decodeShow,
		encodeResponseHTTP,
		options...,
	))

	r.Methods("GET").Path("/product/{id}").Handler(httptransport.NewServer(
		e.GetProduct,
		decodeGetProduct,
		encodeResponseHTTP,
		options...,
	))

	r.Methods("PUT").Path("/product/{id}").Handler(httptransport.NewServer(
		e.UpdateProduct,
		decodeUpdateProduct,
		encodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeSave(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req entity.JSONRequestProduct
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeShow(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeGetProduct(_ context.Context, r *http.Request) (request interface{}, err error) {
	params := mux.Vars(r)["id"]
	return params, nil
}

func decodeUpdateProduct(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req entity.JSONRequestUpdateProduct
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	req.Id = mux.Vars(r)["id"]
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
