package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func ShipmentPredefinedHandler(s service.ShipmentPredefinedService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeShipmentPredefinedEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixOther, global.PathShipmentPredefinedUID)).Handler(httptransport.NewServer(
		ep.Show,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(fmt.Sprint(global.PrefixBase, global.PrefixOther, global.PathShipmentPredefinedUID)).Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateShipmentPredefinedRequest,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixOther, global.PathShipmentPredefined)).Handler(httptransport.NewServer(
		ep.List,
		decodeListShipmentPredefinedRequest,
		encoder.EncodeResponseHTTP,
		options...,
	))
	return pr
}

func decodeUpdateShipmentPredefinedRequest(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateShipmentPredefinedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	req.Uid = mux.Vars(r)[pathUID]
	return req, nil
}

func decodeListShipmentPredefinedRequest(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.ListShipmentPredefinedRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.GetFilter()

	return params, nil
}
