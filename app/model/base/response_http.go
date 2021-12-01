package base

import "reflect"

// swagger:model SuccessResponse
type responseHttp struct {
	// Meta is the API response information
	// in: struct{}
	Meta metaResponse `json:"meta"`
	// Pagination of the paginate respons
	// in: struct{}
	Pagination *Pagination `json:"pagination,omitempty"`
	// Data is our data
	// in: struct{}
	Data data `json:"data"`
}

type metaResponse struct {
	// Code is the response code
	//in: int
	Code int `json:"code"`
	// Message is the response message
	//in: string
	Message string `json:"message"`
}

type data struct {
	Records interface{} `json:"records,omitempty"`
	Record  interface{} `json:"record,omitempty"`
}

func SetHttpResponse(code int, message string, result interface{}, paging *Pagination) interface{} {
	dt := data{}
	isSlice := reflect.ValueOf(result).Kind() == reflect.Slice
	if isSlice {
		dt.Records = result
		dt.Record = nil
	} else {
		dt.Records = nil
		dt.Record = result
	}

	return responseHttp{
		Meta: metaResponse{
			Code:    code,
			Message: message,
		},
		Pagination: paging,
		Data:       dt,
	}
}

func GetHttpResponse(resp interface{}) *responseHttp {
	result, ok := resp.(responseHttp)

	if ok {
		return &result
	}
	return nil
}
