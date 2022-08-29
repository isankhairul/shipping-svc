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

	"github.com/gorilla/schema"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

const (
	pathUID = "uid"
)

func ChannelHttpHandler(s service.ChannelService, ccs service.ChannelCourierService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeChannelEndpoints(s, ccs)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannel, global.PathChannelApp)).Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveChannel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannel, global.PathChannelApp)).Handler(httptransport.NewServer(
		ep.List,
		decodeListChannel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannel, global.PathChannelAppUID)).Handler(httptransport.NewServer(
		ep.Show,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannel, global.PathChannelAppUID)).Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateChannel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannel, global.PathChannelAppUID)).Handler(httptransport.NewServer(
		ep.Delete,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannel, global.PathChannelCourierStatus)).Handler(httptransport.NewServer(
		ep.ListStatus,
		decodeListChannelStatus,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixChannel, global.PathUIDCourierList)).Handler(httptransport.NewServer(
		ep.ChannelCourierList,
		decodeChannelCourierList,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeSaveChannel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	return req, nil
}

func decodeListChannel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ChannelListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.GetFilter()

	return params, nil
}

func decodeUpdateChannel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)
	req.Uid = mux.Vars(r)[pathUID]
	return req, nil
}

func decodeListChannelStatus(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetChannelCourierStatusRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.GetFilter()

	return params, nil
}

func decodeChannelCourierList(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetChannelCourierListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.ChannelUID = mux.Vars(r)[pathUID]
	params.SetFilterMap()

	return params, nil
}
