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
	Save   endpoint.Endpoint
	Show   endpoint.Endpoint
	List   endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

func MakeChannelEndpoints(s service.ChannelService) ChannelEndpoint {
	return ChannelEndpoint{
		Save:   makeSaveChannel(s),
		Show:   makeShowChannel(s),
		List:   getListChannels(s),
		Delete: makeDeleteChannel(s),
		Update: makeUpdateChannel(s),
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
		req := rqst.(request.SaveChannelRequest)
		msg := s.UpdateChannel(req.Uid, req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
