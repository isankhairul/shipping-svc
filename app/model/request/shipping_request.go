package request

import (
	"encoding/json"
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
	ShippingType        string
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
	ProvinceCode string `json:"province_name"`
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
		Package:     *c.ToCreateOrderShipperPackage(),
		ExternalID:  c.OrderNo,
		Coverage:    "domestic",
		BestPrices:  false,
		ServiceType: 1,
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
	// Filter : {"channel_code":["kd","hb"],"channel_name":["name","name"],"courier_name":["shipper","shipper"],"shipping_status":["created","request_pickup"],"order_shipping_date_from":"2022-09-09","order_shipping_date_to":"2022-09-12"}
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

	// Sort direction
	// in: string
	// enum: asc, desc
	Dir string `schema:"dir" binding:"omitempty" json:"dir"`

	Filters GetOrderShippingFilter `json:"-"`
}

type GetOrderShippingFilter struct {
	ChannelCode           []string `json:"channel_code"`
	ChannelName           []string `json:"channel_name"`
	CourierName           []string `json:"courier_name"`
	ShippingStatus        []string `json:"shipping_status"`
	OrderShippingDateFrom string   `json:"order_shipping_date_from"`
	OrderShippingDateTo   string   `json:"order_shipping_date_to"`
}

func (m *GetOrderShippingList) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
}
