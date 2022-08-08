package encoder

import (
	"context"
	"go-klikdokter/app/model/request"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func DecodeEmptyRequestHTTP(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	return nil, nil
}

func DecodePaginationRequestHTTP(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CourierListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.GetFilter()
	return params, nil

}

func UIDRequestHTTP(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["uid"]
	return uid, nil
}
