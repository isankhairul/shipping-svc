package endpoint

import (
	"context"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/endpoint"
)

type CourierCoverageCodeEndpoint struct {
	Save   endpoint.Endpoint
	Show   endpoint.Endpoint
	List   endpoint.Endpoint
	Update endpoint.Endpoint
	Import endpoint.Endpoint
	Delete endpoint.Endpoint
}

func MakeCourierCoverageCodeEndpoints(s service.CourierCoverageCodeService) CourierCoverageCodeEndpoint {
	return CourierCoverageCodeEndpoint{
		Save:   makeSaveCourierCoverageCode(s),
		Show:   makeShowCourierCoverageCode(s),
		List:   makeGetCourierCoverageCodes(s),
		Update: makeUpdateCourierCoverageCodes(s),
		Import: makeImportCourierCoverageCode(s),
		Delete: makeDeleteCourierCoverageCode(s),
	}
}

func makeSaveCourierCoverageCode(s service.CourierCoverageCodeService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveCourierCoverageCodeRequest)
		result, msg := s.CreateCourierCoverageCode(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeShowCourierCoverageCode(s service.CourierCoverageCodeService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, msg := s.GetCourierCoverageCode(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeGetCourierCoverageCodes(s service.CourierCoverageCodeService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.CourierCoverageCodeListRequest)
		result, pagination, msg := s.GetList(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func makeUpdateCourierCoverageCodes(s service.CourierCoverageCodeService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveCourierCoverageCodeRequest)
		result, msg := s.UpdateCourierCoverageCode(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
	}
}

func makeImportCourierCoverageCode(s service.CourierCoverageCodeService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (response interface{}, err error) {
		req := rqst.(request.ImportCourierCoverageCodeRequest)
		result, msg := s.ImportCourierCoverageCode(req)

		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return *result, nil
	}
}

func makeDeleteCourierCoverageCode(s service.CourierCoverageCodeService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		msg := s.DeleteCourierCoverageCode(fmt.Sprint(rqst))
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
	}
}
