package base

import "reflect"

// swagger:model SuccessResponse
type responseHttp struct {
	// Meta is the API response information
	// in: MetaResponse
	Meta metaResponse `json:"meta"`
	// Pagination of the paginate respons
	// in: PaginationResponse
	Pagination *Pagination `json:"pagination,omitempty"`
	// Data is our data
	// in: DataResponse
	Data data `json:"data"`
	// Errors is the response message
	//in: string
	Errors interface{} `json:"errors,omitempty"`
}

// swagger:model MetaResponse
type metaResponse struct {
	// Code is the response code
	// example: 201000
	Code int `json:"code"`
	// Message is the response message
	// example: Success
	Message string `json:"message"`
}

// swagger:model DataResponse
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

// type ResponseFile struct {
// 	Name string
// 	Type string
// 	Data interface{}
// }
