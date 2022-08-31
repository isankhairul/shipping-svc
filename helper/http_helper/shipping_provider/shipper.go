package shipping_provider

import (
	"encoding/json"
	"errors"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"

	"go-klikdokter/helper/http_helper"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

type Shipper interface {
	GetShippingRate(origin, destination *entity.CourierCoverageCode, data *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error)
	GetPricingDomestic(req *request.GetPricingDomestic) (*response.GetPricingDomestic, error)
}
type shipper struct {
	Logger        log.Logger
	Authorization map[string]string
	Base          string
}

func NewShipper(log log.Logger) Shipper {
	return &shipper{
		Authorization: map[string]string{
			viper.GetString("shipper.auth.key"): viper.GetString("shipper.auth.value"),
		},
		Base:   viper.GetString("shipper.base"),
		Logger: log,
	}
}

func (h *shipper) GetPricingDomestic(req *request.GetPricingDomestic) (*response.GetPricingDomestic, error) {

	response := response.GetPricingDomestic{}
	path := viper.GetString("shipper.path.get-pricing-domestic")
	url := h.Base + path

	header := h.Authorization
	header["Content-Type"] = "application/json"

	respByte, err := http_helper.Post(url, header, req, h.Logger)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respByte, &response)

	if err != nil {
		return nil, err
	}

	if response.Metadata.HTTPStatusCode != 200 {
		return nil, errors.New(response.Metadata.HTTPStatus)
	}

	return &response, nil
}

func (h *shipper) GetShippingRate(origin, destination *entity.CourierCoverageCode, data *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error) {

	payload := request.NewGetPricingDomesticRequest(origin, destination, data)

	shipperResponse, err := h.GetPricingDomestic(payload)
	if err != nil {
		return &response.ShippingRateCommonResponse{
			Rate: make(map[string]response.ShippingRateData),
			Msg:  message.ErrGetShipperRate,
		}, err
	}

	resp := shipperResponse.ToShippingRate()
	return resp, nil
}
