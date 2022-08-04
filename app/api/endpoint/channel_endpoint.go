package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type ChannelEndpoint struct {
	Save               endpoint.Endpoint
	Show               endpoint.Endpoint
	List               endpoint.Endpoint
	Update             endpoint.Endpoint
	Delete             endpoint.Endpoint
	ListStatus         endpoint.Endpoint
	ChannelCourierList endpoint.Endpoint
}

func MakeChannelEndpoints(s service.ChannelService, ccs service.ChannelCourierService) ChannelEndpoint {
	return ChannelEndpoint{
		Save:               makeSaveChannel(s),
		Show:               makeShowChannel(s),
		List:               getListChannels(s),
		Delete:             makeDeleteChannel(s),
		Update:             makeUpdateChannel(s),
		ListStatus:         makeGetListChannelStatus(s),
		ChannelCourierList: makeGetChannelCourierList(ccs),
	}
}

func makeSaveChannel(s service.ChannelService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveChannelRequest)
		result, msg := s.CreateChannel(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeShowChannel(s service.ChannelService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetChannel(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func getListChannels(s service.ChannelService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.ChannelListRequest)
		result, pagination, msg := s.GetList(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeDeleteChannel(s service.ChannelService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		msg := s.DeleteChannel(fmt.Sprint(request))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeUpdateChannel(s service.ChannelService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateChannelRequest)
		msg := s.UpdateChannel(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeGetListChannelStatus(s service.ChannelService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetChannelCourierStatusRequest)
		result, pagination, msg := s.GetListStatus(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeGetChannelCourierList(s service.ChannelCourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.GetChannelCourierListRequest)
		result, pagination, msg := s.GetChannelCourierListByChannelUID(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}
