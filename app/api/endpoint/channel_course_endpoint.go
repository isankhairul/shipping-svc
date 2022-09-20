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

type ChannelCourierEndpoint struct {
	SaveChannelCourier   endpoint.Endpoint
	ListChannelCouriers  endpoint.Endpoint
	GetChannelCourier    endpoint.Endpoint
	UpdateChannelCourier endpoint.Endpoint
	DeleteChannelCourier endpoint.Endpoint
}

func MakeChannelCourierEndpoints(s service.ChannelCourierService) ChannelCourierEndpoint {
	return ChannelCourierEndpoint{
		SaveChannelCourier:   SaveChannelCourier(s),
		ListChannelCouriers:  ListChannelCouriers(s),
		GetChannelCourier:    GetChannelCourier(s),
		UpdateChannelCourier: UpdateChannelCourier(s),
		DeleteChannelCourier: DeleteChannelCourier(s),
	}
}

func SaveChannelCourier(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.SaveChannelCourierRequest)
		req.JWTInfo = *jwtInfo
		result, msg := s.CreateChannelCourier(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func GetChannelCourier(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		result, msg := s.GetChannelCourier(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func ListChannelCouriers(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		input := rqst.(request.ChannelCourierListRequest)
		result, pagination, msg := s.ListChannelCouriers(input)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func UpdateChannelCourier(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.UpdateChannelCourierRequest)
		req.JWTInfo = *jwtInfo
		result, msg := s.UpdateChannelCourier(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func DeleteChannelCourier(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		uid := fmt.Sprint(rqst)
		msg = s.DeleteChannelCourier(uid)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
