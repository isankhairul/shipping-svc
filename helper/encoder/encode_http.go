package encoder

import (
	"context"
	"encoding/json"
	"go-klikdokter/helper/message"
	"net/http"
	"strconv"
)

type errorer interface {
	error() error
}

func EncodeResponseHTTP(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	if err, ok := resp.(errorer); ok && err.error() != nil {
		EncodeError(ctx, err.error(), w)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

//Encode error, for HTTP
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code, _ := strconv.Atoi(err.Error())
	switch code {
	case message.CODE_ERR_NOTFOUND, message.CODE_ERR_BADROUTING:
		w.WriteHeader(http.StatusNotFound)
	case message.CODE_ERR_NOAUTH:
		w.WriteHeader(http.StatusUnauthorized)
	case message.CODE_ERR_DB, message.CODE_ERR_BADREQUEST, message.CODE_ERR_VALIDATE:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
