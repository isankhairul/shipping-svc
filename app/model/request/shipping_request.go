package request

import (
	"encoding/json"
	"fmt"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"time"
)

//swagger:parameters ShippingRate
type GetShippingRate struct {
	//in: body
	Body GetShippingRateRequest `json:"body"`
}

//swagger:parameters ShippingRateByShippingType
type GetShippingRateByShippingType struct {
	//in: path
	ShippingType string `json:"shipping-type"`
	//in: body
	Body GetShippingRateRequest `json:"body"`
}

type GetShippingRateRequest struct {
	ShippingType        string            `json:"-"`
	ChannelUID          string            `json:"channel_uid"`
	TotalWeight         float64           `json:"total_weight"`
	TotalWidth          float64           `json:"total_width"`
	TotalHeight         float64           `json:"total_heigth"`
	TotalLength         float64           `json:"total_length"`
	TotalProductPrice   float64           `json:"total_product_price"`
	ContainPrescription bool              `json:"contain_prescription"`
	Origin              AreaDetailPayload `json:"origin"`
	Destination         AreaDetailPayload `json:"destination"`
	CourierServiceUID   []string          `json:"courier_service_uid"`
	ChannelCode         string            `json:"-"`
}

func (g *GetShippingRateRequest) CheckCoordinate() (bool, message.Message) {
	if len(g.Origin.Latitude) == 0 || len(g.Origin.Longitude) == 0 ||
		len(g.Destination.Latitude) == 0 || len(g.Destination.Longitude) == 0 {
		return false, message.CoordinateRequiredMsg
	}

	return true, message.SuccessMsg
}

type AreaDetailPayload struct {
	CountryCode string `json:"country_code"`
	PostalCode  string `json:"postal_code"`
	Subdistrict string `json:"subdistrict"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

type ChannelCourierServicePayloadItem struct {
	CourierServiceUID string `json:"courier_service_uid"`
}

type CreateDeliveryPartner struct {
	Name  string `json:"name"`
	UID   string `json:"uid"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type CreateDeiveryArea struct {
	Address      string `json:"address"`
	CountryCode  string `json:"country_code"`
	PostalCode   string `json:"postal_code"`
	Subdistrict  string `json:"subdistrict"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	ProvinceName string `json:"province_name"`
	CityName     string `json:"city_name"`
	DistrictName string `json:"district_name"`
}

type CreateDeliveryProduct struct {
	UID   string  `json:"uid"`
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

type CreateDeliveryPackage struct {
	Product             []CreateDeliveryProduct `json:"product"`
	TotalWeight         float64                 `json:"total_weight"`
	TotalWidth          float64                 `json:"total_width"`
	TotalLength         float64                 `json:"total_length"`
	TotalHeight         float64                 `json:"total_height"`
	TotalProductPrice   float64                 `json:"total_product_price"`
	ContainPrescription uint                    `json:"contain_prescription"`
}

//swagger:parameters CreateDelivery
type CreateDeliveryRequest struct {
	//in:body
	Body CreateDelivery `json:"body"`
}

type CreateDelivery struct {
	ChannelUID        string                `json:"channel_uid"`
	CouirerServiceUID string                `json:"courier_service_uid"`
	OrderNo           string                `json:"order_no"`
	COD               bool                  `json:"cod"`
	UseInsurance      bool                  `json:"use_insurance"`
	Notes             string                `json:"notes"`
	Merchant          CreateDeliveryPartner `json:"merchant"`
	Customer          CreateDeliveryPartner `json:"customer"`
	Origin            CreateDeiveryArea     `json:"origin"`
	Destination       CreateDeiveryArea     `json:"destination"`
	Package           CreateDeliveryPackage `json:"package"`
	Username          string                `json:"username"`
}

func (c *CreateDelivery) CheckCoordinate() (bool, message.Message) {
	if len(c.Origin.Latitude) == 0 || len(c.Origin.Longitude) == 0 ||
		len(c.Destination.Latitude) == 0 || len(c.Destination.Longitude) == 0 {
		return false, message.CoordinateRequiredMsg
	}

	return true, message.SuccessMsg
}

func (c *CreateDelivery) ToCreateOrderShipperPackage() *CreateOrderShipperPackage {
	result := []CreateOrderShipperPackageItem{}
	for _, v := range c.Package.Product {
		result = append(result, CreateOrderShipperPackageItem{
			ID:    0,
			Name:  v.Name,
			Price: v.Price,
			Qty:   v.Qty,
		})
	}

	return &CreateOrderShipperPackage{
		Items:  result,
		Height: c.Package.TotalHeight,
		Length: c.Package.TotalLength,
		Width:  c.Package.TotalWeight,
		Weight: c.Package.TotalWeight,
		Price:  c.Package.TotalProductPrice,
	}
}

func (c *CreateDelivery) ToCreateOrderShipper() *CreateOrderShipper {
	return &CreateOrderShipper{
		Consignee: CreateOrderShipperPartner{
			Name:        c.Customer.Name,
			PhoneNumber: c.Customer.Phone,
		},
		Consigner: CreateOrderShipperPartner{
			Name:        c.Merchant.Name,
			PhoneNumber: c.Merchant.Phone,
		},
		Courier: CreateOrderShipperCourier{
			COD:          c.COD,
			RateID:       0,
			UseInsurance: c.UseInsurance,
		},
		Destination: CreateOrderShipperArea{
			Address:     c.Destination.Address,
			AreaID:      0,
			CountryID:   0,
			CountryName: c.Destination.CountryCode,
			Direction:   "",
			Lat:         c.Destination.Latitude,
			Long:        c.Destination.Longitude,
			PostalCode:  c.Destination.PostalCode,
		},
		Origin: CreateOrderShipperArea{
			Address:     c.Origin.Address,
			AreaID:      0,
			CountryID:   0,
			CountryName: c.Origin.CountryCode,
			Direction:   "",
			Lat:         c.Origin.Latitude,
			Long:        c.Origin.Longitude,
			PostalCode:  c.Origin.PostalCode,
		},
		Package:    *c.ToCreateOrderShipperPackage(),
		ExternalID: c.OrderNo,
		Coverage:   "domestic",
		//BestPrices:  false,
		//ServiceType: 1,
		PaymentType: "postpay",
	}
}

// swagger:parameters OrderShippingTracking
type GetOrderShippingTracking struct {
	// in: path
	// required: true
	UID string `json:"uid"`

	// in: query
	// required: true
	ChannelUID string `schema:"channel_uid" json:"channel_uid"`
}

// swagger:parameters WebhookUpdateStatusShipper
type WebhookUpdateStatusShipperRequest struct {
	// in:body
	Body WebhookUpdateStatusShipper `json:"body"`
}

// swagger:model WebhookUpdateStatusShipperRequest
type WebhookUpdateStatusShipper struct {
	Auth            string         `json:"auth"`
	OrderID         string         `json:"order_id"`
	TrackingID      string         `json:"tracking_id"`
	OrderTrackingID string         `json:"order_tracking_id"`
	ExternalID      string         `json:"external_id"`
	StatusDate      time.Time      `json:"status_date"`
	Internal        ShippingStatus `json:"internal"`
	External        ShippingStatus `json:"external"`
	InternalStatus  ShipperStatus  `json:"internal_status"`
	ExternalStatus  ShipperStatus  `json:"external_status"`
	Awb             string         `json:"awb,omitempty"`

	// Extend Jwt Info
	global.JWTInfo
}

type ShippingStatus struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ShipperStatus struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// swagger:parameters GetOrderShippingList
type GetOrderShippingList struct {
	// Filter : {"order_shipping_uid":["001","002"],"order_no":["001","002"],"channel_code":["kd","hb"],"channel_name":["name","name"],"courier_name":["shipper","shipper"],"shipping_status":["created","request_pickup"],"order_shipping_date_from":["2022-09-09"],"order_shipping_date_to":["2022-09-12"]}
	// in: query
	Filter string `json:"filter"`

	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Sort fields
	// in: string
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`

	Filters GetOrderShippingFilter `json:"-"`
}

type GetOrderShippingFilter struct {
	ChannelCode                []string `json:"channel_code"`
	ChannelName                []string `json:"channel_name"`
	CourierName                []string `json:"courier_name"`
	CourierServicesName        []string `json:"courier_services_name"`
	OrderNo                    []string `json:"order_no"`
	Airwaybill                 []string `json:"airwaybill"`
	ShippingStatus             []string `json:"shipping_status"`
	OrderShippingDateFromArray []string `json:"order_shipping_date_from"`
	OrderShippingDateToArray   []string `json:"order_shipping_date_to"`
	BookingID                  []string `json:"booking_id"`
	MerchantName               []string `json:"merchant_name"`
	CustomerName               []string `json:"customer_name"`
	OrderShippingUID           []string `json:"order_shipping_uid"`

	OrderShippingDateFrom string `json:"-"`
	OrderShippingDateTo   string `json:"-"`
}

func (m *GetOrderShippingList) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}

	if len(m.Filters.OrderShippingDateFromArray) > 0 {
		m.Filters.OrderShippingDateFrom = m.Filters.OrderShippingDateFromArray[0]
	}

	if len(m.Filters.OrderShippingDateToArray) > 0 {
		m.Filters.OrderShippingDateTo = m.Filters.OrderShippingDateToArray[0]
	}
}

// swagger:parameters GetOrderShippingDetail
type GetOrderShippingDetail struct {
	// in: path
	// required: true
	UID string `json:"uid"`
}

// swagger:parameters CancelOrder
type CancelOrder struct {
	// in: path
	// required: true
	UID string `json:"uid"`

	// in: body
	Body CancelOrderBodyRequest `json:"body"`
}

// swagger:model CancelOrderBodyRequest
type CancelOrderBodyRequest struct {
	// example: Stok barang habis
	Reason   string `json:"reason"`
	Username string `json:"username"`
}

// swagger:parameters CancelPickup
type CancelPickup struct {
	// in: path
	// required: true
	UID string `json:"uid"`

	// in: body
	Body CancelPickupBodyRequest `json:"body"`
}

// swagger:model CancelPickupBodyRequest
type CancelPickupBodyRequest struct {
	Username string `json:"username"`
}

type UpdateOrderShippingBody struct {
	ChannelUID         string                        `json:"channel_uid"`
	CourierCode        string                        `json:"courier_code"`
	CourierServiceUID  string                        `json:"courier_service_uid"`
	OrderNo            string                        `json:"order_no"`
	OrderShippingUID   string                        `json:"order_shipping_uid"`
	Airwaybill         string                        `json:"airwaybill"`
	ShippingStatus     string                        `json:"shipping_status"`
	ShippingStatusName string                        `json:"shipping_status_name"`
	Details            UpdateOrderShippingBodyDetail `json:"details"`
	DriverInfo         UpdateOrderShippingDriverInfo `json:"driver_info"`
	UpdatedBy          string                        `json:"update_by"`
	Timestamp          time.Time                     `json:"timestamp"`
}

type UpdateOrderShippingBodyDetail struct {
	ExternalStatusCode        string `json:"external_status_code"`
	ExternalStatusName        string `json:"external_status_name"`
	ExternalStatusDescription string `json:"external_status_description"`
}

type UpdateOrderShippingDriverInfo struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	LicencePlate string `json:"license_plate"`
	TrackingURL  string `json:"tracking_url"`
}

func (u *UpdateOrderShippingDriverInfo) Description() string {
	desc := ""
	if len(u.Name) > 0 {
		desc += fmt.Sprint(" . Driver : ", u.Name)
	}

	if len(u.Phone) > 0 {
		desc += fmt.Sprint(" . Phone : ", u.Phone)
	}

	if len(u.LicencePlate) > 0 {
		desc += fmt.Sprint(" . Plate : ", u.LicencePlate)
	}

	if len(u.TrackingURL) > 0 {
		desc += fmt.Sprint(" . Tracking : ", u.TrackingURL)
	}

	return desc
}

// swagger:parameters GetOrderShippingLabel
type GetOrderShippingLabel struct {
	// in: path
	// required: true
	ChannelUID string `json:"channel-uid"`
	// in: body
	Body GetOrderShippingLabelBody `json:"body"`
}

type GetOrderShippingLabelBody struct {
	OrderShippingUID []string `json:"order_shipping_uid"`
	HideProduct      bool     `json:"hide_product"`
}

// swagger:parameters RepickupOrder
type RepickupOrder struct {
	// in: body
	Body RepickupOrderRequest `json:"body"`
}

type RepickupOrderRequest struct {
	ChannelUID       string `json:"channel_uid"`
	OrderShippingUID string `json:"order_shipping_uid"`
	Username         string `json:"username"`
}

// swagger:parameters ShippingTracking
type GetShippingTracking struct {
	// in: path
	// required: true
	UID string `json:"uid"`

	// in: query
	// required: true
	ChannelUID string `schema:"channel_uid" json:"channel_uid"`
}
