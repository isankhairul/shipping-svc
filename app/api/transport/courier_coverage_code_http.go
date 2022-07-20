package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/global"
	"go-klikdokter/pkg/util"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func CourierCoverageCodeHttpHandler(s service.CourierCoverageCodeService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeCourierCoverageCodeEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path(util.PrefixBase + "/courier/courier-coverage-code/").Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(util.PrefixBase + "/courier/courier-coverage-code/").Handler(httptransport.NewServer(
		ep.List,
		decodeListCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(util.PrefixBase + "/courier/courier-coverage-code/{id}").Handler(httptransport.NewServer(
		ep.Show,
		decodeShowCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(util.PrefixBase + "/courier/courier-coverage-code/{id}").Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(util.PrefixBase + "/courier/courier-coverage-code/import").Handler(httptransport.NewServer(
		ep.Import,
		decodeImportCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("DELETE").Path(util.PrefixBase + "/courier/courier-coverage-code/{id}").Handler(httptransport.NewServer(
		ep.Delete,
		decodeShowCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeSaveCourierCoverageCode(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveCourierCoverageCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	global.HtmlEscape(&req)

	return req, nil
}

func decodeShowCourierCoverageCode(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}

func decodeListCourierCoverageCode(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CourierCoverageCodeListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func decodeUpdateCourierCoverageCode(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveCourierCoverageCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	global.HtmlEscape(&req)

	req.Uid = mux.Vars(r)["id"]
	return req, nil
}

func decodeImportCourierCoverageCode(ctx context.Context, r *http.Request) (rsqt interface{}, err error) {
	var req request.ImportCourierCoverageCodeRequest
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	// Close the file at the end of the program
	defer file.Close()
	rows, err := util.ReadCsvFile(file)
	if err != nil {
		return nil, err
	}
	req.Rows = rows
	return req, nil
}
