package encoder

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/model/base"
	"go-klikdokter/helper/message"
	"net/http"
)

type errorer interface {
	error() error
}

type errorResponse struct {
	// Meta is the API response information
	// in: struct{}
	Meta struct {
		// Code is the response code
		//in: int
		Code int `json:"code"`
		// Message is the response message
		//in: string
		Message string `json:"message"`
	} `json:"meta"`
	// Data is our data
	// in: struct{}
	Data interface{} `json:"data"`
}

func EncodeResponseHTTP(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if err, ok := resp.(errorer); ok && err.error() != nil {
		EncodeError(ctx, err.error(), w)
		return nil
	}

	result := base.GetHttpResponse(resp)
	code := result.Meta.Code
	switch code {
	case message.CODE_ERR_NOTFOUND, message.CODE_ERR_BADROUTING:
		w.WriteHeader(http.StatusNotFound)
	case message.CODE_ERR_NOAUTH:
		w.WriteHeader(http.StatusUnauthorized)
	case message.CODE_ERR_DB, message.CODE_ERR_BADREQUEST, message.CODE_ERR_VALIDATE:
		w.WriteHeader(http.StatusBadRequest)
	case message.CODE_SUCCESS:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

//Encode error, for HTTP
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	result := &errorResponse{}
	result.Meta.Code = message.CODE_ERR_BADREQUEST
	result.Meta.Message = message.MSG_INVALID_REQUEST
	json.NewEncoder(w).Encode(result)
}
