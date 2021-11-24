package entity

// swagger:model SuccessResponse
type ResponseHttp struct {
	// Code of the business code
	// in: int64
	Code int `json:"code"`
	// Message of the respons
	// in: int64
	Message string `json:"message"`
	// Message of the respons
	// in: string
	Data interface{} `json:"data"`
}

// swagger:model ErrorResponse
type ErrorResponseHttp struct {
	// Code of the business code
	// in: int64
	Code int `json:"code"`
	// Message of the respons
	// in: int64
	Message string `json:"message"`
	// Message of the respons
	// in: string
	Data interface{} `json:"data"`
}
