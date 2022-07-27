package request

import (
	"bytes"
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters CourierCoverageCodeListRequest
type CourierCoverageCodeListRequest struct {
	// Maximum records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Sort fields
	// in: string
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`

	// Courier Name
	// in: string
	CourierName string `schema:"courier_name" binding:"omitempty" json:"courier_name"`

	// Country Code
	// in: string
	CountryCode string `schema:"country_code" binding:"omitempty" json:"country_code"`

	// Postal Code
	// in: string
	PostalCode string `schema:"postal_code" binding:"omitempty" json:"postal_code"`

	// Description
	// in: string
	Description string `schema:"description" binding:"omitempty" json:"description"`

	// Courier coverage code status
	// in: int
	Status *int `json:"status" binding:"omitempty"`
}

// swagger:parameters GetCourierCoverageCodeRequest DeleteCourierCoverageCodeRequest
type CourierCoverageCodeRequest struct {
	// Uid of the article
	// in: path
	Uid string `json:"uid" schema:"uid"`
}

// swagger:parameters UpdateCourierCoverageCodeBody
type ReqUpdateCourierCoverageCodeBody struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`

	// in: body
	Body SaveCourierCoverageCodeRequest `json:"body"`
}

// swagger:parameters SaveCourierCoverageCodeRequest
type ReqSaveCourierCoverageCodeBody struct {
	// in: body
	Body SaveCourierCoverageCodeRequest `json:"body"`
}

type SaveCourierCoverageCodeRequest struct {
	// Courier UID of the Courier
	// required: True
	// in: string
	CourierUID string `json:"courier_uid"`

	// Country code of the Courier Coverage Code
	// required: True
	// in: string
	CountryCode string `json:"country_code"`

	// Postal code of the Courier Coverage Code
	// required: True
	// in: string
	PostalCode string `json:"postal_code"`

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

	// Uid of the courỉe, use this on UPDATE function
	// in: int32
	Uid string `json:"-"`

	// Status of coverage code of the courỉe, use this on UPDATE function
	// in: int32
	Status int `json:"status"`
}

// swagger:parameters DeleteCourierCoverageCodeByIDParam
type DeleteCourierCoverageCodeByIDParam struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`
}

// swagger:parameters CourierCoverageCodeByIDParam
type CourierCoverageCodeByIDParam struct {
	// name: id
	// in: path
	// required: true
	UId string `json:"uid"`
}

// swagger:parameters ImportCourierCoverageCodeRequest
type ImportCourierCoverageCodeRequest struct {
	Rows []map[string]string `json:"-" binding:"omitempty"`
	// in: formData
	// name: file
	// swagger:file
	// require:true
	File *bytes.Buffer `json:"file"`

	FileName string `json:"-"`
}

func (req SaveCourierCoverageCodeRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CountryCode, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.PostalCode, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.CourierUID, validation.Required.Error(message.ErrReq.Message)),
	)
}
