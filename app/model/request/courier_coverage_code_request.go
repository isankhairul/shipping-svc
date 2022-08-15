package request

import (
	"bytes"
	"encoding/json"
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters CourierCoverageCodeListRequest
type CourierCoverageCodeListRequest struct {
	// Filter : {"courier_name":["value","value"],"country_code":["value","value"],"postal_code":["value","value"],"description":["value","value"]
	//,"status":[0,1], "code1":["value","value"], "code2":["value","value"], "code3":["value","value"], "code4":["value","value"], "code5":["value","value"], "code6":["value","value"]}
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
	Sort    string                        `schema:"sort" binding:"omitempty" json:"sort"`
	Filters CourierCoverageCodeListFilter `json:"-"`
}

type CourierCoverageCodeListFilter struct {
	CourierName []string `json:"courier_name"`
	CountryCode []string `json:"country_code"`
	PostalCode  []string `json:"postal_code"`
	Description []string `json:"description"`
	Status      []int    `json:"status"`
	Code1       []string `json:"code1"`
	Code2       []string `json:"code2"`
	Code3       []string `json:"code3"`
	Code4       []string `json:"code4"`
	Code5       []string `json:"code5"`
	Code6       []string `json:"code6"`
}

func (m *CourierCoverageCodeListRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
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
	Status int32 `json:"status"`
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
