package shipping_provider

import (
	"encoding/json"
	"errors"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"strconv"

	"go-klikdokter/helper/http_helper"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

type Shipper interface {
	GetShippingRate(courierID *uint64, input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error)
	GetPricingDomestic(req *request.GetPricingDomestic) (*response.GetPricingDomestic, error)
}
type shipper struct {
	courierCoverage repository.CourierCoverageCodeRepository
	Logger          log.Logger
	Authorization   map[string]string
	Base            string
}

func NewShipper(ccr repository.CourierCoverageCodeRepository, log log.Logger) Shipper {
	return &shipper{
		Authorization: map[string]string{
			viper.GetString("shipper.auth.key"): viper.GetString("shipper.auth.value"),
		},
		Base:            viper.GetString("shipper.base"),
		courierCoverage: ccr,
		Logger:          log,
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

func (h *shipper) GetOriginAndDestination(courierID *uint64, input *request.GetShippingRateRequest) (int, int, message.Message) {
	originReq := &request.FindShipperCourierCoverage{
		CourierID:   *courierID,
		CountryCode: input.Origin.CountryCode,
		PostalCode:  input.Origin.PostalCode,
		Subdistrict: input.Origin.Subdistrict,
	}

	origin, _ := h.courierCoverage.FindShipperCourierCoverage(originReq)
	if origin == nil {
		return 0, 0, message.ErrOriginNotFound
	}

	originAreaID, _ := strconv.Atoi(origin.Subdistrict)
	if originAreaID == 0 {
		return 0, 0, message.ErrOriginNotFound
	}

	destinationReq := &request.FindShipperCourierCoverage{
		CourierID:   *courierID,
		CountryCode: input.Destination.CountryCode,
		PostalCode:  input.Destination.PostalCode,
		Subdistrict: input.Destination.Subdistrict,
	}

	destination, _ := h.courierCoverage.FindShipperCourierCoverage(destinationReq)
	if destination == nil {
		return 0, 0, message.ErrDestinationNotFound
	}

	destinationAreaID, _ := strconv.Atoi(destination.Subdistrict)
	if destinationAreaID == 0 {
		return 0, 0, message.ErrDestinationNotFound
	}

	return originAreaID, destinationAreaID, message.SuccessMsg
}

func (h *shipper) GetShippingRate(courierID *uint64, input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error) {

	origin, destination, msg := h.GetOriginAndDestination(courierID, input)

	//if origin or destination not found
	if msg.Code != message.SuccessMsg.Code {
		return &response.ShippingRateCommonResponse{
			Rate: make(map[string]response.ShippingRateData),
			Msg:  msg,
		}, errors.New(msg.Message)
	}

	payload := request.NewGetPricingDomesticRequest(origin, destination, input)

	shipperResponse, err := h.GetPricingDomestic(payload)

	//if failed to get pricing from shipper api
	if err != nil {
		return &response.ShippingRateCommonResponse{
			Rate: make(map[string]response.ShippingRateData),
		}, err
	}

	//if everithing go well
	resp := shipperResponse.ToShippingRate()
	return resp, nil
}
