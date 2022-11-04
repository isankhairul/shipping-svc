package response

import (
	"fmt"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util/datatype"
	"sort"
	"time"
)

type GetShippingRatePriceRange struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type GetShippingRateError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func SetShippingRateErrorMessage(msg message.Message) GetShippingRateError {
	return GetShippingRateError{fmt.Sprint(msg.Code), msg.Message}
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

//swagger:model ShippingRate
type GetShippingRateResponse struct {
	ShippingTypeCode        string                    `json:"shipping_type_code"`
	ShippingTypeName        string                    `json:"shipping_type_name"`
	ShippingTypeDescription string                    `json:"shipping_type_description"`
	PriceRange              GetShippingRatePriceRange `json:"price_range"`
	EtdMin                  float64                   `json:"etd_min"`
	EtdMax                  float64                   `json:"etd_max"`
	AvailableCode           int                       `json:"available_code"`
	Error                   GetShippingRateError      `json:"error"`
	Services                []GetShippingRateService  `json:"services"`
}

type GetShippingRateResponseList struct {
	//in: body
	Response []GetShippingRateResponse `json:"response"`
}

// common response to get shipping rate from courier
type ShippingRateCommonResponse struct {
	// data of each courier_Service
	// key: courier_code:courier_shipping_code
	Rate map[string]ShippingRateData

	// summary per shippingType
	// key: shipping_type
	Summary map[string]ShippingRateSummary

	// error applied to all courier service of the courier
	// key: courier_code
	CourierMsg map[string]message.Message
}

func (s *ShippingRateCommonResponse) Add(data *ShippingRateCommonResponse) {
	if data == nil {
		return
	}

	for k, v := range data.Rate {
		s.Rate[k] = v
	}

	for k, v := range data.CourierMsg {
		s.CourierMsg[k] = v
	}
}

func (s *ShippingRateCommonResponse) FindShippingCode(courierCode, shippingCode string) ShippingRateData {
	courierShippingCode := global.CourierShippingCodeKey(courierCode, shippingCode)

	//check message from shipping provider
	msg, ok := s.CourierMsg[courierCode]
	if ok && msg != message.SuccessMsg {
		return ShippingRateData{
			AvailableCode: 400,
			Error:         SetShippingRateErrorMessage(msg),
		}
	}

	data, ok := s.Rate[courierShippingCode]
	if !ok {
		return ShippingRateData{
			AvailableCode: 400,
			Error:         SetShippingRateErrorMessage(message.ErrShippingRateNotFound),
		}
	}

	return data
}

func (s *ShippingRateCommonResponse) SummaryPerShippingType(shippingType string, price, etdMax, etdMin float64, status int) {
	summaryData, ok := s.Summary[shippingType]

	if !ok {
		summaryData = ShippingRateSummary{
			PriceRange:    GetShippingRatePriceRange{},
			AvailableCode: 400,
			Error:         SetShippingRateErrorMessage(message.ErrShippingRateNotFound),
		}
	}

	if status == 200 {
		maxPrice := summaryData.PriceRange.MaxPrice
		minPrice := summaryData.PriceRange.MinPrice
		eMax := summaryData.EtdMax
		eMin := summaryData.EtdMin
		summaryData.AvailableCode = 200
		summaryData.Error = SetShippingRateErrorMessage(message.SuccessMsg)

		if maxPrice == 0 || price > maxPrice {
			summaryData.PriceRange.MaxPrice = price
		}

		if minPrice == 0 || price < minPrice {
			summaryData.PriceRange.MinPrice = price
		}

		if eMax == 0 || etdMax > eMax {
			summaryData.EtdMax = etdMax
		}

		if eMin == 0 || etdMin < eMin {
			summaryData.EtdMin = etdMin
		}
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
	Etd_Min          float64
	Etd_Max          float64
	MinDay           int
	MaxDay           int
	UnitPrice        float64
	TotalPrice       float64
	InsuranceFee     float64
	MustUseInsurance bool
	InsuranceApplied bool
	Distance         float64
}

func (s *ShippingRateData) UpdateMessage(msg message.Message) {

	if msg == message.SuccessMsg {
		s.AvailableCode = 200
		s.Error = SetShippingRateErrorMessage(message.SuccessMsg)
		return
	}

	s.AvailableCode = 400
	s.Error = SetShippingRateErrorMessage(msg)
	s.Weight = 0
	s.Volume = 0
	s.VolumeWeight = 0
	s.FinalWeight = 0
	s.MinDay = 0
	s.MaxDay = 0
	s.Etd_Min = 0
	s.Etd_Max = 0
	s.UnitPrice = 0
	s.TotalPrice = 0
	s.InsuranceFee = 0
	s.MustUseInsurance = false
	s.InsuranceApplied = false
	s.Distance = 0
}

type ShippingRateSummary struct {
	PriceRange    GetShippingRatePriceRange
	EtdMin        float64
	EtdMax        float64
	AvailableCode int
	Error         GetShippingRateError
}

//swagger:model CreateDeliveryResponse
type CreateDelivery struct {
	OrderShippingUID string `json:"order_shipping_uid,omitempty"`
	OrderNoAPI       string `json:"order_no_api,omitempty"`
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
	DateTime time.Time `json:"-"`
	//example: 2022-01-31
	Date string `json:"date"`
	//example: 12:30
	Time string `json:"time"`
	//example: Status
	Status string `json:"status"`
	//example: Order Masuk ke sistem
	Note string `json:"note"`
}

type trackingOrderSorter struct {
	data []GetOrderShippingTracking
	by   trackingOrderSorterBy
}

func (t *trackingOrderSorter) Len() int {
	return len(t.data)
}

func (t *trackingOrderSorter) Swap(i, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
}

func (t *trackingOrderSorter) Less(i, j int) bool {
	return t.by(&t.data[i], &t.data[j])
}

type trackingOrderSorterBy func(arg1, arg2 *GetOrderShippingTracking) bool

func (b trackingOrderSorterBy) Sort(data []GetOrderShippingTracking) {
	d := &trackingOrderSorter{
		data: data,
		by:   b,
	}

	sort.Sort(d)
}

func SortOrderStatusByTimeDesc(data []GetOrderShippingTracking) []GetOrderShippingTracking {
	trackingOrderSorterBy(func(data1, data2 *GetOrderShippingTracking) bool {
		return data2.DateTime.Before(data1.DateTime)
	}).Sort(data)
	return data
}

//swagger:response GetOrderShippingList
type GetOrderShippingListResponse struct {
	//in:body
	Body []GetOrderShippingList `json:"body"`
}

//swagger:model GetOrderShippingListResponse
type GetOrderShippingList struct {
	ChannelCode        string    `gorm:"column:channel_code" json:"channel_code"`
	OrderShippingDate  time.Time `gorm:"column:order_shipping_date" json:"order_shipping_date"`
	ChannelName        string    `gorm:"column:channel_name" json:"channel_name"`
	OrderShippingUID   string    `gorm:"column:order_shipping_uid" json:"order_shipping_uid"`
	OrderNo            string    `gorm:"column:order_no" json:"order_no"`
	CourierName        string    `gorm:"column:courier_name" json:"courier_name"`
	CourierServiceName string    `gorm:"column:courier_services_name" json:"courier_services_name"`
	Airwaybill         string    `gorm:"column:airwaybill" json:"airwaybill"`
	BookingID          string    `gorm:"column:booking_id" json:"booking_id"`
	MerchantName       string    `gorm:"column:merchant_name" json:"merchant_name"`
	CustomerName       string    `gorm:"column:customer_name" json:"customer_name"`
	ShippingStatus     string    `gorm:"column:shipping_status" json:"shipping_status"`
	ShippingStatusName string    `gorm:"column:shipping_status_name" json:"shipping_status_name"`
}

//swagger:response GetOrderShippingDetail
type GetOrderShippingDetailResponse struct {
	//in:body
	Body []GetOrderShippingDetail `json:"body"`
}

//swagger:model GetOrderShippingDetailResponse
type GetOrderShippingDetail struct {
	//example: hh6845hjjisdfhidsf
	ChannelUID string `json:"channel_uid"`
	//example: kd
	ChannelCode string `json:"channel_code"`
	//example: Klikdokter
	ChannelName string `json:"channel_name"`
	//example: hh6845hjjisdfhidsf
	OrderShippingUID string `json:"order_shipping_uid"`
	//example: 2022-09-28
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
	//example: 3.2
	TotalLength float64 `json:"total_length"`
	//example: 3.3
	TotalWidth float64 `json:"total_width"`
	//example: 3.4
	TotalHeight float64 `json:"total_height"`
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
	//example: 2210
	MerchantSubdistrict string `json:"merchant_subdistrict"`
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
	//example: 2210
	CustomerSubdistrict string `json:"customer_subdistrict"`
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
	Prescription int `json:"prescription"`
}

//swagger:model GetOrderShippingDetailResponseHistory
type GetOrderShippingDetailHistory struct {
	CreatedAt time.Time `json:"created_at"`
	//example: request_pickup
	Status string `json:"status"`
	//example: Pickup Request with Airwaybill No : 1077400000002
	Notes string `json:"notes"`
	//example: username
	CreatedBy string `json:"created_by"`
	//example: request_pickup
	StatusName string `json:"status_name"`
}

//swagger:model GetOrderShippingLabelResponse
type GetOrderShippingLabelResponse struct {
	ChannelCode          string                       `json:"channel_code"`
	ChannelName          string                       `json:"channel_name"`
	ChannelImage         datatype.JSONB               `json:"channel_image"`
	OrderShippingUID     string                       `json:"order_shipping_uid"`
	OrderShippingDate    time.Time                    `json:"order_shipping_date"`
	OrderNo              string                       `json:"order_no"`
	OrderNoAPI           string                       `json:"order_no_api"`
	CourierName          string                       `json:"courier_name"`
	CourierImage         datatype.JSONB               `json:"courier_image"`
	CourierServiceName   string                       `json:"courier_service_name"`
	CourierServiceImage  datatype.JSONB               `json:"courier_service_image"`
	Airwaybill           string                       `json:"airwaybill"`
	BookingID            string                       `json:"booking_id"`
	TotalProductPrice    float64                      `json:"total_product_price"`
	TotalLength          float64                      `json:"total_length"`
	TotalWidth           float64                      `json:"total_width"`
	TotalHeight          float64                      `json:"total_height"`
	TotalWeight          float64                      `json:"total_weight"`
	TotalVolume          float64                      `json:"total_volume"`
	FinalWeight          float64                      `json:"final_weight"`
	ShippingCost         float64                      `json:"shipping_cost"`
	Insurance            bool                         `json:"insurance"`
	InsuranceCost        float64                      `json:"insurance_cost"`
	TotalShippingCost    float64                      `json:"total_shipping_cost"`
	ShippingNotes        string                       `json:"shipping_notes"`
	MerchantUID          string                       `json:"merchant_uid"`
	MerchantName         string                       `json:"merchant_name"`
	MerchantEmail        string                       `json:"merchant_email"`
	MerchantPhone        string                       `json:"merchant_phone"`
	MerchantAddress      string                       `json:"merchant_address"`
	MerchantDistrictName string                       `json:"merchant_district_name"`
	MerchantCityName     string                       `json:"merchant_city_name"`
	MerchantProvinceName string                       `json:"merchant_province_name"`
	MerchantPostalCode   string                       `json:"merchant_postal_code"`
	CustomerUID          string                       `json:"customer_uid"`
	CustomerName         string                       `json:"customer_name"`
	CustomerEmail        string                       `json:"customer_email"`
	CustomerPhone        string                       `json:"customer_phone"`
	CustomerAddress      string                       `json:"customer_address"`
	CustomerDistrictName string                       `json:"customer_district_name"`
	CustomerCityName     string                       `json:"customer_city_name"`
	CustomerProvinceName string                       `json:"customer_province_name"`
	CustomerPostalCode   string                       `json:"customer_postal_code"`
	CustomerNotes        string                       `json:"customer_notes"`
	OrderShippingItems   []GetOrderShippingDetailItem `json:"order_shipping_items"`
}

//swagger:model RepickupOrderResponse
type RepickupOrderResponse struct {
	OrderShippingUID string `json:"order_shipping_uid,omitempty"`
	OrderNoAPI       string `json:"order_no_api,omitempty"`
	PickupCode       string `json:"pickup_code,omitempty"`
}

// swagger:response DownloadOrderShipping
type DownloadOrderShipping struct {
	Channel              string
	OrderShippingDate    time.Time
	OrderShippingUid     string
	OrderNo              string
	CourierName          string
	CourierService       string
	Airwaybill           string
	BookingId            string
	CustomerName         string
	CustomerPhoneNumber  string
	CustomerEmail        string
	CustomerAddress      string
	CustomerProvinceName string
	CustomerCityName     string
	CustomerDistrictName string
	CustomerSubdistrict  string
	CustomerPostalCode   string
	CustomerNotes        string
	MerchantName         string
	MerchantPhoneNumber  string
	MerchantEmail        string
	MerchantAddress      string
	MerchantProvinceName string
	MerchantCityName     string
	MerchantDistrictName string
	MerchantSubdistrict  string
	MerchantPostalCode   string
	TotalWeight          string
	TotalVolume          string
	TotalProductPrice    string
	TotalFinalWeight     string
	ContainPrescription  string
	Insurance            string
	InsuranceCost        string
	ShippingCost         string
	TotalShippingCost    string
	ActualShippingCost   string
	ShippingNotes        string
	ShippingStatusName   string
	OrderStatusHistory   string
}
