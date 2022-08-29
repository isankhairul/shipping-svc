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

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourierCoverageCode)).Handler(httptransport.NewServer(
		ep.Save,
		decodeSaveCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourierCoverageCode)).Handler(httptransport.NewServer(
		ep.List,
		decodeListCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourierCoverageCode, global.PathUID)).Handler(httptransport.NewServer(
		ep.Show,
		encoder.UIDRequestHTTP,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("PUT").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourierCoverageCode, global.PathUID)).Handler(httptransport.NewServer(
		ep.Update,
		decodeUpdateCourierCoverageCode,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourierCoverageCode, global.PathImport)).Handler(httptransport.NewServer(
		ep.Import,
		decodeImportCourierCoverageCode,
		encoder.EncodeResponseCSV,
		options...,
	))

	pr.Methods("DELETE").Path(fmt.Sprint(global.PrefixBase, global.PrefixCourierCoverageCode, global.PathUID)).Handler(httptransport.NewServer(
		ep.Delete,
		encoder.UIDRequestHTTP,
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
	//global.HtmlEscape(&req)

	return req, nil
}

func decodeListCourierCoverageCode(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var params request.CourierCoverageCodeListRequest

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err = schema.NewDecoder().Decode(&params, r.Form); err != nil {
		return nil, err
	}

	params.GetFilter()

	return params, nil
}

func decodeUpdateCourierCoverageCode(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveCourierCoverageCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	//add this to htmlescape body post
	//global.HtmlEscape(&req)

	req.Uid = mux.Vars(r)[pathUID]
	return req, nil
}

func decodeImportCourierCoverageCode(ctx context.Context, r *http.Request) (rsqt interface{}, err error) {
	var req request.ImportCourierCoverageCodeRequest
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	// Close the file at the end of the program
	defer file.Close()
	rows, err := util.ReadCsvFile(file)

	if err != nil {
		return nil, err
	}

	limit := 5000
	if len(rows) > limit {
		return nil, fmt.Errorf("The number of rows in your dataset is greater than the maximum allowed (%d)", limit)
	}

	req.Rows = rows
	req.FileName = header.Filename
	return req, nil
}
