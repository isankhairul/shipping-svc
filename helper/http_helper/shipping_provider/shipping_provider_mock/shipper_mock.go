package shipping_provider_mock

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"

	"github.com/stretchr/testify/mock"
)

type ShipperMock struct {
	Mock mock.Mock
}

func (h *ShipperMock) GetPricingDomestic(req *request.GetPricingDomestic) (*response.GetPricingDomestic, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.GetPricingDomestic), nil
}

func (h *ShipperMock) GetShippingRate(origin, destination *entity.CourierCoverageCode, data *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.ShippingRateCommonResponse), nil
}
