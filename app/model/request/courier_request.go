package request

import (
	"encoding/json"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util/datatype"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters SaveCourierRequest
type ReqCourierBody struct {
	//  in: body
	Body SaveCourierRequest `json:"body"`
}

// swagger:model SaveCourierRequestBody
type SaveCourierRequest struct {
	// Name of the courier
	// in: string
	CourierName string `json:"courier_name"`

	// Code of the courier
	// in: string
	Code string `json:"code"`

	// Uid of the courỉe, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`

	// type of courier
	// in: string
	CourierType string `json:"courier_type"`

	// Logo of courier
	// in: string
	Logo string `json:"logo"`

	// Hide purpose of the Courier
	// in: integer
	HidePurpose int32 `json:"hide_purpose"`

	// Courier Api Integration of the Courier
	// in: integer
	CourierApiIntegration int32 `json:"courier_api_intergration"`

	// Geo Coodinate of the Courier
	// in: string
	UseGeocoodinate int32 `json:"use_geocoodinate"`

	// Provide Airwaybill of the Courier
	// in: integer
	ProvideAirwaybill int32 `json:"provide_airwaybill"`

	// Courier status
	// in: int
	// required: false
	Status int32 `json:"status" binding:"omitempty"`

	// Image UID
	// in: string
	ImageUID string `json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `json:"image_path"`
}

// swagger:parameters CourierByUIdParam
type CourierByUIdParam struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`
}

// swagger:parameters DeleteCourierByUIdParam
type DeleteCourierByUIdParam struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`
}

// swagger:parameters CourierListRequest
type CourierListRequest struct {
	// Filter : {"courier_type":["value","value"],"courier_name":["value","value"],"courier_code":["value","value"],"status":[0,1]}
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

	Filters CourierListFilter `json:"-"`
}

type CourierListFilter struct {
	CourierType []string `json:"courier_type"`
	CourierName []string `json:"courier_name"`
	CourierCode []string `json:"courier_code"`
	Status      []int    `json:"status"`
}

func (m *CourierListRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
}

// swagger:parameters UpdateCourierRequest
type ReqUpdateCourierBody struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`

	//  in: body
	Body UpdateCourierRequest `json:"body"`
}

type UpdateCourierRequest struct {
	// Uid of the courier, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`

	// Name of the courier
	// in: string
	CourierName string `json:"courier_name"`

	// Code of the courier
	// in: string
	Code string `json:"code"`

	// type of courier
	// in: string
	CourierType string `json:"courier_type"`

	// Logo of courier
	// in: string
	Logo string `json:"logo"`

	// Hide purpose of the Courier
	// in: integer
	HidePurpose int `json:"hide_purpose"`

	// Courier Api Integration of the Courier
	// in: integer
	CourierApiIntegration int `json:"courier_api_intergration"`

	// Geo Coodinate of the Courier
	// in: string
	UseGeocoodinate int `json:"use_geocoodinate"`

	// Provide Airwaybill of the Courier
	// in: integer
	ProvideAirwaybill int `json:"provide_airwaybill"`

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`

	// Courier status
	// in: int
	Status int `json:"status" binding:"omitempty"`
}

func (req SaveCourierRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CourierName, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Code, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.CourierType, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Logo, validation.Required.Error(message.ErrReq.Message)),
	)
}
