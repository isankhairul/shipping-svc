package shipping_provider_mock

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/helper/message"

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

func (h *ShipperMock) GetShippingRate(courierID *uint64, input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error) {
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

func (h *ShipperMock) CreateOrder(req *request.CreateOrderShipper) (*response.CreateOrderShipperResponse, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.CreateOrderShipperResponse), nil
}

func (h *ShipperMock) GetTimeslot(req *request.GetPickUpTimeslot) (*response.GetPickUpTimeslotResponse, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.GetPickUpTimeslotResponse), nil
}

func (h *ShipperMock) CreatePickUpOrder(req *request.CreatePickUpOrderShipper) (*response.CreatePickUpOrderShipperResponse, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.CreatePickUpOrderShipperResponse), nil
}

func (h *ShipperMock) CreatePickUpOrderWithTimeSlots(orderID ...string) (*response.CreatePickUpOrderShipperResponse, message.Message) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(message.Message)
		}
	}

	return arguments.Get(0).(*response.CreatePickUpOrderShipperResponse), message.SuccessMsg
}

func (h *ShipperMock) CreateDelivery(shipperOrderID string, courierService *entity.CourierService, req *request.CreateDelivery) (*response.CreateDeliveryThirdPartyData, message.Message) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return arguments.Get(0).(*response.CreateDeliveryThirdPartyData), arguments.Get(1).(message.Message)
		}
	}

	return arguments.Get(0).(*response.CreateDeliveryThirdPartyData), message.SuccessMsg
}

func (h *ShipperMock) GetOrderDetail(orderID string) (*response.GetOrderDetailResponse, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.GetOrderDetailResponse), nil
}

func (h *ShipperMock) GetTracking(orderID string) ([]response.GetOrderShippingTracking, message.Message) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(message.Message)
		}
	}

	return arguments.Get(0).([]response.GetOrderShippingTracking), message.SuccessMsg
}

func (h *ShipperMock) CancelPickupRequest(pickupCode string) (*response.MetadataResponse, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.MetadataResponse), nil
}

func (h *ShipperMock) CancelOrder(orderID string, req *request.CancelOrder) (*response.MetadataResponse, error) {
	arguments := h.Mock.Called()

	if len(arguments) > 1 {
		if arguments.Get(1) != nil {
			return nil, arguments.Get(1).(error)
		}
	}

	if arguments.Get(0) == nil {
		return nil, nil
	}

	return arguments.Get(0).(*response.MetadataResponse), nil
}
