package request

import (
	"encoding/json"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters CreateChannelCourierService
type ReqChannelCourierServiceBody struct {
	//  in: body
	Body SaveChannelCourierServiceRequest `json:"body"`
}

type SaveChannelCourierServiceRequest struct {
	// UID of the Courier Service
	// in: int
	// required: true
	CourierServiceUID string `json:"courier_service_uid"`
	// UID of the channel courier
	// in: int
	// required: true
	ChannelCourierUID string `json:"channel_courier_uid"`

	// Priority Sort of ChannelCourier
	// in: number
	// required: true
	PriceInternal float64 `json:"price_internal" binding:"omitempty"`

	// Courier status
	// in: number
	// required: true
	Status int32 `json:"status" binding:"omitempty"`

	// Extend Jwt Info
	global.JWTInfo
}

func (req SaveChannelCourierServiceRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ChannelCourierUID, validation.Required.Error(message.ErrChannelCourierID.Message)),
		validation.Field(&req.CourierServiceUID, validation.Required.Error(message.ErrCourierServiceID.Message)),
	)
}

// swagger:parameters UpdateChannelCourierService
type UpdateChannelCourierServiceRequest struct {
	// name: uid
	// in: path
	// required: true
	UID string `json:"uid"`
	//  in: body
	Body UpdateChannelCourierService `json:"body"`
}

type UpdateChannelCourierService struct {
	// Priority Sort of ChannelCourier
	// required: true
	PriceInternal float64 `json:"price_internal" binding:"omitempty"`

	// Courier status
	// required: true
	Status int `json:"status" binding:"omitempty"`

	// Extend Jwt Info
	global.JWTInfo
}

// swagger:parameters GetChannelCourierServiceList
type ChannelCourierServiceListRequest struct {
	//Filter : {"channel_name":["value","value"],"courier_uid":["value","value"],"courier_name":["value","value"],"shipping_type":["value","value"],"shipping_type_name":["value","value"],"shipping_code":["value","value"],"shipping_name":["value","value"],"status":[0,1]}
	// in: query
	Filter string `json:"filter"`

	// Maximun records per page
	// in: query
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	// in: query
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Sort fields and direction
	// in: query
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`

	Filters ChannelCourierServiceFilter `json:"-"`
}

// swagger:parameters GetChannelCourierServiceByUID DeleteChannelCourierServiceByUID
type ChannelCourierServiceByUID struct {
	ChannelCourierServiceByUIDPath
}

type ChannelCourierServiceByUIDPath struct {
	// name: uid
	// in: path
	// required: true
	UID string `json:"uid"`
}

func (m *ChannelCourierServiceListRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
}

type ChannelCourierServiceFilter struct {
	ChannelName      []string `json:"channel_name"`
	CourierName      []string `json:"courier_name"`
	Status           []int    `json:"status"`
	ShippingName     []string `json:"shipping_name"`
	ShippingCode     []string `json:"shipping_code"`
	ShippingType     []string `json:"shipping_type"`
	CourierUID       []string `json:"courier_uid"`
	ShippingTypeName []string `json:"shipping_type_name"`
}
