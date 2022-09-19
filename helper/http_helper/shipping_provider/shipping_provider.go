package shipping_provider

import "strings"

const (
	ShipperCode      = "shipper"
	InternalCourier  = "internal"
	MerchantCourier  = "merchant"
	ThirPartyCourier = "third_party"

	StatusCreated       = "created"
	StatusRequestPickup = "request_pickup"
	StatusCancelled     = "cancelled"
)

var shipperCancelableStatus = []string{
	StatusRequestPickup,
}

func IsOrderCancelable(courierCode, status string) bool {
	var statusList []string

	switch courierCode {
	case ShipperCode:
		statusList = shipperCancelableStatus
	}

	for _, v := range statusList {
		if strings.EqualFold(v, status) {
			return true
		}
	}

	return false
}
