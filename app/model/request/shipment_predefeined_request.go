package request

import (
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters ListShipmentPredefinedRequest
type ListShipmentPredefinedRequest struct {
	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields, example: name asc, uom desc
	// in: string
	Sort string `schema:"sort" binding:"omitempty"`

	// Type
	// in: string
	Type string `schema:"type" binding:"omitempty"`

	// Code
	// in: string
	Code string `schema:"code" binding:"omitempty"`

	// Title
	// in: string
	Title string `schema:"title" binding:"omitempty"`

	// Courier status
	// in: int
	Status *int `binding:"omitempty"`
}

// swagger:parameters UpdateShipmentPredefinedRequest
type ReqShipmentPredefinedRequest struct {
	// Uid of the shipmentPredefinedRequest
	// in: path
	// required: true
	UId string `json:"uid"`
	//  in: body
	Body UpdateShipmentPredefinedRequest `json:"body"`
}

type CreateShipmentPredefinedRequest struct {
	// Name of the ShippmentPredefined
	// in: string
	Type string `json:"type"`

	// Code of the ShippmentPredefined
	// in: string
	Code string `json:"code"`

	// Type of the Courier
	// in: string
	Title string `json:"title" binding:"omitempty"`

	// Note of the Courier
	// in: string
	Note string `json:"note" binding:"omitempty"`

	// Status of the ShippmentPredefined
	// in: integer
	Status int `json:"status" binding:"omitempty"`
}

// swaggers:
type UpdateShipmentPredefinedRequest struct {
	Uid string `json:"-" binding:"omitempty"`

	// Name of the ShippmentPredefined
	// in: string
	Type string `json:"type"`

	// Code of the ShippmentPredefined
	// in: string
	Code string `json:"code"`

	// Note of the ShippmentPredefined
	// in: string
	Note string `json:"note"`

	// Type of the Courier
	// in: string
	Title string `json:"title" binding:"omitempty"`
	// Status of the ShippmentPredefined
	// in: integer
	Status int `json:"status" binding:"omitempty"`
}

func (req UpdateShipmentPredefinedRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Type, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Code, validation.Required.Error(message.ErrReq.Message)),
	)
}
