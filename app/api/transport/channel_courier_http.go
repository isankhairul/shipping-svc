package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func ChannelCourierHttpHandler(s service.ChannelCourierService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeChannelCourierEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourier, global.PathUID)).Handler(httptransport.NewServer(
		ep.GetChannelCourier,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourier)).Handler(httptransport.NewServer(
		ep.SaveChannelCourier,
		decodeSaveChannelCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourier, global.PathUID)).Handler(httptransport.NewServer(
		ep.UpdateChannelCourier,
		decodeUpdateChannelCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourier)).Handler(httptransport.NewServer(
		ep.ListChannelCouriers,
		decodePaginationRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourier, global.PathUID)).Handler(httptransport.NewServer(
		ep.DeleteChannelCourier,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodePaginationRequestHTTP(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ChannelCourierListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.GetFilter()
	return params, nil
}

func decodeSaveChannelCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveChannelCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//global.HtmlEscape(&req)

	return req, nil
}

func decodeUpdateChannelCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateChannelCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)
	req.Uid = mux.Vars(r)[pathUID]
	return req, nil
}
