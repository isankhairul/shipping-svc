package response

// swagger:response ImportStatus
type ImportStatus struct {
	// UID of the Courier Coverage Code
	// in: body
	UID string `json:"uid,omitempty"`

	// Courier UID of the Courier
	// in: body
	CourierUID string `json:"courier_uid"`

	// Country code of the Courier Coverage Code
	// in: body
	CountryCode string `json:"country_code" binding:"omitempty"`

	// Postal code of the Courier Coverage Code
	// in: body
	PostalCode string `json:"postal_code"`

	// Description of the Courier Coverage Code
	// in: body
	Description string `json:"description"`

	// Code 1 of the Courier Coverage Code
	// in: body
	Code1 string `json:"code1"`

	// Code 2 of the Courier Coverage Code
	// in: body
	Code2 string `json:"code2"`

	// Code 3 of the Courier Coverage Code
	// in: body
	Code3 string `json:"code3"`

	// Code 4 of the Courier Coverage Code
	// in: body
	Code4 string `json:"code4"`

	// Code 5 of the Courier Coverage Code
	// in: body
	Code5 string `json:"code5"`

	// Code 6 of the Courier Coverage Code
	// in: body
	Code6 string `json:"code6"`

	// Import Status
	// in: body
	Status bool `json:"status"`

	// Message Status
	// In: body
	Message string `json:"message"`
}
