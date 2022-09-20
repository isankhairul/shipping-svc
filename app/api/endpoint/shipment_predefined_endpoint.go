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

type ShipmentPredefinedEndpoint struct {
	List   endpoint.Endpoint
	Update endpoint.Endpoint
	Show   endpoint.Endpoint
}

func MakeShipmentPredefinedEndpoints(s service.ShipmentPredefinedService) ShipmentPredefinedEndpoint {
	return ShipmentPredefinedEndpoint{
		List:   listShipmentPredefined(s),
		Update: updateShipmentPredefined(s),
		Show:   makeShowShipmentPredefined(s),
	}
}

func listShipmentPredefined(s service.ShipmentPredefinedService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.ListShipmentPredefinedRequest)
		result, pagination, msg := s.GetAll(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, result, pagination), nil
	}
}

func updateShipmentPredefined(s service.ShipmentPredefinedService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		req := rqst.(request.UpdateShipmentPredefinedRequest)
		req.JWTInfo = *jwtInfo
		ret, msg := s.UpdateShipmentPredefined(req)
		if msg.Code == 4000 {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		return base.SetHttpResponse(msg.Code, msg.Message, ret, nil), nil
	}
}

func makeShowShipmentPredefined(s service.ShipmentPredefinedService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {

		// Retrieve JWT Info
		_, msg := global.SetJWTInfoFromContext(ctx)
		if msg.Code != message.SuccessMsg.Code {
			return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
		}

		ret, msg := s.GetByUID(rqst.(string))
		return base.SetHttpResponse(msg.Code, msg.Message, ret, nil), nil
	}
}
