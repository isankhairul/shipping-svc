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

const (
	shippingTypePath = "shipping-type"
	topicName        = "topic-name"
	channelUID       = "channel-uid"
)

func ShippingHttpHandler(s service.ShippingService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeShippingEndpoint(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathShippingRate)).Handler(httptransport.NewServer(
		ep.GetShippingRate,
		decodeGetShippingRate,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathShippingRateShippingType)).Handler(httptransport.NewServer(
		ep.GetShippingRateByShippingType,
		decodeGetShippingRate,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShipping)).Handler(httptransport.NewServer(
		ep.CreateDelivery,
		decodeCreateDelivery,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderTracking)).Handler(httptransport.NewServer(
		ep.GetOrderShippingTracking,
		decodeGetOrderTracking,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShipping)).Handler(httptransport.NewServer(
		ep.GetOrderShippingList,
		decodeGetOrderShippingList,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShippingUID)).Handler(httptransport.NewServer(
		ep.GetOrderShippingDetail,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathCancelPickupUID)).Handler(httptransport.NewServer(
		ep.CancelPickUp,
		decodeCancelPickup,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathCancelOrderUID)).Handler(httptransport.NewServer(
		ep.CancelOrder,
		decodeCancelOrder,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathUpdateOrderTopicName)).Handler(httptransport.NewServer(
		ep.UpdateOrderShipping,
		decodeUpdateOrderShipping,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShippingLabel)).Handler(httptransport.NewServer(
		ep.GetOrderShippingLabel,
		decodeOrderShippingLabel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathRepickup)).Handler(httptransport.NewServer(
		ep.RepickupOrder,
		decodeRepickupOrder,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeGetShippingRate(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetShippingRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.ShippingType = mux.Vars(r)[shippingTypePath]
	return req, nil
}

func decodeCreateDelivery(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.CreateDelivery
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetOrderTracking(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetOrderShippingTracking
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.UID = mux.Vars(r)[pathUID]
	return params, nil
}

func decodeGetOrderShippingList(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetOrderShippingList
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.GetFilter()
	return params, nil
}
func decodeCancelOrder(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CancelOrder
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params.Body); err != nil {
		return nil, err
	}

	params.UID = mux.Vars(r)[pathUID]
	return params, nil
}

func decodeCancelPickup(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CancelPickup
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params.Body); err != nil {
		return nil, err
	}

	params.UID = mux.Vars(r)[pathUID]
	return params, nil
}

func decodeUpdateOrderShipping(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.UpdateOrderShipping
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params.Body); err != nil {
		return nil, err
	}

	params.TopicName = mux.Vars(r)[topicName]
	return params, nil
}

func decodeOrderShippingLabel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetOrderShippingLabel
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params.Body); err != nil {
		return nil, err
	}

	params.ChannelUID = mux.Vars(r)[channelUID]
	return params, nil
}

func decodeRepickupOrder(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.RepickupOrderRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, err
	}
	return params, nil
}
