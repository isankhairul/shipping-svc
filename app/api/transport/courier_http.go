package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"net/http"

	"github.com/gorilla/schema"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func CourierHttpHandler(s service.CourierService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeCourierEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path("/courier/courier").Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/courier/courier").Handler(httptransport.NewServer(
		ep.List,
		decodeListCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/courier/courier/{id}").Handler(httptransport.NewServer(
		ep.Show,
		decodeShowCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path("/courier/courier/{id}").Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path("/courier/courier/{id}").Handler(httptransport.NewServer(
		ep.Delete,
		decodeDeleteCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeSaveCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	global.HtmlEscape(&req)

	return req, nil
}

func decodeShowCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}

func decodeListCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CourierListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeUpdateCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	global.HtmlEscape(&req)

	req.Uid = mux.Vars(r)["id"]
	return req, nil
}

func decodeDeleteCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}
