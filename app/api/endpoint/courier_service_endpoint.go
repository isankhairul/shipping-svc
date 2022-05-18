package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/endpoint"
)

type CourierServiceEndpoint struct {
	Save   endpoint.Endpoint
	Show   endpoint.Endpoint
	List   endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

func MakeCourierServiceEndpoints(s service.CourierServiceService) CourierServiceEndpoint {
	return CourierServiceEndpoint{
		Save:   makeSaveCourierService(s),
		Show:   makeShowCourierService(s),
		List:   getListCourierServices(s),
		Delete: makeDeleteCourierService(s),
		Update: makeUpdateCourierService(s),
	}
}

func makeSaveCourierService(s service.CourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveCourierServiceRequest)
		result, msg := s.CreateCourierService(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeShowCourierService(s service.CourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetCourierService(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func getListCourierServices(s service.CourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CourierServiceListRequest)
		result, pagination, msg := s.GetList(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeDeleteCourierService(s service.CourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		msg := s.DeleteCourierService(fmt.Sprint(request))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}

func makeUpdateCourierService(s service.CourierServiceService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveCourierServiceRequest)
		msg := s.UpdateCourierService(req.Uid, req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
