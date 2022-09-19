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
