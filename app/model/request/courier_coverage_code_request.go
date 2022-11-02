package request

import (
	"bytes"
	"encoding/json"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters CourierCoverageCodeListRequest
type CourierCoverageCodeListRequest struct {
	// Filter : {"courier_name":["value","value"],"country_code":["value","value"],"postal_code":["value","value"],"description":["value","value"]
	//,"district_numeric_code":["value", "value"],"district_name":["value","value"],"subdistrict":["value","value"],"subdistrict_name":["value","value"]
	//,"province_numeric_code":["value","value"],"province_name":["value","value"],"city_numeric_code":["value","value"],"city_name":["value","value"]
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
	// enum:courier_name,courier_name desc,country_code,country_code desc,postal_code,postal_code desc,district_numeric_code,district_numeric_code desc,district,district desc,subdistrict_name,subdistrict_name desc,subdistrict,subdistrict desc,province_numeric_code,province_numeric_code desc,province_name,province_name desc,city_numeric_code,city_numeric_code desc,city_name,city_name desc,description,description desc,status,status desc,code1,code1 desc,code2,code2 desc,code3,code3 desc,code4,code4 desc,code5,code5 desc,code6,code6 desc
	Sort    string                        `schema:"sort" binding:"omitempty" json:"sort"`
	Filters CourierCoverageCodeListFilter `json:"-"`
}

type CourierCoverageCodeListFilter struct {
	CourierName         []string `json:"courier_name"`
	CountryCode         []string `json:"country_code"`
	ProvinceNumericCode []string `json:"province_numeric_code"`
	ProvinceName        []string `json:"province_name"`
	CityNumericCode     []string `json:"city_numeric_code"`
	CityName            []string `json:"city_name"`
	PostalCode          []string `json:"postal_code"`
	DistrictNumericCode []string `json:"districtNumericCode"`
	DistrictName        []string `json:"districtName"`
	Subdistrict         []string `json:"subdistrict"`
	SubdistrictName     []string `json:"subdistrict_name"`
	Description         []string `json:"description"`
	Status              []int    `json:"status"`
	Code1               []string `json:"code1"`
	Code2               []string `json:"code2"`
	Code3               []string `json:"code3"`
	Code4               []string `json:"code4"`
	Code5               []string `json:"code5"`
	Code6               []string `json:"code6"`
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
	// required: true
	// in: string
	CourierUID string `json:"courier_uid"`

	// Country code of the Courier Coverage Code
	// required: true
	// in: string
	CountryCode string `json:"country_code"`

	// Province Numeric Code of the Courier Coverage Code
	// required: true
	// in: string
	// example: 99
	ProvinceNumericCode string `json:"province_numeric_code"`

	// Province Name of the Courier Coverage Code
	// required: true
	// in: string
	// example: DKI Jakarta
	ProvinceName string `json:"province_name"`

	// City Numeric Code of the Courier Coverage Code
	// required: true
	// in: string
	// example: 99
	CityNumericCode string `json:"city_numeric_code"`

	// City Name of the Courier Coverage Code
	// required: true
	// in: string
	// example: Jakarta Selatan
	CityName string `json:"city_name"`

	// Postal code of the Courier Coverage Code
	// required: true
	// in: string
	PostalCode string `json:"postal_code"`

	// District numeric code of the Courier Coverage Code
	// in: string
	// require: true
	// example: 151338
	DistrictNumericCode string `json:"district_numeric_code"`

	// District name of the Courier Coverage Code
	// in: string
	// require: true
	// example: "Denpasar Barat"
	DistrictName string `json:"district_name"`

	// Subdistrict numeric code of the Courier Coverage Code
	// required: true
	// in: string
	// example: 447
	Subdistrict string `json:"subdistrict"`

	// Subdistrict name of the Courier Coverage Code
	// in: string
	// require: true
	// example: "Tegal Harum"
	SubdistrictName string `json:"subdistrict_name"`

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

	// Extend Jwt Info
	global.JWTInfo
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

	// Extend Jwt Info
	global.JWTInfo
}

func (req SaveCourierCoverageCodeRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CountryCode, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.PostalCode, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.CourierUID, validation.Required.Error(message.ErrReq.Message)),
	)
}
