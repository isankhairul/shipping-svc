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
}

// swagger:parameters courier
type GetCourierRequest struct {
	// name: id
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters list courier
type CourierListRequest struct {
	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields, example: name asc, uom desc
	// in: string
	Sort string `schema:"sort" binding:"omitempty"`

	// Courier type
	// in: string
	CourierType string `schema:"courier_type" binding:"omitempty"`

	// Courier status
	// in: string
	Status string `schema:"status" binding:"omitempty"`
}

// swagger:parameters courier
type UpdateCourierRequest struct {
	// Uid of the courier, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`

	// Name of the courier
	// in: string
	CourierName string `json:"courier_name"`

	// type of courier
	// in: string
	CourierType string `json:"courier_type"`

	// Logo of courier
	// in: string
	Logo string `json:"logo"`

	// Courier status
	// in: string
	Status string `json:"status" binding:"omitempty"`
}

func (req SaveCourierRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CourierName, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Code, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.CourierType, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Logo, validation.Required.Error(message.ErrReq.Message)),
	)
}
