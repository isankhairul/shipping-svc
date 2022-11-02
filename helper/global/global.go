package global

import (
	"fmt"
	"html"
	"reflect"
	"strings"
)

var (
	//set in main.go
	PrefixBase = ""

	//Handler Prefix
	PrefixChannel               = "/channel/"
	PrefixCourier               = "/courier/"
	PrefixCourierCoverageCode   = "/courier/courier-coverage-code/"
	PrefixChannelCourier        = "/channel/channel-courier/"
	PrefixChannelCourierService = "/channel/channel-courier-service/"
	PrefixShipping              = "/shipping/"
	PrefixWebhook               = "/public/webhook/"
	PrefixOther                 = "/other/"

	//Path
	PathUID    = "{uid}"
	PathImport = "import"

	PathChannelApp    = "channel-app"
	PathChannelAppUID = "channel-app/{uid}"

	PathCourier    = "courier"
	PathCourierUID = "courier/{uid}"

	PathCourierService    = "courier-services"
	PathCourierServiceUID = "courier-services/{uid}"

	PathShipmentPredefined    = "shipment-predefined"
	PathShipmentPredefinedUID = "shipment-predefined/{uid}"

	PathChannelCourierStatus = "channel-status-courier-status"
	PathUIDCourierList       = "{uid}/courier-list"

	PathShippingRate             = "shipping-rate"
	PathShippingRateShippingType = "shipping-rate/{shipping-type}"
	PathOrderShipping            = "order-shipping"
	PathOrderShippingDownload    = "order-shipping/download"
	PathOrderTracking            = "order-tracking/{uid}"
	PathShipper                  = "shipper"
	PathGrab                     = "grab"
	PathOrderShippingUID         = "order-shipping/{uid}"
	PathCancelPickupUID          = "cancel-pickup/{uid}"
	PathCancelOrderUID           = "cancel-order/{uid}"
	PathOrderShippingLabel       = "order-shipping-label/{channel-uid}"
	PathRepickup                 = "repickup"
	PathShippingTracking         = "tracking/{uid}"

	ServerPort = "server.port"
)

func HtmlEscape(req interface{}) {
	value := reflect.ValueOf(req).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type() != reflect.TypeOf("") {
			continue
		}

		str := field.Interface().(string)
		field.SetString(html.EscapeString(str))
	}
}

func AddLike(column string, value []string) string {
	var condition string
	for _, v := range value {
		condition += fmt.Sprintf(" LOWER(%s) ILIKE '%%%s%%' OR", column, v)
	}
	return strings.TrimRight(condition, " OR")
}

func CourierShippingCodeKey(courierCode, shippingCode string) string {
	return strings.ToLower(fmt.Sprint(courierCode, ":", shippingCode))
}
