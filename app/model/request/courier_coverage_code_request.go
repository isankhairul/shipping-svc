package request

import (
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters CourierCoverageCodeListRequest
type CourierCoverageCodeListRequest struct {
	// Maximun records per page
	// in: path
	Limit int `schema:"limit" binding:"omitempty,numeric,min=10,max=100"`

	// Page No
	// in: path
	Page int `schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields, example: name asc, uom desc
	// in: path
	Sort string `schema:"sort" binding:"omitempty"`
}

// swagger:parameters GetCourierCoverageCodeRequest DeleteCourierCoverageCodeRequest
type CourierCoverageCodeRequest struct {
	// Uid of the article
	// in: path
	Uid string `json:"uid" schema:"uid"`
}

// swagger:parameters SaveCourierCoverageCodeRequest
type ReqSaveCourierCoverageCodeBody struct {
	// in: body
	Body SaveCourierCoverageCodeRequest `json:"body"`
}

type SaveCourierCoverageCodeRequest struct {
	// Uid of the courỉe, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`

	// Courier UID of the Courier
	// required: True
	// in: string
	CourierUID string `json:"courier_uid" binding:"omitempty"`

	// Country code of the Courier Coverage Code
	// required: True
	// in: string
	CountryCode string `json:"country_code" binding:"omitempty"`

	// Postal code of the Courier Coverage Code
	// required: True
	// in: string
	PostalCode string `json:"postal_code" binding:"omitempty"`

	// Description of the Courier Coverage Code
	// in: string
	Description string `json:"description"`

	// Code 1 of the Courier Coverage Code
	// in: string
	Code1 string `json:"code1"`

	// Code 2 of the Courier Coverage Code
	// in: string
	Code2 string `json:"code2"`

	// Code 3 of the Courier Coverage Code
	// in: string
	Code3 string `json:"code3"`

	// Code 4 of the Courier Coverage Code
	// in: string
	Code4 string `json:"code4"`

	// Code 5 of the Courier Coverage Code
	// in: string
	Code5 string `json:"code5"`

	// Code 6 of the Courier Coverage Code
	// in: string
	Code6 string `json:"code6"`
}

// swagger:parameters CourierCoverageCodeByIDParam
type CourierCoverageCodeByIDParam struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`
}

func (req SaveCourierCoverageCodeRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CountryCode, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.PostalCode, validation.Required.Error(message.ErrReq.Message)),
	)
}
