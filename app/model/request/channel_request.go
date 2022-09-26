package request

import (
	"encoding/json"
	"go-klikdokter/helper/global"
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

	// Extend Jwt Info
	global.JWTInfo
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
	// Filter : {"channel_code":["value","value"],"channel_name":["value","value"],"status":[0,1]}
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

	Filters ChannelListFilter `json:"-"`
}

type ChannelListFilter struct {

	// Channel Code
	// in: string
	ChannelCode []string `json:"channel_code"`

	// Channel Name
	// in: string
	ChannelName []string `json:"channel_name"`

	// Channel Status
	// in: int
	Status []int `json:"status"`
}

func (m *ChannelListRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
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

	// Extend Jwt Info
	global.JWTInfo
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
	// Filter : {"channel_code":["value","value"],"channel_name":["value","value"],"courier_name":["value","value"],"status_title":["value","value"], "status_code":["value","value"], "courier_status":["value","value"]}
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

	Filters GetChannelCourierStatusFilter `json:"-"`
}

type GetChannelCourierStatusFilter struct {
	ChannelCode   []string `json:"channel_code"`
	ChannelName   []string `json:"channel_name"`
	CourierName   []string `json:"courier_name"`
	StatusCode    []string `json:"status_code"`
	StatusTitle   []string `json:"status_title"`
	CourierStatus []string `json:"courier_status"`
}

func (m *GetChannelCourierStatusRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
}

// swagger:parameters GetChannelCourierList
type GetChannelCourierListRequest struct {
	//in: path
	ChannelUID string `schema:"channel-uid" json:"channel-uid"`

	// Maximun records per page
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Filter: {"courier_type_code": ["third_party"],"courier_code":["shipper"],"courier_name":["Gojek"],"shipping_type_code":["instant","reguler"],"shipping_name":["Same Day"],"status":[0,1]}
	Filter string `schema:"filter" binding:"omitempty" json:"filter"`

	// Sort fields
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`

	// sort direction
	// enum: asc,desc
	Dir string `schema:"dir" binding:"omitempty" json:"dir"`

	FilterMap map[string]interface{} `json:"-"`
}

func (m *GetChannelCourierListRequest) SetFilterMap() {
	filters := map[string]interface{}{}
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &filters)
		m.FilterMap = filters
	}
}
