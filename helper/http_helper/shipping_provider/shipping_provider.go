package shipping_provider

import (
	"fmt"
	"go-klikdokter/app/model/request"
	"go-klikdokter/pkg/util"
	"strings"

	"github.com/spf13/viper"
)

const (
	ShipperCode       = "shipper"
	GrabCode          = "grab"
	InternalCourier   = "internal"
	MerchantCourier   = "merchant"
	ThirPartyCourier  = "third_party"
	AggregatorCourier = "aggregator"

	StatusCreated       = "created"
	StatusRequestPickup = "request_pickup"
	StatusCancelled     = "cancelled"
)

var shipperPickupOrderCancelableStatus = []string{
	StatusRequestPickup,
}

var shipperOrderCancelableStatus = []string{
	StatusCreated,
	StatusRequestPickup,
}

func IsPickUpOrderCancelable(courierCode, status string) bool {
	var statusList []string

	switch courierCode {
	case ShipperCode:
		statusList = shipperPickupOrderCancelableStatus
	}

	for _, v := range statusList {
		if strings.EqualFold(v, status) {
			return true
		}
	}

	return false
}

func IsOrderCancelable(courierCode, status string) bool {
	var statusList []string

	switch courierCode {
	case ShipperCode:
		statusList = shipperOrderCancelableStatus
	}

	for _, v := range statusList {
		if strings.EqualFold(v, status) {
			return true
		}
	}

	return false
}

func ShipperWebhookAuth() string {
	// <api_key> + <endpoint_url> + <response_format>
	apiKey := viper.GetString("shipper.auth.value")
	endpoint := viper.GetString("shipper.webhook.update-status-endpoint")
	format := "json"

	return util.MD5Hash(apiKey + endpoint + format)
}

func GrabWebhookAuth(req *request.WebhookUpdateStatusGrabHeader) bool {
	//clientID := viper.GetString("grab.auth.client-id")
	clientSecret := viper.GetString("grab.auth.client-secret")

	//input := fmt.Sprint(req.AuthorizationID, req.Authorization)
	//auth := fmt.Sprint(clientID, clientSecret)

	input := fmt.Sprint(req.Authorization)
	auth := fmt.Sprint(clientSecret)

	return input == auth
}
