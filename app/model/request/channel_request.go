package request

import (
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters SaveChannelRequest
type ReqChannelBody struct {
	//  in: body
	Body SaveChannelRequest `json:"body"`
}

type SaveChannelRequest struct {
	// ChannelName of the Channel
	// in: string
	ChannelName string `json:"channel_name" binding:"omitempty"`

	// ChannelCode of the Channel
	// in: string
	ChannelCode string `json:"channel_code" binding:"omitempty"`

	// Description of the Channel
	// in: string
	Description string `json:"description" binding:"omitempty"`

	// status of Channel
	// in: int
	Status *int `json:"status" binding:"omitempty"`

	// Logo of Channel
	// in: string
	Logo string `json:"logo" binding:"omitempty"`
}

// swagger:parameters Channel
type GetChannelRequest struct {
	// name: id
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters Channels
type ChannelListRequest struct {
	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields, example: name asc, uom desc
	// in: string
	Sort string `schema:"sort" binding:"omitempty"`

	// Channel Code
	// in: string
	ChannelCode string `schema:"ChannelCode" binding:"omitempty"`

	// Channel Name
	// in: string
	ChannelName string `schema:"ChannelName" binding:"omitempty"`

	// Channel Status
	// in: int
	Status int `schema:"Status" binding:"omitempty"`
}

// swagger:parameters UpdateChannelRequest
type ReqChannelBodyUpdate struct {
	// Uid of the Channel
	// in: path
	// required: true
	UId string `json:"uid"`
	//  in: body
	Body UpdateChannelRequest `json:"body"`
}

type UpdateChannelRequest struct {
	// Uid of the courier, use this on UPDATE function
	// in: int32
	Uid string `json:"-" binding:"omitempty"`

	// ChannelName of the Channel
	// in: string
	ChannelName string `json:"channel_name"`

	// ChannelCode of Channel
	// in: string
	ChannelCode string `json:"channel_code"`

	// Description of the Channel
	// in: string
	Description string `json:"description" binding:"omitempty"`

	// Logo of Channel
	// in: string
	Logo string `json:"logo"`

	// Status of Channel
	// in: int
	Status int `json:"status" binding:"omitempty"`
}

// swagger:parameters ChannelRequestGetByUid
type ChannelGetByUid struct {
	// Uid of the Channel
	// in: path
	// required: true
	UId string `json:"uid"`
}

// swagger:parameters ChannelRequestDeleteByUid
type ChannelDeleteByUid struct {
	// Uid of the Channel
	// in: path
	// required: true
	UId string `json:"uid"`
}

func (req SaveChannelRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ChannelName, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.ChannelCode, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Description, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Status, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Logo, validation.Required.Error(message.ErrReq.Message)),
	)
}
