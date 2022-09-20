package request

import (
	"encoding/json"
	"go-klikdokter/helper/global"
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

	//Extend JWT Info
	global.JWTInfo
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

	//Extend JWT Info
	global.JWTInfo
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
	//Filter : {"channel_name":["value","value"],"courier_name":["value","value"],"channel_code":["value","value"],"status":[0,1]}
	Filter string `schema:"filter" json:"filter,inline"`

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

	Filters ChannelCourierListFilter `json:"-"`
}

type ChannelCourierListFilter struct {
	ChannelName []string `json:"channel_name"`
	ChannelCode []string `json:"channel_code"`
	CourierName []string `json:"courier_name"`
	Status      []int    `json:"status"`
}

func (m *ChannelCourierListRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
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
