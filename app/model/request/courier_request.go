package request

import (
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters SaveCourierRequest
type ReqCourierBody struct {
	//  in: body
	Body SaveCourierRequest `json:"body"`
}

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
	HidePurpose int `gorm:"not null;default:0" json:"hide_purpose"`

	// Courier Api Integration of the Courier
	// in: integer
	CourierApiIntegration int `gorm:"not null;default:1" json:"courier_api_intergration"`

	// Geo Coodinate of the Courier
	// in: string
	UseGeocoodinate int `gorm:"not null;default:0" json:"use_geocoodinate"`

	// Provide Airwaybill of the Courier
	// in: integer
	ProvideAirwaybill int `gorm:"not null;default:0" json:"provide_airwaybill"`

	// Courier status
	// in: int
	// required: false
	Status int `json:"status" binding:"omitempty"`
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
	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Sort fields, example: name asc, uom desc
	// in: string
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`

	// Courier type
	// in: string
	CourierType string `schema:"courier_type" binding:"omitempty" json:"courier_type"`

	// Courier name
	// in: string
	CourierName string `schema:"courier_name" binding:"omitempty" json:"courier_name"`

	// Courier code
	// in: string
	CourierCode string `schema:"courier_code" binding:"omitempty" json:"courier_code"`

	// Courier status
	// in: int
	Status *int `binding:"omitempty"`
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
