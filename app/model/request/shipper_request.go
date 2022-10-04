package request

import "time"

func NewGetPricingDomesticRequest(origin, destination int, input *GetShippingRateRequest) *GetPricingDomestic {
	req := GetPricingDomestic{
		Height: input.TotalHeight,
		Length: input.TotalLength,
		Weight: input.TotalWeight,
		Width:  input.TotalWidth,
		Origin: AreaDetail{
			AreaID:    origin,
			Latitude:  input.Origin.Latitude,
			Longitude: input.Origin.Longitude,
		},
		Destination: AreaDetail{
			AreaID:    destination,
			Latitude:  input.Destination.Latitude,
			Longitude: input.Destination.Longitude,
		},
		Page:      1,
		COD:       false,
		ForOrder:  false,
		ItemValue: input.TotalProductPrice,
	}
	return &req
}

type GetPricingDomestic struct {
	COD         bool       `json:"cod"`
	ForOrder    bool       `json:"for_order"`
	Height      float64    `json:"height"`
	ItemValue   float64    `json:"item_value"`
	Length      float64    `json:"length"`
	Weight      float64    `json:"weight"`
	Width       float64    `json:"width"`
	Page        int        `json:"page"`
	Origin      AreaDetail `json:"origin"`
	Destination AreaDetail `json:"destination"`
}

type AreaDetail struct {
	AreaID    int    `json:"area_id"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lng"`
}

type FindShipperCourierCoverage struct {
	CourierID   uint64
	CountryCode string
	PostalCode  string
	Subdistrict string
}
type CreateOrderShipperPartner struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type CreateOrderShipperCourier struct {
	COD          bool `json:"cod"`
	RateID       int  `json:"rate_id"`
	UseInsurance bool `json:"use_insurance"`
}

type CreateOrderShipperArea struct {
	Address     string `json:"address"`
	AreaID      uint64 `json:"area_id"`
	CountryID   uint64 `json:"country_id"`
	CountryName string `json:"country_name"`
	Direction   string `json:"direction"`
	Lat         string `json:"lat"`
	Long        string `json:"lng"`
	PostalCode  string `json:"postal_code"`
}

type CreateOrderShipperPackage struct {
	Items       []CreateOrderShipperPackageItem `json:"items"`
	Height      float64                         `json:"height"`
	Length      float64                         `json:"length"`
	PackageType int                             `json:"package_type"`
	Price       float64                         `json:"price"`
	Width       float64                         `json:"width"`
	Weight      float64                         `json:"weight"`
}

type CreateOrderShipperPackageItem struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type CreateOrderShipper struct {
	Consignee   CreateOrderShipperPartner `json:"consignee"`
	Consigner   CreateOrderShipperPartner `json:"consigner"`
	Courier     CreateOrderShipperCourier `json:"courier"`
	Destination CreateOrderShipperArea    `json:"destination"`
	Origin      CreateOrderShipperArea    `json:"origin"`
	Package     CreateOrderShipperPackage `json:"package"`
	Coverage    string                    `json:"coverage"`
	ExternalID  string                    `json:"external_id"`
	BestPrices  bool                      `json:"best_prices"`
	ServiceType int                       `json:"service_type"`
	PaymentType string                    `json:"payment_type"`
}

type GetPickUpTimeslot struct {
	TimeZone string `json:"time_zone"`
}

type CreatePickUpOrderShipper struct {
	Data CreatePickUpOrderShipperData `json:"data"`
}

type CreatePickUpOrderShipperData struct {
	OrderActivation CreatePickUpOrderShipperOrderActivation `json:"order_activation"`
}

type CreatePickUpOrderShipperOrderActivation struct {
	OrderID   []string  `json:"order_id"`
	Timezone  string    `json:"timezone"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type CancelOrderShipperRequest struct {
	Reason string `json:"reason"`
}
