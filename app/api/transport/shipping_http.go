package transport

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"net/http"
	"time"

	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	shippingTypePath = "shipping-type"
	channelUID       = "channel-uid"
)

func ShippingHttpHandler(s service.ShippingService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeShippingEndpoint(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathShippingRate)).Handler(httptransport.NewServer(
		ep.GetShippingRate,
		decodeGetShippingRate,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathShippingRateShippingType)).Handler(httptransport.NewServer(
		ep.GetShippingRateByShippingType,
		decodeGetShippingRate,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShipping)).Handler(httptransport.NewServer(
		ep.CreateDelivery,
		decodeCreateDelivery,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderTracking)).Handler(httptransport.NewServer(
		ep.GetOrderShippingTracking,
		decodeGetOrderTracking,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShipping)).Handler(httptransport.NewServer(
		ep.GetOrderShippingList,
		decodeGetOrderShippingList,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShippingDownload)).Handler(httptransport.NewServer(
		ep.DownloadOrderShipping,
		decodeOrderShippingDownload,
		encodeOrderShippingDownload,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShippingUID)).Handler(httptransport.NewServer(
		ep.GetOrderShippingDetail,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathCancelPickupUID)).Handler(httptransport.NewServer(
		ep.CancelPickUp,
		decodeCancelPickup,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathCancelOrderUID)).Handler(httptransport.NewServer(
		ep.CancelOrder,
		decodeCancelOrder,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathOrderShippingLabel)).Handler(httptransport.NewServer(
		ep.GetOrderShippingLabel,
		decodeOrderShippingLabel,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathRepickup)).Handler(httptransport.NewServer(
		ep.RepickupOrder,
		decodeRepickupOrder,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixShipping, global.PathShippingTracking)).Handler(httptransport.NewServer(
		ep.GetShippingTracking,
		decodeGetOrderTracking,
		encoder.EncodeResponseHTTP,
		options...,
	))
	return pr
}

func decodeGetShippingRate(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.GetShippingRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.ShippingType = mux.Vars(r)[shippingTypePath]
	return req, nil
}

func decodeCreateDelivery(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.CreateDelivery
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetOrderTracking(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetOrderShippingTracking
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.UID = mux.Vars(r)[pathUID]
	return params, nil
}

func decodeGetOrderShippingList(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetOrderShippingList
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.GetFilter()
	return params, nil
}

func decodeOrderShippingDownload(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.DownloadOrderShipping
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}
	params.GetFilter()
	return params, nil
}

func decodeCancelOrder(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CancelOrder
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params.Body); err != nil {
		return nil, err
	}

	params.UID = mux.Vars(r)[pathUID]
	return params, nil
}

func decodeCancelPickup(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CancelPickup
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params.Body); err != nil {
		return nil, err
	}

	params.UID = mux.Vars(r)[pathUID]
	return params, nil
}

func decodeOrderShippingLabel(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.GetOrderShippingLabel
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params.Body); err != nil {
		return nil, err
	}

	params.ChannelUID = mux.Vars(r)[channelUID]
	return params, nil
}

func decodeRepickupOrder(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.RepickupOrderRequest
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, err
	}
	return params, nil
}

func encodeOrderShippingDownload(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	httpResponse := base.GetHttpResponse(resp)
	code := httpResponse.Meta.Code

	if code != message.SuccessMsg.Code {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		switch code {
		case message.ErrPageNotFound.Code, message.ErrBadRouting.Code:
			w.WriteHeader(http.StatusNotFound)
		case message.ErrNoAuth.Code:
			w.WriteHeader(http.StatusUnauthorized)
		case message.ErrDB.Code, message.ErrBadRouting.Code, message.ErrReq.Code:
			w.WriteHeader(http.StatusBadRequest)
		case message.SuccessMsg.Code, message.ShippingProviderMsg.Code:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return json.NewEncoder(w).Encode(resp)
	}

	orderShippings := httpResponse.Data.Records.([]response.DownloadOrderShipping)
	csvRows := mapToCsvArrays(orderShippings)

	fileName := "order-shipping-" + time.Now().In(util.Loc).Format(time.RFC3339)

	csvWriter := csv.NewWriter(w)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	if err := csvWriter.WriteAll(csvRows); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func mapToCsvArrays(orderShippings []response.DownloadOrderShipping) [][]string {
	csvRows := make([][]string, len(orderShippings)+1)
	csvRows[0] = []string{
		"channel", "order_shipping_date", "order_shipping_uid", "order_no", "courier_name", "courier_service",
		"airwaybill", "booking_id", "customer_name", "customer_phone_number", "customer_email", "customer_address",
		"customer_province_name", "customer_city_name", "customer_district_name", "customer_subdistrict",
		"customer_postal_code", "customer_notes", "merchant_name", "merchant_phone_number", "merchant_email",
		"merchant_address", "merchant_province_name", "merchant_city_name", "merchant_district_name",
		"merchant_subdistrict", "merchant_postal_code", "total_weight", "total_volume", "total_product_price",
		"total_final_weight", "contain_prescription", "insurance", "insurance_cost", "shipping_cost",
		"total_shipping_cost", "actual_shipping_cost", "shipping_notes", "shipping_status_name", "order_status_history",
	}

	for i, os := range orderShippings {
		csvRows[i+1] = []string{
			os.Channel,
			os.OrderShippingDate.In(util.Loc).Format(util.LayoutDefault),
			os.OrderShippingUid,
			os.OrderNo,
			os.CourierName,
			os.CourierService,
			os.Airwaybill,
			os.BookingId,
			os.CustomerName,
			os.CustomerPhoneNumber,
			os.CustomerEmail,
			os.CustomerAddress,
			os.CustomerProvinceName,
			os.CustomerCityName,
			os.CustomerDistrictName,
			os.CustomerSubdistrict,
			os.CustomerPostalCode,
			os.CustomerNotes,
			os.MerchantName,
			os.MerchantPhoneNumber,
			os.MerchantEmail,
			os.MerchantAddress,
			os.MerchantProvinceName,
			os.MerchantCityName,
			os.MerchantDistrictName,
			os.MerchantSubdistrict,
			os.MerchantPostalCode,
			os.TotalWeight,
			os.TotalVolume,
			os.TotalProductPrice,
			os.TotalFinalWeight,
			os.ContainPrescription,
			os.Insurance,
			os.InsuranceCost,
			os.ShippingCost,
			os.TotalShippingCost,
			os.ActualShippingCost,
			os.ShippingNotes,
			os.ShippingStatusName,
			os.OrderStatusHistory,
		}
	}

	return csvRows
}
