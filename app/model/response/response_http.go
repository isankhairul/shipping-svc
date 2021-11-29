package response

// swagger:model SuccessResponse
type responseHttp struct {
	// in: int64
	Meta metaResponse `json:"meta"`
	// Pagination of the paginate respons
	// in: string
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	// in: string
	Data interface{} `json:"data"`
}

type metaResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

func SetResponse(code int, message string, reqId string, data interface{}, paging *PaginationResponse) interface{} {
	return responseHttp{
		Meta: metaResponse{
			Code:      code,
			Message:   message,
			RequestId: reqId,
		},
		Pagination: paging,
		Data:       data,
	}
}
