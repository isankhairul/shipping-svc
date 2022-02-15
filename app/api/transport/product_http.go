package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"net/http"

	"github.com/gorilla/schema"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func ProductHttpHandler(s service.ProductService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeProductEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path("/products/").Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/products/list").Handler(httptransport.NewServer(
		ep.List,
		decodeListProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/products/{id}").Handler(httptransport.NewServer(
		ep.Show,
		decodeShowProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path("/products/{id}").Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path("/products/{id}").Handler(httptransport.NewServer(
		ep.Update,
		decodeDeleteProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeSaveProduct(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeShowProduct(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}

func decodeListProduct(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ProductListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeUpdateProduct(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Uid = mux.Vars(r)["id"]
	return req, nil
}

func decodeDeleteProduct(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}
