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

func ChannelCourierServiceHttpHandler(s service.ChannelCourierServiceService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeChannelCourierServiceEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourierService, global.PathUID)).Handler(httptransport.NewServer(
		ep.GetChannelCourierService,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourierService)).Handler(httptransport.NewServer(
		ep.SaveChannelCourierService,
		decodeSaveChannelCourierServices,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourierService, global.PathUID)).Handler(httptransport.NewServer(
		ep.UpdateChannelCourierService,
		decodeUpdateChannelCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourierService)).Handler(httptransport.NewServer(
		ep.ListChannelCourierServices,
		decodeListChannelCourierServicesRequest,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourierService, global.PathUID)).Handler(httptransport.NewServer(
		ep.DeleteChannelCourierService,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeListChannelCourierServicesRequest(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ChannelCourierServiceListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.GetFilter()
	return params, nil
}

func decodeSaveChannelCourierServices(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveChannelCourierServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//global.HtmlEscape(&req)

	return req, nil
}

func decodeUpdateChannelCourierService(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateChannelCourierServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Body); err != nil {
		return nil, err
	}

	req.UID = mux.Vars(r)[pathUID]
	return req, nil
}
