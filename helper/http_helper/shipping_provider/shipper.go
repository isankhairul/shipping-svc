package shipping_provider

import (
	"encoding/json"
	"errors"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"strconv"
	"strings"

	"go-klikdokter/helper/http_helper"
	"go-klikdokter/helper/message"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
)

type Shipper interface {
	GetShippingRate(courierID *uint64, input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error)
	GetPricingDomestic(req *request.GetPricingDomestic) (*response.GetPricingDomestic, error)
	CreateOrder(req *request.CreateOrderShipper) (*response.CreateOrderShipperResponse, error)
	CreatePickUpOrder(req *request.CreatePickUpOrderShipper) (*response.CreatePickUpOrderShipperResponse, error)
	CreateDelivery(ShipperOrderID string, courierService *entity.CourierService, req *request.CreateDelivery) (*response.CreateDeliveryThirdPartyData, message.Message)
	GetOrderDetail(orderID string) (*response.GetOrderDetailResponse, error)
	GetTracking(orderID string) ([]response.GetOrderShippingTracking, message.Message)
	CancelPickupRequest(pickupCode string) (*response.MetadataResponse, error)
	CancelOrder(orderID string, req *request.CancelOrder) (*response.MetadataResponse, error)
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
		return nil, errors.New(response.Metadata.Errors[0].Message)
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
		return 0, 0, message.OriginNotFoundMsg
	}

	originAreaID, _ := strconv.Atoi(origin.Code1)
	if originAreaID == 0 {
		return 0, 0, message.OriginNotFoundMsg
	}

	destinationReq := &request.FindShipperCourierCoverage{
		CourierID:   *courierID,
		CountryCode: input.Destination.CountryCode,
		PostalCode:  input.Destination.PostalCode,
		Subdistrict: input.Destination.Subdistrict,
	}

	destination, _ := h.courierCoverage.FindShipperCourierCoverage(destinationReq)
	if destination == nil {
		return 0, 0, message.DestinationNotFoundMsg
	}

	destinationAreaID, _ := strconv.Atoi(destination.Code1)
	if destinationAreaID == 0 {
		return 0, 0, message.DestinationNotFoundMsg
	}

	return originAreaID, destinationAreaID, message.SuccessMsg
}

func (h *shipper) GetShippingRate(courierID *uint64, input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error) {

	origin, destination, msg := h.GetOriginAndDestination(courierID, input)

	//if origin or destination not found
	if msg.Code != message.SuccessMsg.Code {
		return &response.ShippingRateCommonResponse{
			Rate:       make(map[string]response.ShippingRateData),
			CourierMsg: map[string]message.Message{ShipperCode: msg},
		}, errors.New(msg.Message)
	}

	payload := request.NewGetPricingDomesticRequest(origin, destination, input)

	shipperResponse, err := h.GetPricingDomestic(payload)

	//if failed to get pricing from shipper api
	if err != nil {
		msg = message.ShippingProviderMsg
		msg.Message = err.Error()
		return &response.ShippingRateCommonResponse{
			Rate:       make(map[string]response.ShippingRateData),
			CourierMsg: map[string]message.Message{ShipperCode: msg},
		}, err
	}

	//if everithing go well
	resp := shipperResponse.ToShippingRate()
	return resp, nil
}

func (h *shipper) CreateOrder(req *request.CreateOrderShipper) (*response.CreateOrderShipperResponse, error) {

	response := response.CreateOrderShipperResponse{}
	path := viper.GetString("shipper.path.order")
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

	if response.Metadata.HTTPStatusCode != 200 && response.Metadata.HTTPStatusCode != 201 {
		return nil, errors.New(response.Metadata.Errors[0].Message)
	}
	return &response, nil
}

func (h *shipper) GetTimeslot(req *request.GetPickUpTimeslot) (*response.GetPickUpTimeslotResponse, error) {
	response := response.GetPickUpTimeslotResponse{}
	path := viper.GetString("shipper.path.pick-up-timeslot")
	url := h.Base + path

	header := h.Authorization
	header["Content-Type"] = "application/json"

	params := map[string]string{
		"time_zone": req.TimeZone,
	}

	respByte, err := http_helper.Get(url, header, params, h.Logger)

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

func (h *shipper) CreatePickUpOrder(req *request.CreatePickUpOrderShipper) (*response.CreatePickUpOrderShipperResponse, error) {

	response := response.CreatePickUpOrderShipperResponse{}
	path := viper.GetString("shipper.path.pick-up-timeslot")
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

func (h *shipper) CreatePickUpOrderWithTimeSlots(orderID ...string) (*response.CreatePickUpOrderShipperResponse, message.Message) {
	logger := log.With(h.Logger, "Shipper", "CreatePickUpOrderWithTimeSlots")
	timeslots, err := h.GetTimeslot(&request.GetPickUpTimeslot{TimeZone: "Asia/Jakarta"})

	if err != nil {
		_ = level.Error(logger).Log("h.CreateOrder", err.Error())
		return nil, message.ErrGetPickUpTimeslot
	}

	req := &request.CreatePickUpOrderShipper{
		Data: request.CreatePickUpOrderShipperData{
			OrderActivation: request.CreatePickUpOrderShipperOrderActivation{
				OrderID:   orderID,
				Timezone:  timeslots.Data.Timezone,
				StartTime: timeslots.Data.Timeslots[0].StartTime,
				EndTime:   timeslots.Data.Timeslots[0].EndTime,
			},
		},
	}

	pickup, err := h.CreatePickUpOrder(req)
	if err != nil {
		_ = level.Error(logger).Log("h.CreatePickUpOrder", err.Error())
		return nil, message.ErrCreatePickUpOrder
	}

	return pickup, message.SuccessMsg
}

func (h *shipper) CreateDelivery(shipperOrderID string, courierService *entity.CourierService, req *request.CreateDelivery) (*response.CreateDeliveryThirdPartyData, message.Message) {
	logger := log.With(h.Logger, "Shipper", "CreateDelivery")
	resp := &response.CreateDeliveryThirdPartyData{}
	order := &response.CreateOrderShipperResponse{Data: response.CreateOrderShipper{OrderID: shipperOrderID}}
	var err error
	if len(shipperOrderID) == 0 {
		input := request.GetShippingRateRequest{
			Origin: request.AreaDetailPayload{
				CountryCode: req.Origin.CountryCode,
				PostalCode:  req.Origin.PostalCode,
				Subdistrict: req.Origin.Subdistrict,
			},
			Destination: request.AreaDetailPayload{
				CountryCode: req.Destination.CountryCode,
				PostalCode:  req.Destination.PostalCode,
				Subdistrict: req.Destination.Subdistrict,
			},
		}
		origin, destination, msg := h.GetOriginAndDestination(&courierService.CourierID, &input)

		if msg != message.SuccessMsg {
			return nil, msg
		}

		orderRequet := req.ToCreateOrderShipper()
		orderRequet.Origin.AreaID = uint64(origin)
		orderRequet.Destination.AreaID = uint64(destination)
		orderRequet.Courier.RateID, _ = strconv.Atoi(courierService.ShippingCode)
		orderRequet.Package.PackageType = viper.GetInt("shipper.setting.package-type")
		order, err = h.CreateOrder(orderRequet)

		if err != nil {
			_ = level.Error(logger).Log("h.CreateOrder", err.Error())
			msg = message.ShippingProviderMsg
			msg.Message = err.Error()
			return nil, msg
		}

		resp.Insurance = order.Data.Courier.UseInsurance
		resp.InsuranceCost = order.Data.Courier.InsuranceAmount
		resp.ShippingCost = order.Data.Courier.Amount
		resp.TotalShippingCost = order.Data.Courier.Amount
		resp.ActualShippingCost = order.Data.Courier.Amount
		resp.BookingID = order.Data.OrderID
		resp.Status = StatusCreated
	}

	pickup, msg := h.CreatePickUpOrderWithTimeSlots(order.Data.OrderID)

	if msg == message.SuccessMsg {
		resp.Status = StatusRequestPickup
		resp.PickUpCode = pickup.Data.OrderActivation[0].PickUpCode
		resp.PickUpTime = pickup.Data.OrderActivation[0].PickUpTime
	}

	return resp, message.SuccessMsg
}

func (h *shipper) GetOrderDetail(orderID string) (*response.GetOrderDetailResponse, error) {
	response := response.GetOrderDetailResponse{}
	path := viper.GetString("shipper.path.order-detail")
	path = strings.ReplaceAll(path, "{orderID}", orderID)
	url := h.Base + path

	header := h.Authorization
	header["Content-Type"] = "application/json"

	respByte, err := http_helper.Get(url, header, map[string]string{}, h.Logger)

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

func (h *shipper) GetTracking(orderID string) ([]response.GetOrderShippingTracking, message.Message) {
	logger := log.With(h.Logger, "Shipper", "GetTracking")
	orderDetail, err := h.GetOrderDetail(orderID)
	if err != nil {
		_ = level.Error(logger).Log("h.GetOrderDetail", err.Error())
		return nil, message.ErrGetOrderDetail
	}

	return orderDetail.Data.ToOrderShippingTracking(), message.SuccessMsg
}

func (h *shipper) CancelPickupRequest(pickupCode string) (*response.MetadataResponse, error) {
	response := response.MetadataResponse{}
	path := viper.GetString("shipper.path.cancel-pickup")
	url := h.Base + path

	header := h.Authorization
	header["Content-Type"] = "application/json"

	respByte, err := http_helper.Patch(url, header, map[string]string{"pickup_Code": pickupCode}, h.Logger)

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

func (h *shipper) CancelOrder(orderID string, req *request.CancelOrder) (*response.MetadataResponse, error) {
	response := response.MetadataResponse{}
	path := viper.GetString("shipper.path.order-detail")
	path = strings.ReplaceAll(path, "{orderID}", orderID)
	url := h.Base + path

	header := h.Authorization
	header["Content-Type"] = "application/json"

	respByte, err := http_helper.Delete(url, header, request.CancelOrderShipperRequest{Reason: req.Body.Reason}, h.Logger)

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
