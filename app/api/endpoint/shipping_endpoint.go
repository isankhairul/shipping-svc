package endpoint

import (
	"context"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type ShippingEndpoint struct {
	GetShippingRate               endpoint.Endpoint
	GetShippingRateByShippingType endpoint.Endpoint
	CreateDelivery                endpoint.Endpoint
	GetOrderShippingTracking      endpoint.Endpoint
}

func MakeShippingEndpoint(s service.ShippingService) ShippingEndpoint {
	return ShippingEndpoint{
		GetShippingRate:               makeGetShippingRate(s),
		GetShippingRateByShippingType: makeGetShippingRateByShippingType(s),
		CreateDelivery:                makeCreateDelivery(s),
		GetOrderShippingTracking:      makeGetOrderShippingTracking(s),
	}
}

func makeGetShippingRate(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetShippingRateRequest)
		result, msg := s.GetShippingRate(req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetShippingRateByShippingType(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetShippingRateRequest)
		result, msg := s.GetShippingRateByShippingType(req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeCreateDelivery(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CreateDelivery)
		result, msg := s.CreateDelivery(&req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
func makeGetOrderShippingTracking(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetOrderShippingTracking)
		result, msg := s.OrderShippingTracking(&req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
