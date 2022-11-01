package shipping_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/helper/http_helper"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
)

type Grab interface {
	GetShippingRate(input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error)
	CreateDelivery(courierService *entity.CourierService, req *request.CreateDelivery) (*response.CreateDeliveryThirdPartyData, message.Message)
	GetTracking(orderID string) ([]response.GetOrderShippingTracking, message.Message)
}

type grab struct {
	Logger log.Logger
}

func NewGrab(log log.Logger) Grab {
	return &grab{
		Logger: log,
	}
}

func (g *grab) GetToken() string {
	req := &request.GrabAuthRequest{
		ClientID:     viper.GetString("grab.auth.client-id"),
		ClientSecret: viper.GetString("grab.auth.client-secret"),
		GrantType:    viper.GetString("grab.auth.grant-type"),
		Scope:        viper.GetString("grab.auth.scope"),
	}

	url := grabUrl(viper.GetString("grab.path.auth"))
	headers := map[string]string{
		"Cache-Control": "no-cache",
		"Content-Type":  "application/json",
	}

	respByte, err := http_helper.Post(url, headers, req, g.Logger)
	if err != nil {
		return ""
	}

	resp := map[string]interface{}{}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return ""
	}

	token, ok := resp["access_token"]
	if !ok {
		return ""
	}

	return fmt.Sprint("Bearer ", token)
}

func (g *grab) GetShippingRate(input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error) {

	if checkCoordinate, msg := input.CheckCoordinate(); !checkCoordinate {
		return &response.ShippingRateCommonResponse{
			Rate:       make(map[string]response.ShippingRateData),
			CourierMsg: map[string]message.Message{GrabCode: msg},
		}, errors.New(msg.Message)
	}

	originLat, _ := strconv.ParseFloat(input.Origin.Latitude, 64)
	originLong, _ := strconv.ParseFloat(input.Origin.Longitude, 64)
	destinationLat, _ := strconv.ParseFloat(input.Destination.Latitude, 64)
	destinationLong, _ := strconv.ParseFloat(input.Destination.Longitude, 64)

	volumeWeight := util.CalculateVolumeWeightKg(float64(input.TotalWidth),
		float64(input.TotalLength),
		float64(input.TotalHeight))

	volumeWeightGram := volumeWeight * 1000
	weightGram := input.TotalWeight * 1000
	req := &request.GrabDeliveryQuotes{
		Origin: request.Origin{
			Address: "",
			Coordinates: request.Coordinates{
				Latitude:  originLat,
				Longitude: originLong,
			},
		},
		Destination: request.Destination{
			Address: "",
			Coordinates: request.Coordinates{
				Latitude:  destinationLat,
				Longitude: destinationLong,
			},
		},
		Packages: []request.Package{
			{
				Name:        fmt.Sprintf("grab-shipping-rate %s", input.ChannelCode),
				Description: fmt.Sprintf("shipping-item %s", input.ChannelCode),
				Quantity:    1,
				Price:       int(input.TotalProductPrice),
				Dimensions: request.Dimensions{
					Height: int(input.TotalHeight),
					Width:  int(input.TotalWidth),
					Depth:  int(input.TotalLength),
					Weight: int(math.Max(volumeWeightGram, weightGram)),
				},
			},
		},
	}

	resp, err := g.GetDeliveryQuote(req)
	if err != nil {
		msg := message.ShippingProviderMsg
		msg.Message = err.Error()
		return &response.ShippingRateCommonResponse{
			Rate:       make(map[string]response.ShippingRateData),
			CourierMsg: map[string]message.Message{GrabCode: msg},
		}, errors.New(msg.Message)
	}

	return resp.ToShippingRate(), nil
}

func (g *grab) GetDeliveryQuote(req *request.GrabDeliveryQuotes) (*response.GrabDeliveryQuotes, error) {
	url := grabUrl(viper.GetString("grab.path.get-delivery-quote"))
	headers, err := g.setRequestHeader()
	if err != nil {
		return nil, err
	}

	respByte, err := http_helper.Post(url, headers, req, g.Logger)
	if err != nil {
		return nil, err
	}

	resp := &response.GrabDeliveryQuotes{}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Quotes) > 0 {
		return resp, nil
	}

	errResp := &response.GrabError{}
	err = json.Unmarshal(respByte, &errResp)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(errResp.GetReason())
}

func (g *grab) CreateDelivery(courierService *entity.CourierService, req *request.CreateDelivery) (*response.CreateDeliveryThirdPartyData, message.Message) {
	if ok, msg := req.CheckCoordinate(); !ok {
		return nil, msg
	}
	originLat, _ := strconv.ParseFloat(req.Origin.Latitude, 64)
	originLong, _ := strconv.ParseFloat(req.Origin.Longitude, 64)
	destinationLat, _ := strconv.ParseFloat(req.Destination.Latitude, 64)
	destinationLong, _ := strconv.ParseFloat(req.Destination.Longitude, 64)

	now := time.Now().Add(5 * time.Second)
	grabReq := &request.CreateDeliveryGrab{
		MerchantOrderID: req.OrderNo,
		ServiceType:     strings.ToUpper(courierService.ShippingCode),
		Sender: request.GrabSenderRecipient{
			FirstName:   req.Merchant.Name,
			LastName:    "",
			Title:       "",
			CompanyName: "",
			Email:       req.Merchant.Email,
			Phone:       req.Merchant.Phone,
			SmsEnabled:  false,
			Instruction: req.Notes,
		},
		Recipient: request.GrabSenderRecipient{
			FirstName:   req.Customer.Name,
			LastName:    "",
			Title:       "",
			CompanyName: "",
			Email:       req.Customer.Email,
			Phone:       req.Customer.Phone,
			SmsEnabled:  false,
			Instruction: req.Notes,
		},
		Packages: []request.Package{},
		Origin: request.Origin{
			Address: req.Origin.Address,
			Coordinates: request.Coordinates{
				Latitude:  originLat,
				Longitude: originLong,
			},
			Keywords: "",
		},
		Destination: request.Destination{
			Address: req.Destination.Address,
			Coordinates: request.Coordinates{
				Latitude:  destinationLat,
				Longitude: destinationLong,
			},
			Keywords: "",
		},
		PaymentMethod: "CASHLESS",
		Schedule: request.Schedule{
			PickupTimeFrom: now.Format(time.RFC3339),
			PickupTimeTo:   now.Add(time.Hour).Format(time.RFC3339),
		},
	}

	for _, v := range req.Package.Product {
		grabReq.Packages = append(grabReq.Packages, request.Package{
			Name:        v.Name,
			Description: "",
			Quantity:    v.Qty,
			Price:       int(v.Price),
			Dimensions: request.Dimensions{
				Height: int(req.Package.TotalHeight),
				Width:  int(req.Package.TotalWidth),
				Depth:  int(req.Package.TotalLength),
				Weight: int(req.Package.TotalWeight * 1000),
			},
		})
	}

	order, err := g.CreateOrder(grabReq)
	if err != nil {
		msg := message.ShippingProviderMsg
		msg.Message = err.Error()
		return nil, msg
	}

	return &response.CreateDeliveryThirdPartyData{
		BookingID:          order.DeliveryID,
		ShippingCost:       order.Quote.Amount,
		TotalShippingCost:  order.Quote.Amount,
		ActualShippingCost: order.Quote.Amount,
		Status:             StatusRequestPickup,
		PickUpCode:         order.DeliveryID,
		Airwaybill:         order.DeliveryID,
		//PickUpTime: order.Timeline.Pickup,
		// Insurance: false,
		// InsuranceCost: 0,
	}, message.SuccessMsg

}

func (g *grab) CreateOrder(req *request.CreateDeliveryGrab) (*response.CreateDeliveryGrab, error) {
	url := grabUrl(viper.GetString("grab.path.create-delivery"))
	headers, err := g.setRequestHeader()
	if err != nil {
		return nil, err
	}

	respByte, err := http_helper.Post(url, headers, req, g.Logger)
	if err != nil {
		return nil, err
	}

	resp := &response.CreateDeliveryGrab{}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.DeliveryID) > 0 {
		return resp, nil
	}

	errResp := &response.GrabError{}
	err = json.Unmarshal(respByte, &errResp)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(errResp.GetReason())
}

func (g *grab) GetOrderDetail(deliveryID string) (*response.GrabDeliveryDetail, error) {
	url := grabUrl(viper.GetString("grab.path.delivery-detail"))
	url = strings.ReplaceAll(url, "{deliveryID}", deliveryID)
	headers, err := g.setRequestHeader()
	if err != nil {
		return nil, err
	}

	respByte, err := http_helper.Get(url, headers, map[string]string{}, g.Logger)

	if err != nil {
		return nil, err
	}

	resp := response.GrabDeliveryDetail{}
	err = json.Unmarshal(respByte, &resp)

	if err != nil {
		return nil, err
	}

	if len(resp.DeliveryID) > 0 {
		return &resp, nil
	}

	errResp := &response.GrabError{}
	err = json.Unmarshal(respByte, &errResp)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(errResp.GetReason())
}

func (g *grab) GetTracking(orderID string) ([]response.GetOrderShippingTracking, message.Message) {
	logger := log.With(g.Logger, "Grab", "GetTracking")

	orderDetail, err := g.GetOrderDetail(orderID)
	if err != nil {
		_ = level.Error(logger).Log("g.GetOrderDetail", err.Error())
		return nil, message.ErrGetOrderDetail
	}

	return orderDetail.ToOrderShippingTracking(), message.SuccessMsg
}

func (g *grab) setRequestHeader() (map[string]string, error) {
	auth := g.GetToken()
	if len(auth) == 0 {
		return make(map[string]string), errors.New("grab unauthorized")
	}

	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": auth,
	}, nil
}

func grabUrl(path string) string {
	base := viper.GetString("grab.base")
	return base + path
}
