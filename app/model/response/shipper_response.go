package response

import (
	"fmt"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"time"
)

type GetPricingDomestic struct {
	Metadata   ShipperMetaData        `json:"metadata"`
	Data       GetPricingDomesticData `json:"data"`
	Pagination ShipperPagination      `json:"pagination"`
}

type CreateOrderShipperResponse struct {
	Metadata ShipperMetaData    `json:"metadata"`
	Data     CreateOrderShipper `json:"data"`
}

type GetPickUpTimeslotResponse struct {
	Metadata ShipperMetaData   `json:"metadata"`
	Data     GetPickUpTimeslot `json:"data"`
}

type CreatePickUpOrderShipperResponse struct {
	Metadata ShipperMetaData          `json:"metadata"`
	Data     CreatePickUpOrderShipper `json:"data"`
}

type GetOrderDetailResponse struct {
	Metadata ShipperMetaData `json:"metadata"`
	Data     GetOrderDetail  `json:"data"`
}

type MetadataResponse struct {
	Metadata ShipperMetaData `json:"metadata"`
}

type ShipperMetaData struct {
	Path           string            `json:"path"`
	HTTPStatusCode int               `json:"http_status_code"`
	HTTPStatus     string            `json:"http_status"`
	Timestamp      uint64            `json:"timestamp"`
	Errors         []message.Message `json:"errors"`
}

type ShipperPagination struct {
	CurrentPage     int      `json:"current_page"`
	CurrentElements int      `json:"current_elements"`
	TotalPages      int      `json:"total_pages"`
	TotalElements   int      `json:"total_elements"`
	SortBy          []string `json:"sort_by"`
}

type GetPricingDomesticData struct {
	Origin      DataAreaDetail `json:"origin"`
	Destination DataAreaDetail `json:"destination"`
	Pricings    []PricingsItem `json:"pricings"`
}

type DataAreaDetail struct {
	// "area_id": 2482,
	// "area_name": "Sampora",
	// "suburb_id": 240,
	// "suburb_name": "Cisauk",
	// "city_id": 21,
	// "city_name": "Tangerang, Kab.",
	// "province_id": 3,
	// "province_name": "Banten",
	// "country_id": 228,
	// "country_name": "INDONESIA",
	// "lat": -6.3124856,
	// "lng": 106.6497623
	AreaID    int     `json:"area_id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type PricingsItem struct {
	Logistic         PricingLogisticDetail `json:"logistic"`
	Rate             PricingRateDetail     `json:"rate"`
	Weight           float64               `json:"weight"`
	Volume           float64               `json:"volume"`
	VolumeWeight     float64               `json:"volume_weight"`
	FinalWeight      float64               `json:"final_weight"`
	MinDay           int                   `json:"min_day"`
	MaxDay           int                   `json:"max_day"`
	UnitPrice        float64               `json:"unit_price"`
	TotalPrice       float64               `json:"total_price"`
	Discount         float64               `json:"discount"`
	DiscountValue    float64               `json:"discount_value"`
	DiscountedPrice  float64               `json:"discounted_price"`
	InsuranceFee     float64               `json:"insurance_fee"`
	MustUseInsurance bool                  `json:"must_use_insurance"`
	LiabilityValue   float64               `json:"liability_value"`
	FinalPrice       float64               `json:"final_price"`
	Currency         string                `json:"currency"`
	InsuranceApplied bool                  `json:"insurance_applied"`
	BasePrice        float64               `json:"base_price"`
	SurchargeFee     float64               `json:"surcharge_fee"`
}

type PricingLogisticDetail struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	Code        string `json:"code"`
	CompanyName string `json:"company_name"`
}

type PricingRateDetail struct {
	ID              uint64 `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Description     string `json:"description"`
	FullDescription string `json:"full_description"`
	IsHubless       bool   `json:"is_hubless"`
}

func (g *GetPricingDomestic) ToShippingRate() *ShippingRateCommonResponse {
	if g == nil {
		return nil
	}

	data := map[string]ShippingRateData{}

	for _, v := range g.Data.Pricings {
		courierShippingCode := global.CourierShippingCodeKey("shipper", fmt.Sprint(v.Rate.ID))
		data[courierShippingCode] = ShippingRateData{
			AvailableCode:    200,
			Error:            GetShippingRateError{Message: message.SuccessMsg.Message},
			Weight:           v.Weight,
			Volume:           v.Volume,
			VolumeWeight:     v.VolumeWeight,
			FinalWeight:      v.FinalWeight,
			MinDay:           v.MinDay,
			MaxDay:           v.MaxDay,
			UnitPrice:        v.UnitPrice,
			TotalPrice:       v.TotalPrice,
			InsuranceFee:     v.InsuranceFee,
			MustUseInsurance: v.MustUseInsurance,
			InsuranceApplied: v.InsuranceApplied,
			Distance: util.CalculateDistanceInKm(g.Data.Origin.Latitude,
				g.Data.Origin.Longitude,
				g.Data.Destination.Latitude,
				g.Data.Destination.Longitude),
		}
	}

	return &ShippingRateCommonResponse{
		Rate:       data,
		Summary:    make(map[string]ShippingRateSummary),
		CourierMsg: make(map[string]message.Message),
	}
}

type CreateOrderShipperPartner struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type CreateOrderShipperCourier struct {
	RateID          int     `json:"rate_id"`
	UseInsurance    bool    `json:"use_insurance"`
	Amount          float64 `json:"amount"`
	InsuranceAmount float64 `json:"insurance_amount"`
	COD             bool    `json:"cod"`
}

type CreateOrderShipperArea struct {
	Address      string `json:"address"`
	AreaID       uint64 `json:"area_id"`
	AreaName     string `json:"area_name"`
	CityID       uint64 `json:"city_id"`
	CityName     string `json:"city_name"`
	CountryID    uint64 `json:"country_id"`
	CountryName  string `json:"country_name"`
	Lat          string `json:"lat"`
	Long         string `json:"lng"`
	PostalCode   string `json:"postal_code"`
	ProvinceID   uint64 `json:"province_id"`
	ProvinceName string `json:"province_name"`
	SuburbID     uint64 `json:"suburb_id"`
	SuburbName   string `json:"suburb_name"`
	EmailAddress string `json:"email_address"`
	CompanyName  string `json:"company_name"`
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
	Qty   uint64  `json:"qty"`
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
	OrderID     string                    `json:"order_id"`
	PaymentType string                    `json:"payment_type"`
}

type GetPickUpTimeslot struct {
	Timezone  string     `json:"time_zone"`
	Timeslots []Timeslot `json:"time_slots"`
}

type Timeslot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type CreatePickUpOrderOrderActivation struct {
	OrderID     string    `json:"order_id"`
	PickUpCode  string    `json:"pickup_code"`
	IsActivated bool      `json:"is_activate"`
	PickUpTime  time.Time `json:"pickup_time"`
}

type CreatePickUpOrderShipper struct {
	OrderActivation []CreatePickUpOrderOrderActivation `json:"order_activations"`
}

type GetOrderDetail struct {
	Consignee        interface{}              `json:"consignee"`
	Consigner        interface{}              `json:"consigner"`
	Origin           interface{}              `json:"origin"`
	Destination      interface{}              `json:"destination"`
	ExternalID       string                   `json:"external_id"`
	OrderID          string                   `json:"order_id"`
	Courier          interface{}              `json:"courier"`
	Package          interface{}              `json:"package"`
	PaymentType      string                   `json:"payment_type"`
	Driver           interface{}              `json:"driver"`
	LabelCheckSum    string                   `json:"label_check_sum"`
	CreationDate     time.Time                `json:"creation_date"`
	LastUpdatedDate  time.Time                `json:"last_updated_date"`
	AWBNumber        string                   `json:"awb_number"`
	Trackings        []GetOrderDetailTracking `json:"trackings"`
	IsActive         bool                     `json:"is_active"`
	IsHubless        bool                     `json:"is_hubless"`
	PickUpCode       string                   `json:"pickup_code"`
	PickUpTime       string                   `json:"pickup_time"`
	ShipmentStatus   interface{}              `json:"shipment_status"`
	ProofOfDelivery  interface{}              `json:"proof_of_delivery"`
	TimeSlotSelected interface{}              `json:"time_slot_selected"`
}

type GetOrderDetailTracking struct {
	ShipperStatus  GetOrderDetailTrackingStatus `json:"shipper_status"`
	LogisticStatus GetOrderDetailTrackingStatus `json:"logistic_status"`
	CreatedDate    time.Time                    `json:"created_date"`
}

type GetOrderDetailTrackingStatus struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g *GetOrderDetail) ToOrderShippingTracking() []GetOrderShippingTracking {
	resp := []GetOrderShippingTracking{}
	codes := make(map[string]bool)
	for _, v := range g.Trackings {
		if _, ok := codes[v.LogisticStatus.Name]; ok {
			continue
		}
		codes[v.LogisticStatus.Name] = true
		resp = append(resp, GetOrderShippingTracking{
			Note: v.LogisticStatus.Name,
			Date: v.CreatedDate.In(util.Loc).Format(util.LayoutDateOnly),
			Time: v.CreatedDate.In(util.Loc).Format(util.LayoutTimeOnly),
		})
	}
	return resp
}
