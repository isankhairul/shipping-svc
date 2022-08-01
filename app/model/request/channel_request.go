package request

import (
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util/datatype"

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
	Status *int32 `json:"status" binding:"omitempty"`

	// Logo of Channel
	// in: string
	Logo string `json:"logo" binding:"omitempty"`

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`
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

	// Sort fields
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

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`
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

// swagger:parameters GetChannelCourierStatus
type GetChannelCourierStatusRequest struct {
	// Channel name
	// in: query
	// collection format: multi
	ChannelName []string `schema:"channel_name" binding:"omitempty" json:"channel_name"`

	// Courier name
	// in: query
	// collection format: multi
	CourierName []string `schema:"courier_name" binding:"omitempty" json:"courier_name"`

	// Status code
	// in: query
	// collection format: multi
	StatusCode []string `schema:"status_code" binding:"omitempty" json:"status_code"`

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
