package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/pkg/util"
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

	pr.Methods("POST").Path(util.PrefixBase + "/courier/courier").Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(util.PrefixBase + "/courier/courier").Handler(httptransport.NewServer(
		ep.List,
		encoder.DecodePaginationRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(util.PrefixBase + "/courier/courier/{id}").Handler(httptransport.NewServer(
		ep.Show,
		decodeShowCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(util.PrefixBase + "/courier/courier/{id}").Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(util.PrefixBase + "/courier/courier/{id}").Handler(httptransport.NewServer(
		ep.Delete,
		decodeDeleteCourier,
		encoder.EncodeResponseHTTP,
		options...,
	))

	//courier-services
	pr.Methods("POST").Path(util.PrefixBase + "/courier/courier-services").Handler(httptransport.NewServer(
		ep.SaveCourierSerivce,
		decodeSaveCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(util.PrefixBase + "/courier/courier-services").Handler(httptransport.NewServer(
		ep.ListCourierSerivce,
		decodeListCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(util.PrefixBase + "/courier/courier-services/{id}").Handler(httptransport.NewServer(
		ep.ShowCourierSerivce,
		decodeShowCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(util.PrefixBase + "/courier/courier-services/{id}").Handler(httptransport.NewServer(
		ep.UpdateCourierSerivce,
		decodeUpdateCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(util.PrefixBase + "/courier/courier-services/{id}").Handler(httptransport.NewServer(
		ep.DeleteCourierSerivce,
		decodeDeleteCourierService,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(util.PrefixBase + "/courier/shipping-type").Handler(httptransport.NewServer(
		ep.ListShippingType,
		decodeShowCourierService,
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

func decodeShowCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}

func decodeUpdateCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.UpdateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	req.Uid = mux.Vars(r)["id"]
	return req, nil
}

func decodeDeleteCourier(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
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

// GetListCourierServiceById godoc
// @Summary API GetListCourierService
// @Description API GetListCourierService
// @Security AuthorizationHeader
// @Tags CourierService
// @Accept json
// @Param data body request.GetCourierServiceRequest true "Request data"
// @Produce json
// @Success 200 {object}  entity.CourierService
// @Router /courier/courier-serivces/{uid} [get]
func decodeShowCourierService(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
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

	req.Uid = mux.Vars(r)["id"]
	return req, nil
}

// DeleteCourierService godoc
// @Summary API DeleteCourierService
// @Description API DeleteCourierService
// @Security AuthorizationHeader
// @Tags CourierService
// @Accept json
// @Param data body request.GetCourierServiceRequest true "Request data"
// @Produce json
// @Success 200 {object}  message.Message
// @Router /courier/courier-serivces/{uid} [delete]
func decodeDeleteCourierService(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}
