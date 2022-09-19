package response

import (
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util/datatype"
	"time"
)

type GetShippingRatePriceRange struct {
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
}

type GetShippingRateError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type GetShippingRateCourir struct {
	CourierID       uint64 `json:"-"`
	CourierUID      string `json:"courier_uid"`
	CourierCode     string `json:"courier_code"`
	CourierName     string `json:"courier_name"`
	CourierTypeCode string `json:"courier_type_code"`
	CourierTypeName string `json:"courier_type_name"`
}

type GetShippingRateService struct {
	Courier                 GetShippingRateCourir `json:"courier"`
	CourierServiceUID       string                `json:"courier_service_uid"`
	ShippingCode            string                `json:"shipping_code"`
	ShippingName            string                `json:"shipping_name"`
	ShippingTypeDescription string                `json:"shipping_type_description"`
	Logo                    datatype.JSONB        `json:"logo"`
	ShippingTypeCode        string                `json:"shipping_type_code"`
	ShippingTypeName        string                `json:"shipping_type_name"`
	Etd_Min                 float64               `json:"etd_min"`
	Etd_Max                 float64               `json:"etd_max"`
	AvailableCode           int                   `json:"available_code"`
	Error                   GetShippingRateError  `json:"error"`
	Weight                  float64               `json:"weight"`
	Volume                  float64               `json:"volume"`
	VolumeWeight            float64               `json:"volume_weight"`
	FinalWeight             float64               `json:"final_weight"`
	MinDay                  int                   `json:"min_day"`
	MaxDay                  int                   `json:"max_day"`
	UnitPrice               float64               `json:"unit_price"`
	TotalPrice              float64               `json:"total_price"`
	InsuranceFee            float64               `json:"insurance_fee"`
	MustUseInsurance        bool                  `json:"must_use_insurance"`
	InsuranceApplied        bool                  `json:"insurance_applied"`
	Distance                float64               `json:"distance"`
}

func (g *GetShippingRateService) FromShipper(val PricingsItem) {
	g.Weight = val.Weight
	g.Volume = val.Volume
	g.VolumeWeight = val.VolumeWeight
	g.FinalWeight = val.FinalWeight
	g.MinDay = val.MinDay
	g.MaxDay = val.MaxDay
	g.UnitPrice = val.UnitPrice
	g.TotalPrice = val.TotalPrice
	g.InsuranceFee = val.InsuranceFee
	g.MustUseInsurance = val.MustUseInsurance
	g.InsuranceApplied = val.InsuranceApplied
}

type GetShippingRateResponse struct {
	ShippingTypeCode        string                    `json:"shipping_type_code"`
	ShippingTypeName        string                    `json:"shipping_type_name"`
	ShippingTypeDescription string                    `json:"shipping_type_description"`
	PriceRange              GetShippingRatePriceRange `json:"price_range"`
	EtdMin                  *float64                  `json:"etd_min"`
	EtdMax                  *float64                  `json:"etd_max"`
	AvailableCode           int                       `json:"available_code"`
	Error                   GetShippingRateError      `json:"error"`
	Services                []GetShippingRateService  `json:"services"`
}

//swagger:response ShippingRate
type GetShippingRateResponseList struct {
	//in: body
	Response []GetShippingRateResponse `json:"response"`
}

type ShippingRateCommonResponse struct {
	Rate    map[string]ShippingRateData
	Summary map[string]ShippingRateSummary
	Msg     message.Message
}

func (s *ShippingRateCommonResponse) Add(data map[string]ShippingRateData) {
	for k, v := range data {
		s.Rate[k] = v
	}
}

func (s *ShippingRateCommonResponse) FindShippingCode(shippingCode string) ShippingRateData {
	data, ok := s.Rate[shippingCode]

	if s.Msg.Code != message.SuccessMsg.Code && s.Msg.Code != 0 {
		return ShippingRateData{
			AvailableCode: 400,
			Error: GetShippingRateError{
				Message: s.Msg.Message,
			},
		}
	}

	if !ok {
		return ShippingRateData{
			AvailableCode: 400,
			Error: GetShippingRateError{
				Message: message.ErrShippingRateNotFound.Message,
			},
		}
	}
	return data
}

func (s *ShippingRateCommonResponse) SummaryPerShippingType(shippingType string, price, etdMax, etdMin float64) {
	summaryData, ok := s.Summary[shippingType]

	if !ok {
		summaryData = ShippingRateSummary{
			PriceRange: GetShippingRatePriceRange{},
		}
	}

	maxPrice := summaryData.PriceRange.MaxPrice
	minPrice := summaryData.PriceRange.MinPrice
	eMax := summaryData.EtdMax
	eMin := summaryData.EtdMin

	if maxPrice == nil || price > *maxPrice {
		summaryData.PriceRange.MaxPrice = &price
	}

	if minPrice == nil || price < *minPrice {
		summaryData.PriceRange.MinPrice = &price
	}

	if eMax == nil || etdMax > *eMax {
		summaryData.EtdMax = &etdMax
	}

	if eMin == nil || etdMin < *eMin {
		summaryData.EtdMin = &etdMin
	}

	s.Summary[shippingType] = summaryData
}

type ShippingRateData struct {
	AvailableCode    int
	Error            GetShippingRateError
	Weight           float64
	Volume           float64
	VolumeWeight     float64
	FinalWeight      float64
	MinDay           int
	MaxDay           int
	UnitPrice        float64
	TotalPrice       float64
	InsuranceFee     float64
	MustUseInsurance bool
	InsuranceApplied bool
	Distance         float64
}

func (s *ShippingRateData) SetMessage(isErr bool, msg message.Message) {
	if s.AvailableCode == 400 {
		return
	}

	if !isErr {
		s.AvailableCode = 200
		s.Error = GetShippingRateError{Message: message.SuccessMsg.Message}
		return
	}

	s.AvailableCode = 400
	s.Error = GetShippingRateError{Message: msg.Message}
}

type ShippingRateSummary struct {
	PriceRange GetShippingRatePriceRange
	EtdMin     *float64
	EtdMax     *float64
}

//swagger:response CreateDelivery
type CreateDeliveryResponse struct {
	//in:body
	Body CreateDelivery `json:"body"`
}

//swagger:model CreateDeliveryResponse
type CreateDelivery struct {
	OrderShippingUID string `json:"order_shipping_uid"`
	OrderNoAPI       string `json:"order_no_api"`
}

type CreateDeliveryThirdPartyData struct {
	Insurance          bool
	InsuranceCost      float64
	ShippingCost       float64
	TotalShippingCost  float64
	ActualShippingCost float64
	BookingID          string
	Status             string
	Airwaybill         string

	PickUpTime time.Time
	PickUpCode string
}

//swagger:response OrderShippingTracking
type GetOrderShippingTrackingResponse struct {
	//in:body
	Body []GetOrderShippingTracking `json:"body"`
}

//swagger:model GetOrderShippingTrackingResponse
type GetOrderShippingTracking struct {
	//example: 2022-01-31
	Date string `json:"date"`
	//example: 12:30
	Time string `json:"time"`
	//example: Order Masuk ke sistem
	Note string `json:"note"`
}

//swagger:response GetOrderShippingList
type GetOrderShippingListResponse struct {
	//in:body
	Body []GetOrderShippingList `json:"body"`
}

//swagger:model GetOrderShippingListResponse
type GetOrderShippingList struct {
	ChannelCode        string `gorm:"column:channel_code" json:"channel_code"`
	ChannelName        string `gorm:"column:channel_name" json:"channel_name"`
	OrderShippingUID   string `gorm:"column:order_shipping_uid" json:"order_shipping_uid"`
	OrderNo            string `gorm:"column:order_no" json:"order_no"`
	CourierName        string `gorm:"column:courier_name" json:"courier_name"`
	CourierServiceName string `gorm:"column:courier_services_name" json:"courier_services_name"`
	Airwaybill         string `gorm:"column:airwaybill" json:"airwaybill"`
	BookingID          string `gorm:"column:booking_id" json:"booking_id"`
	MerchantName       string `gorm:"column:merchant_name" json:"merchant_name"`
	CustomerName       string `gorm:"column:customer_name" json:"customer_name"`
	ShippingStatus     string `gorm:"column:shipping_status" json:"shipping_status"`
	ShippingStatusName string `gorm:"column:shipping_status_name" json:"shipping_status_name"`
}
//swagger:response GetOrderShippingDetail
type GetOrderShippingDetailResponse struct {
	//in:body
	Body []GetOrderShippingDetail `json:"body"`
}

//swagger:model GetOrderShippingDetailResponse
type GetOrderShippingDetail struct {
	//example: kd
	ChannelCode string `json:"channel_code"`
	//example: Klikdokter
	ChannelName string `json:"channel_name"`
	//example: hh6845hjjisdfhidsf
	OrderShippingUID string `json:"order_shipping_uid"`

	OrderShippingDate time.Time `json:"order_shipping_date"`
	//example: 1000363553.1
	OrderNo string `json:"order_no"`
	//example: 1000363553.1
	OrderNoAPI string `json:"order_no_api"`
	//example: delivered
	ShippingStatus string `json:"shipping_status"`
	//example: Delivered
	ShippingStatusName string `json:"shipping_status_name"`
	//example: Shipper
	CourierName string `json:"courier_name"`
	//example: Sicepat Next Day
	CourierServiceName string `json:"courier_services_name"`
	//example: 13421BBFGXZ
	Airwaybill string `json:"airwaybill"`
	//example: 11222449
	BookingID string `json:"booking_id"`
	//example: 150000
	TotalProductPrice float64 `json:"total_product_price"`
	//example: 3.5
	TotalWeight float64 `json:"total_weight"`
	//example: 4.5
	TotalVolume float64 `json:"total_volume"`
	//example: 4.5
	FinalWeight float64 `json:"final_weight"`
	//example: 20000
	ShippingCost float64 `json:"shipping_cost"`
	//example: false
	Insurance bool `json:"insurance"`
	//example: 10000
	InsuranceCost float64 `json:"insurance_cost"`
	//example: 2000
	TotalShippingCost float64 `json:"total_shipping_cost"`
	//example: Notes
	ShippingNotes string `json:"shipping_notes"`
	//example: fhdsfg0376762345dfg
	MerchantUID string `json:"merchant_uid"`
	//example: Fahmi Store
	MerchantName string `json:"merchant_name"`
	//example: fahmi@test.com
	MerchantEmail string `json:"merchant_email"`
	//example: 6208734567345
	MerchantPhone string `json:"merchant_phone"`
	//example: Jl. BSD Grand Boulevard, BSD Green Office Park, BSD City My Republic Plaza (Green Office Park 6
	MerchantAddress string `json:"merchant_address"`
	//example: Cisauk
	MerchantDistrictName string `json:"merchant_district_name"`
	//example: Kabupaten Tangerang
	MerchantCityName string `json:"merchant_city_name"`
	//example: Banten
	MerchantProvinceName string `json:"merchant_province_name"`
	//example: 15345
	MerchantPostalCode string `json:"merchant_postal_code"`
	//example: lsdkosdfg0376762345dfg
	CustomerUID string `json:"customer_uid"`
	//example: Bapak Budiman
	CustomerName string `json:"customer_name"`
	//example: budi@test.com
	CustomerEmail string `json:"customer_email"`
	//example: 620876324562323
	CustomerPhone string `json:"customer_phone"`
	//example: Graha Kirana, Jl. Mitra Sunter Bulevar No.16, RW.11, Sunter Jaya
	CustomerAddress string `json:"customer_address"`
	//example: Tanjung Priok
	CustomerDistrictName string `json:"customer_district_name"`
	//example: Jakarta Utara
	CustomerCityName string `json:"customer_city_name"`
	//example: DKI Jakarta
	CustomerProvinceName string `json:"customer_province_name"`
	//example: 13360
	CustomerPostalCode string `json:"customer_postal_code"`
	//example: antar paket keruang mail room, samping lobby
	CustomerNotes        string                          `json:"customer_notes"`
	OrderShippingItem    []GetOrderShippingDetailItem    `json:"order_shipping_item"`
	OrderShippingHistory []GetOrderShippingDetailHistory `json:"order_shipping_history"`
}

//swagger:model GetOrderShippingDetailResponseItem
type GetOrderShippingDetailItem struct {
	//example: Chil Kid Strawberry 900 gr
	ItemName string `json:"item_name"`
	//example: sssabcgz88877
	ProductUID string `json:"product_uid"`
	//example: 3
	Qty int `json:"qty"`
	//example: 110000
	Price float64 `json:"price"`
	//example: 1.3
	Weight float64 `json:"weight"`
	//example: 2.5
	Volume float64 `json:"volume"`
	//example: 1
	Prescrition int `json:"prescription"`
}

//swagger:model GetOrderShippingDetailResponseHistory
type GetOrderShippingDetailHistory struct {
	CreatedAt time.Time `json:"created_at"`
	//example: request_pickup
	Status string `json:"status"`
	//example: Pickup Request with Airwaybill No : 1077400000002
	Notes string `json:"notes"`
}
