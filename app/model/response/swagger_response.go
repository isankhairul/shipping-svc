package response

// swagger:model MetaResponse
type MetaSingleSuccessResponse struct {
	// Code of Response
	// in: int
	// example: 1000
	Code int `json:"code"`

	// Message of Response
	// in: string
	// example: Success
	Message string `json:"message"`
}
