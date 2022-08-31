package global

import (
	"fmt"
	"html"
	"math"
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

	PathShippingRate = "shipping-rate"

	ServerPort = "server.port"
)

const (
	earthRaidusKm = 6371 // radius of the earth in kilometers.
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

func ConvertDegToRad(d float64) float64 {
	return d * math.Pi / 180
}

func DistanceKM(lat1, long1, lat2, long2 float64) float64 {
	lat1 = ConvertDegToRad(lat1)
	long1 = ConvertDegToRad(long1)
	lat2 = ConvertDegToRad(lat2)
	long2 = ConvertDegToRad(long2)

	dLat := lat2 - lat1
	dLong := long2 - long1

	//a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLong/2), 2)

	//c = 2 ⋅ atan2( √a, √(1−a) )
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	//d = R * c
	d := earthRaidusKm * c
	return d
}

func CourierShippingCodeKey(courierCode, shippingCode string) string {
	return fmt.Sprint(courierCode, ":", shippingCode)
}
