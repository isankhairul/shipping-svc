package endpoint

import (
	"context"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type ChannelCourierServiceEndpoint struct {
	SaveChannelCourierService   endpoint.Endpoint
	ListChannelCourierServices  endpoint.Endpoint
	GetChannelCourierService    endpoint.Endpoint
	UpdateChannelCourierService endpoint.Endpoint
	DeleteChannelCourierService endpoint.Endpoint
}

func MakeChannelCourierServiceEndpoints(s service.ChannelCourierServiceService) ChannelCourierServiceEndpoint {
	return ChannelCourierServiceEndpoint{
		SaveChannelCourierService:   SaveChannelCourierService(s),
		ListChannelCourierServices:  ListChannelCourierServices(s),
		GetChannelCourierService:    GetChannelCourierService(s),
		UpdateChannelCourierService: UpdateChannelCourierService(s),
		DeleteChannelCourierService: DeleteChannelCourierService(s),
	}
}

func SaveChannelCourierService(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveChannelCourierServiceRequest)
		result, msg := s.CreateChannelCourierService(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func ListChannelCourierServices(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.ChannelCourierServiceListRequest)
		result, paging, msg := s.ListChannelCouriersService(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, paging), nil
	}
}

func GetChannelCourierService(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetChannelCourierService(rqst.(string))

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func UpdateChannelCourierService(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateChannelCourierServiceRequest)
		result, msg := s.UpdateChannelCourierService(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func DeleteChannelCourierService(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		msg := s.DeleteChannelCourierService(rqst.(string))

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
