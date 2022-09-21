package endpoint

import (
	"context"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/endpoint"
)

var (
	emptyArrayResult = []string{}
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

		// Retrieve JWT Info
		jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.SaveChannelCourierServiceRequest)
		req.JWTInfo = *jwtInfo
		result, msg := s.CreateChannelCourierService(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func ListChannelCourierServices(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.ChannelCourierServiceListRequest)
		result, paging, msg := s.ListChannelCouriersService(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, paging), nil
	}
}

func GetChannelCourierService(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		result, msg := s.GetChannelCourierService(rqst.(string))

		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, emptyArrayResult, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func UpdateChannelCourierService(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.UpdateChannelCourierServiceRequest)
		req.Body.JWTInfo = *jwtInfo
		result, msg := s.UpdateChannelCourierService(req)

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func DeleteChannelCourierService(s service.ChannelCourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		msg = s.DeleteChannelCourierService(rqst.(string))

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
