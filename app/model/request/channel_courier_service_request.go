package request

import (
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
}

// swagger:parameters GetChannelCourierServiceList
type ChannelCourierServiceListRequest struct {

	// Channel name
	// in: query
	// collection format: multi
	ChannelName []string `schema:"channel_name" binding:"omitempty" json:"channel_name"`

	// Courier name
	// in: query
	// collection format: multi
	CourierName []string `schema:"courier_name" binding:"omitempty" json:"courier_name"`

	// Channel Courier Service status
	// in: query
	// collection format: multi
	Status []int `binding:"omitempty" json:"status"`

	// Shipping name
	// in: query
	// collection format: multi
	ShippingName []string `schema:"shipping_name" binding:"omitempty" json:"shipping_name"`

	// Shipping code
	// in: query
	// collection format: multi
	ShippingCode []string `schema:"shipping_code" binding:"omitempty" json:"shipping_code"`

	// Shipping type
	// collection format: multi
	ShippingType []string `schema:"shipping_type" binding:"omitempty" json:"shipping_type"`

	// Maximun records per page
	// in: query
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	// in: query
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Sort fields and direction
	// in: query
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`
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
