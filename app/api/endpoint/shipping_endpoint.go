package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/endpoint"
)

type ShippingEndpoint struct {
	GetShippingRate               endpoint.Endpoint
	GetShippingRateByShippingType endpoint.Endpoint
	CreateDelivery                endpoint.Endpoint
	GetOrderShippingTracking      endpoint.Endpoint
	UpdateStatusShipper           endpoint.Endpoint
	GetOrderShippingList          endpoint.Endpoint
	GetOrderShippingDetail        endpoint.Endpoint
	CancelPickUp                  endpoint.Endpoint
	CancelOrder                   endpoint.Endpoint
}

func MakeShippingEndpoint(s service.ShippingService) ShippingEndpoint {
	return ShippingEndpoint{
		GetShippingRate:               makeGetShippingRate(s),
		GetShippingRateByShippingType: makeGetShippingRateByShippingType(s),
		CreateDelivery:                makeCreateDelivery(s),
		GetOrderShippingTracking:      makeGetOrderShippingTracking(s),
		UpdateStatusShipper:           makeUpdateStatusShipper(s),
		GetOrderShippingList:          makeGetOrderShippingList(s),
		GetOrderShippingDetail:        makeGetOrderShippingDetail(s),
		CancelPickUp:                  makeCancelPickup(s),
		CancelOrder:                   makeCancelOrder(s),
	}
}

func makeGetShippingRate(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.GetShippingRateRequest)
		result, msg := s.GetShippingRate(req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetShippingRateByShippingType(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.GetShippingRateRequest)
		result, msg := s.GetShippingRateByShippingType(req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeCreateDelivery(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.CreateDelivery)
		req.JWTInfo = *jwtInfo
		result, msg := s.CreateDelivery(&req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
func makeGetOrderShippingTracking(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.GetOrderShippingTracking)
		result, msg := s.OrderShippingTracking(&req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
func makeUpdateStatusShipper(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		req := rqst.(request.WebhookUpdateStatusShipper)
		_, msg := s.UpdateStatusShipper(&req)
		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeGetOrderShippingList(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.GetOrderShippingList)
		result, pagination, msg := s.GetOrderShippingList(&req)
		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeGetOrderShippingDetail(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		result, msg := s.GetOrderShippingDetailByUID(fmt.Sprint(rqst))
		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeCancelPickup(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		msg = s.CancelPickup(fmt.Sprint(rqst))
		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeCancelOrder(s service.ShippingService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.CancelOrder)
		req.Body.JWTInfo = *jwtInfo
		msg = s.CancelOrder(&req)
		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
