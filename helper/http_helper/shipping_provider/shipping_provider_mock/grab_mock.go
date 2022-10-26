package shipping_provider_mock

import (
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"

	"github.com/stretchr/testify/mock"
)

type GrabMock struct {
	Mock mock.Mock
}

func (h *GrabMock) GetShippingRate(req *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error) {
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
