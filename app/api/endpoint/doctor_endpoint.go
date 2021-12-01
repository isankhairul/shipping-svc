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

type DoctorEndpoint struct {
	SaveDoctor endpoint.Endpoint
	Show       endpoint.Endpoint
}

func MakeDoctorEndpoints(s service.DoctorService) DoctorEndpoint {
	return DoctorEndpoint{
		SaveDoctor: makeSaveDoctor(s),
		Show:       makeShowDoctor(s),
	}
}

func makeSaveDoctor(s service.DoctorService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveDoctorRequest)
		result, code, msg := s.CreateDoctor(req)
		if msg != "" {
			return base.SetHttpResponse(code, msg, nil, nil), nil
		}

		return base.SetHttpResponse(code, message.MSG_SUCCESS, result, nil), nil
	}
}

func makeShowDoctor(s service.DoctorService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, code, msg := s.GetDoctor(fmt.Sprint(rqst))
		if msg != "" {
			return base.SetHttpResponse(code, msg, nil, nil), nil
		}

		return base.SetHttpResponse(code, message.MSG_SUCCESS, result, nil), nil
	}
}
