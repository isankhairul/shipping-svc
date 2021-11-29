package base

// swagger:model SuccessResponse
type responseHttp struct {
	// in: int64
	Meta metaResponse `json:"meta"`
	// Pagination of the paginate respons
	// in: string
	Pagination *Pagination `json:"pagination,omitempty"`
	// in: string
	Data interface{} `json:"data"`
}

type metaResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SetHttpResponse(code int, message string, data interface{}, paging *Pagination) interface{} {
	return responseHttp{
		Meta: metaResponse{
			Code:    code,
			Message: message,
		},
		Pagination: paging,
		Data:       data,
	}
}
