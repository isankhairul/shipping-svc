package request

import (
	"encoding/json"
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters ListShipmentPredefinedRequest
type ListShipmentPredefinedRequest struct {
	// Filter : {"type":"value","code":"value","title":"value","status":1}
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

	Filters ListShipmentPredefinedFilter `json:"-"`
}

type ListShipmentPredefinedFilter struct {
	Type   string `json:"type" binding:"omitempty"`
	Code   string `json:"code" binding:"omitempty"`
	Title  string `json:"title" binding:"omitempty"`
	Status *int   `binding:"omitempty"`
}

func (m *ListShipmentPredefinedRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
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
	Status int32 `json:"status" binding:"omitempty"`
}

func (req UpdateShipmentPredefinedRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Type, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Code, validation.Required.Error(message.ErrReq.Message)),
	)
}

// swagger:parameters GetShipmentPredefinedByUID
type ShipmentPredefinedByUIDRequest struct {
	// in: path
	// required: true
	UID string `json:"uid"`
}
