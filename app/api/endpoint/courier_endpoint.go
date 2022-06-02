package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type CourierEndpoint struct {
	Save                endpoint.Endpoint
	Show                endpoint.Endpoint
	List                endpoint.Endpoint
	Update              endpoint.Endpoint
	Delete              endpoint.Endpoint
	ListChannelCouriers endpoint.Endpoint

	//Courier-Serivce
	SaveCourierSerivce   endpoint.Endpoint
	ShowCourierSerivce   endpoint.Endpoint
	ListCourierSerivce   endpoint.Endpoint
	UpdateCourierSerivce endpoint.Endpoint
	DeleteCourierSerivce endpoint.Endpoint
}

func MakeCourierEndpoints(s service.CourierService, cc service.ChannelCourierService) CourierEndpoint {
	return CourierEndpoint{
		Save:                makeSaveCourier(s),
		Show:                makeShowCourier(s),
		List:                getListCouriers(s),
		Delete:              makeDeleteCourier(s),
		Update:              makeUpdateCourier(s),
		ListChannelCouriers: ListChannelCouriers(cc),

		SaveCourierSerivce:   makeSaveCourierService(s),
		ShowCourierSerivce:   makeShowCourierService(s),
		ListCourierSerivce:   getListCourierServices(s),
		DeleteCourierSerivce: makeDeleteCourierService(s),
		UpdateCourierSerivce: makeUpdateCourierService(s),
	}
}

func makeSaveCourier(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveCourierRequest)
		result, msg := s.CreateCourier(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeShowCourier(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetCourier(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func getListCouriers(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CourierListRequest)
		result, pagination, msg := s.GetList(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeDeleteCourier(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		msg := s.DeleteCourier(fmt.Sprint(request))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeUpdateCourier(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateCourierRequest)
		ret, msg := s.UpdateCourier(req.Uid, req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, ret, nil), nil
	}
}

func makeSaveCourierService(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveCourierServiceRequest)
		result, msg := s.CreateCourierService(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeShowCourierService(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetCourierService(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func getListCourierServices(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CourierServiceListRequest)
		result, pagination, msg := s.GetListCourierService(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeDeleteCourierService(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		msg := s.DeleteCourierService(fmt.Sprint(request))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeUpdateCourierService(s service.CourierService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.UpdateCourierServiceRequest)
		result, msg := s.UpdateCourierService(req.Uid, req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}
