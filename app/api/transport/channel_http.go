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

func ChannelHttpHandler(s service.ChannelService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeChannelEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path("/channel/channel-app").Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveChannel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/channel/channel-app").Handler(httptransport.NewServer(
		ep.List,
		decodeListChannel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/channel/channel-app/{id}").Handler(httptransport.NewServer(
		ep.Show,
		decodeShowChannel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path("/channel/channel-app/{id}").Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateChannel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path("/channel/channel-app/{id}").Handler(httptransport.NewServer(
		ep.Delete,
		decodeDeleteChannel,
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
	global.HtmlEscape(&req)

	return req, nil
}

func decodeShowChannel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}

func decodeListChannel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ChannelListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeUpdateChannel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	global.HtmlEscape(&req)
	req.Uid = mux.Vars(r)["id"]
	return req, nil
}

func decodeDeleteChannel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}
