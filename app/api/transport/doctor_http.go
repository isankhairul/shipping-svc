package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"go-klikdokter/helper/encoder"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func DoctorHttpHandler(s service.DoctorService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeDoctorEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path("/kd/v2/doctor").Handler(httptransport.NewServer(
		ep.SaveDoctor,
		decodeSaveDoctor,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/kd/v2/doctor/{id}").Handler(httptransport.NewServer(
		ep.Show,
		decodeShowProduct,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeSaveDoctor(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.SaveDoctorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeShowDoctor(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	uid := mux.Vars(r)["id"]
	return uid, nil
}
