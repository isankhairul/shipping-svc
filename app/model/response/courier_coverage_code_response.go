package response

// swagger:model ImportStatus
type ImportStatus struct {
	// UID of the Courier Coverage Code
	UID string `json:"uid,omitempty"`

	// Courier UID of the Courier
	CourierUID string `json:"courier_uid"`

	// Country code of the Courier Coverage Code
	CountryCode string `json:"country_code"`

	// Province Numeric Code of the Courier Coverage Code
	ProvinceNumericCode string `json:"province_numeric_code"`

	// Province Name of the Courier Coverage Code
	ProvinceName string `json:"province_name"`

	// City Numeric Code of the Courier Coverage Code
	CityNumericCode string `json:"city_numeric_code"`

	// City Name of the Courier Coverage Code
	CityName string `json:"city_name"`

	// Postal code of the Courier Coverage Code
	PostalCode string `json:"postal_code"`

	// Subdistrict of the Courier Coverage Code
	Subdistrict string `json:"subdistrict"`

	// Description of the Courier Coverage Code
	Description string `json:"description"`

	// Code 1 of the Courier Coverage Code
	Code1 string `json:"code1"`

	// Code 2 of the Courier Coverage Code
	Code2 string `json:"code2"`

	// Code 3 of the Courier Coverage Code
	Code3 string `json:"code3"`

	// Code 4 of the Courier Coverage Code
	Code4 string `json:"code4"`

	// Code 5 of the Courier Coverage Code
	Code5 string `json:"code5"`

	// Code 6 of the Courier Coverage Code
	Code6 string `json:"code6"`

	// Import Status
	Status bool `json:"-"`

	// Message Status

	Message string `json:"message"`
}

// swagger:model ImportSummary
type ImportSummary struct {
	//example: 100
	TotalRow int `json:"total_row"`
	//example: 55
	SuccessRow int `json:"success_row"`
	//example: 45
	FailedRow int `json:"failed_row"`
}

// swagger:response ImportCourierCoverageCode
type CourierCoverageCodeImportResponseBody struct {
	//in: body
	Response CourierCoverageCodeImportResponse `json:"response"`
}

// swagger:model ImportCourierCoverageCode
type CourierCoverageCodeImportResponse struct {
	FailedData []ImportStatus `json:"failed_data"`
	Summary    ImportSummary  `json:"summary"`
}
