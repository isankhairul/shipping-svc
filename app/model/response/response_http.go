package response

type responseHttp struct {
	Meta       metaResponse        `json:"meta"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	Data       interface{}         `json:"data"`
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
