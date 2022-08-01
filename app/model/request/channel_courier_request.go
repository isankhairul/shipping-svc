package request

import (
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters SaveChannelCourierRequest
type ReqChannelCourierBody struct {
	//  in: body
	Body SaveChannelCourierRequest `json:"body"`
}

type SaveChannelCourierRequest struct {
	// ID of the courier
	// in: int
	// required: true
	CourierUID string `json:"courier_uid"`

	// ID of the channel
	// in: string
	// required: true
	ChannelUID string `json:"channel_uid"`

	// Priority Sort of ChannelCourier
	// in: int
	// required: true
	PrioritySort int32 `json:"priority_sort" binding:"omitempty"`

	// Courier status
	// in: int
	// required: true
	Status int32 `json:"status" binding:"omitempty"`

	// Hide purpose of the Courier
	// in: integer
	// required: true
	HidePurpose int32 `json:"hide_purpose" binding:"omitempty"`
}

type CourierServiceDTO struct {
	// Priority Sort of ChannelCourier
	// in: number
	// required: true
	PriceInternal float64 `json:"price_internal" binding:"omitempty"`

	// Courier status
	// in: number
	// required: true
	Status int `json:"status" binding:"omitempty"`

	// Courier Service Uid
	// in: string
	// required: true
	CourierServiceUid string `json:"courier_service_uid" binding:"omitempty"`
}

// swagger:parameters UpdateChannelCourierRequest
type ReqUpdateChannelCourierBody struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`
	//  in: body
	Body UpdateChannelCourierRequest `json:"body"`
}

type UpdateChannelCourierRequest struct {
	// Uid of the courỉe, use this on UPDATE function
	// in: int32
	Uid string `json:"-" binding:"omitempty"`

	// Priority Sort of ChannelCourier
	// in: int
	// required: true
	PrioritySort int32 `json:"priority_sort" binding:"omitempty"`

	// Courier status
	// in: int
	// required: true
	Status int32 `json:"status" binding:"omitempty"`

	// Hide purpose of the Courier
	// in: integer
	// required: true
	HidePurpose int32 `json:"hide_purpose" binding:"omitempty"`
}

func (req SaveChannelCourierRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ChannelUID, validation.Required.Error(message.ErrChannelID.Message)),
		validation.Field(&req.CourierUID, validation.Required.Error(message.ErrCourierID.Message)),
		validation.Field(&req.PrioritySort, validation.In(1, 999).Error(message.ErrPrioritySort.Message)),
	)
}

// swagger:parameters ChannelCourierListRequest
type ChannelCourierListRequest struct {

	// Channel name records
	// in: string
	// collection format: multi
	ChannelName []string `schema:"channel_name" binding:"omitempty" json:"channel_name"`

	// Channel code records
	// in: string
	// collection format: multi
	ChannelCode []string `schema:"channel_code" binding:"omitempty" json:"channel_code"`

	// Courier name records
	// in: string
	// collection format: multi
	CourierName []string `schema:"courier_name" binding:"omitempty" json:"courier_name"`

	// Channel Courier status
	// in: int
	// collection format: multi
	Status []int `binding:"omitempty" json:"status"`

	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Sort fields
	// in: string
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`

	// Sort fields
	// in: string
	Dir string `schema:"dir" binding:"omitempty" json:"dir"`
}

// swagger:parameters GetChannelCourierByUid
type GetChannelCourierByUid struct {
	ChannelCourierByUid
}

// swagger:parameters DeleteChannelCourierByUid
type DeleteChannelCourierByUid struct {
	ChannelCourierByUid
}

type ChannelCourierByUid struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`
}
