package transport

import (
	"context"
	"encoding/json"
	"go-klikdokter/app/api/endpoint"
	"go-klikdokter/app/model/base/encoder"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func DoctorHttpHandler(s service.DoctorService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeDoctorEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		//httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path("/doctors/").Handler(httptransport.NewServer(
		ep.SaveDoctor,
		decodeSaveDoctor,
		encoder.EncodeResponseHTTP,
		options...,
	))

	pr.Methods("GET").Path("/doctors/{id}").Handler(httptransport.NewServer(
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
