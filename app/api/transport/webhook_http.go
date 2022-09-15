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
)

func WebhookHttpHandler(s service.ShippingService, logger log.Logger) http.Handler {
	pr := mux.NewRouter()

	ep := endpoint.MakeShippingEndpoint(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoder.EncodeError),
	}

	pr.Methods("POST").Path(fmt.Sprint(global.PrefixBase, global.PrefixWebhook, global.PathShipper)).Handler(httptransport.NewServer(
		ep.UpdateStatusShipper,
		decodeUpdateStatusShipper,
		encoder.EncodeResponseHTTP,
		options...,
	))

	return pr
}

func decodeUpdateStatusShipper(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
	var req request.WebhookUpdateStatusShipper
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}
