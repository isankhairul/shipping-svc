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

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func CourierHttpHandler(s service.CourierService, cc service.ChannelCourierService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeCourierEndpoints(s, cc)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourier)).Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourier)).Handler(httptransport.NewServer(
		ep.List,
		encoder.DecodePaginationRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierUID)).Handler(httptransport.NewServer(
		ep.Show,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierUID)).Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierUID)).Handler(httptransport.NewServer(
		ep.Delete,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	//courier-services
	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierService)).Handler(httptransport.NewServer(
		ep.SaveCourierSerivce,
		decodeSaveCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierService)).Handler(httptransport.NewServer(
		ep.ListCourierSerivce,
		decodeListCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierServiceUID)).Handler(httptransport.NewServer(
		ep.ShowCourierSerivce,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierServiceUID)).Handler(httptransport.NewServer(
		ep.UpdateCourierSerivce,
		decodeUpdateCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, global.PathCourierServiceUID)).Handler(httptransport.NewServer(
		ep.DeleteCourierSerivce,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourier, "shipping-type")).Handler(httptransport.NewServer(
		ep.ListShippingType,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))
	return pr
}

func decodeSaveCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	return req, nil
}

func decodeUpdateCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	req.Uid = mux.Vars(r)[pathUID]
	return req, nil
}

// CreateCourierService godoc
// @Summary API CreateCourierService
// @Description API CreateCourierService
// @Security AuthorizationHeader
// @Tags CourierService
// @Accept json
// @Param data body request.SaveCourierServiceRequest true "Request data"
// @Produce json
// @Success 200 {object}  entity.CourierService
// @Router /courier/courier-serivces [post]
func decodeSaveCourierService(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveCourierServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	return req, nil
}

// GetListCourierService godoc
// @Summary API GetListCourierService
// @Description API GetListCourierService
// @Security AuthorizationHeader
// @Tags CourierService
// @Accept json
// @Param data body request.CourierServiceListRequest true "Request data"
// @Produce json
// @Success 200 {object}  []entity.CourierService
// @Router /courier/courier-serivces [get]
func decodeListCourierService(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CourierServiceListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.GetFilter()

	return params, nil
}

// UpdateCourierService godoc
// @Summary API UpdateCourierService
// @Description API UpdateCourierService
// @Security AuthorizationHeader
// @Tags CourierService
// @Accept json
// @Param data body request.UpdateCourierServiceRequest true "Request data"
// @Produce json
// @Success 200 {object}  message.Message
// @Router /courier/courier-serivces/{uid} [put]
func decodeUpdateCourierService(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateCourierServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	req.Uid = mux.Vars(r)[pathUID]
	return req, nil
}
