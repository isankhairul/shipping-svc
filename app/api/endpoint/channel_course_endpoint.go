package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

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
		req := rqst.(request.SaveChannelCourierRequest)
		result, msg := s.CreateChannelCourier(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func GetChannelCourier(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetChannelCourier(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func ListChannelCouriers(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
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
		req := rqst.(request.UpdateChannelCourierRequest)
		result, msg := s.UpdateChannelCourier(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func DeleteChannelCourier(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		uid := fmt.Sprint(rqst)
		msg := s.DeleteChannelCourier(uid)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
