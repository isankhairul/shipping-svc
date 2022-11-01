package service

import (
	"encoding/json"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/http_helper"
	"go-klikdokter/helper/http_helper/shipping_provider"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/cache"
	"go-klikdokter/pkg/util"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
)

type ShippingService interface {
	GetShippingRate(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message)
	GetShippingRateByShippingType(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message)
	CreateDelivery(input *request.CreateDelivery) (*response.CreateDelivery, message.Message)
	OrderShippingTracking(req *request.GetOrderShippingTracking) ([]response.GetOrderShippingTracking, message.Message)
	UpdateStatusShipper(req *request.WebhookUpdateStatusShipper) (*entity.OrderShipping, message.Message)
	GetOrderShippingList(req *request.GetOrderShippingList) ([]response.GetOrderShippingList, *base.Pagination, message.Message)
	GetOrderShippingDetailByUID(uid string) (*response.GetOrderShippingDetail, message.Message)
	CancelPickup(req *request.CancelPickup) message.Message
	CancelOrder(req *request.CancelOrder) message.Message
	GetOrderShippingLabel(req *request.GetOrderShippingLabel) ([]response.GetOrderShippingLabelResponse, message.Message)
	RepickupOrder(req *request.RepickupOrderRequest) (*response.RepickupOrderResponse, message.Message)
	ShippingTracking(req *request.GetOrderShippingTracking) ([]response.GetOrderShippingTracking, message.Message)
	UpdateStatusGrab(req *request.WebhookUpdateStatusGrabRequest) message.Message
}

type shippingServiceImpl struct {
	logger                    log.Logger
	baseRepo                  repository.BaseRepository
	channelRepo               repository.ChannelRepository
	courierServiceRepo        repository.CourierServiceRepository
	courierCoverageCode       repository.CourierCoverageCodeRepository
	shipper                   shipping_provider.Shipper
	redis                     cache.RedisCache
	orderShipping             repository.OrderShippingRepository
	courierRepo               repository.CourierRepository
	shippingCourierStatusRepo repository.ShippingCourierStatusRepository
	daprEndpoint              http_helper.DaprEndpoint
	grab                      shipping_provider.Grab
}

func NewShippingService(
	l log.Logger,
	br repository.BaseRepository,
	chrp repository.ChannelRepository,
	csrp repository.CourierServiceRepository,
	cccrp repository.CourierCoverageCodeRepository,
	sh shipping_provider.Shipper,
	rc cache.RedisCache,
	osr repository.OrderShippingRepository,
	cr repository.CourierRepository,
	scs repository.ShippingCourierStatusRepository,
	de http_helper.DaprEndpoint,
	gr shipping_provider.Grab,
) ShippingService {
	return &shippingServiceImpl{
		l, br, chrp, csrp, cccrp, sh, rc, osr, cr, scs, de, gr,
	}
}

// swagger:operation POST /shipping/shipping-rate/{shipping-type} Shipping ShippingRateByShippingType
// Get Shipping Rate By Shipping Type
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/ShippmentPredefinedDetail'
func (s *shippingServiceImpl) GetShippingRateByShippingType(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message) {
	if input.ShippingType == "" {
		return []response.GetShippingRateResponse{}, message.ErrShippingTypeRequired
	}
	return s.GetShippingRate(input)
}

// swagger:operation POST /shipping/shipping-rate Shipping ShippingRate
// Get Shipping Rate
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/ShippingRate'
func (s *shippingServiceImpl) GetShippingRate(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message) {
	logger := log.With(s.logger, "ShippingService", "GetShippingRate")

	if len(input.CourierServiceUID) == 0 {
		return []response.GetShippingRateResponse{}, message.ErrCourierServiceIsRequired
	}

	// find Channel By UID
	channel, err := s.channelRepo.FindByUid(&input.ChannelUID)

	if err != nil {
		_ = level.Error(logger).Log("s.channelRepo.FindByUid", err.Error())
		return []response.GetShippingRateResponse{}, message.ErrChannelNotFound
	}

	if channel == nil {
		return []response.GetShippingRateResponse{}, message.ErrChannelNotFound
	}

	input.ChannelCode = channel.ChannelCode

	// find Courier Servies By Channel UID and Courier Servies UID Slice
	courierServices, err := s.courierServiceRepo.FindCourierServiceByChannelAndUIDs(input.ChannelUID, input.CourierServiceUID, input.ContainPrescription, input.ShippingType)

	if err != nil {
		_ = level.Error(logger).Log("s.courierServiceRepo.FindCourierServiceByChannelAndUIDs", err.Error())
		return []response.GetShippingRateResponse{}, message.CourierServiceNotFoundMsg
	}

	if len(courierServices) == 0 {
		return []response.GetShippingRateResponse{}, message.CourierServiceNotFoundMsg
	}

	price := s.getAllCourierPrice(courierServices, &input)

	return toGetShippingRateResponseList(&input, courierServices, price), message.SuccessMsg
}

// function to populate price data
func (s *shippingServiceImpl) getAllCourierPrice(courierServices []entity.ChannelCourierServiceForShippingRate, req *request.GetShippingRateRequest) *response.ShippingRateCommonResponse {
	var resp = &response.ShippingRateCommonResponse{
		Rate:       make(map[string]response.ShippingRateData),
		Summary:    make(map[string]response.ShippingRateSummary),
		CourierMsg: make(map[string]message.Message),
	}

	// Populate valid courier list distinct
	var courierList []entity.Courier
	courier := make(map[uint64]string)
	for _, v := range courierServices {

		_, isValid := v.IsValidCourier()
		_, isExist := courier[v.CourierID]

		if isValid && !isExist {
			courier[v.CourierID] = v.CourierCode
			courierList = append(courierList, entity.Courier{
				BaseIDModel: base.BaseIDModel{
					ID:  v.CourierID,
					UID: v.CourierUID,
				},
				Code:        v.CourierCode,
				CourierType: v.CourierTypeCode,
			})
		}
	}

	// Get all internal and merchant price if any
	internalAndMerchantPrice := s.internalAndMerchantPrice(courierList, courierServices, req)

	// Get third party price if any
	thirdPartyPrice := s.getThirdPartyPrice(courierList, req)

	resp.Add(internalAndMerchantPrice)
	resp.Add(thirdPartyPrice)
	return resp
}

// function to populate internal and merchant price
func (s *shippingServiceImpl) internalAndMerchantPrice(courierList []entity.Courier, courierService []entity.ChannelCourierServiceForShippingRate, req *request.GetShippingRateRequest) *response.ShippingRateCommonResponse {
	var (
		courierIDs []uint64
		price      = &response.ShippingRateCommonResponse{Rate: make(map[string]response.ShippingRateData)}
	)

	for _, v := range courierList {
		if v.CourierType != shipping_provider.InternalCourier && v.CourierType != shipping_provider.MerchantCourier {
			continue
		}

		courierIDs = append(courierIDs, v.ID)
	}

	if len(courierIDs) == 0 {
		return price
	}

	var (
		volume       = util.CalculateVolume(req.TotalHeight, req.TotalWidth, req.TotalLength)
		volumeWeight = util.CalculateVolumeWeightKg(req.TotalHeight, req.TotalWidth, req.TotalLength)
		finalWeight  = math.Max(req.TotalWeight, volumeWeight)

		lat1, _  = strconv.ParseFloat(req.Origin.Latitude, 64)
		long1, _ = strconv.ParseFloat(req.Origin.Longitude, 64)
		lat2, _  = strconv.ParseFloat(req.Destination.Latitude, 64)
		long2, _ = strconv.ParseFloat(req.Destination.Longitude, 64)

		distance = util.CalculateDistanceInKm(lat1, long1, lat2, long2)
	)

	// Check courier coverage
	origin := s.courierCoverageCode.FindInternalAndMerchantCourierCoverage(courierIDs, req.Origin.CountryCode, req.Origin.PostalCode)
	destination := s.courierCoverageCode.FindInternalAndMerchantCourierCoverage(courierIDs, req.Destination.CountryCode, req.Destination.PostalCode)

	for _, v := range courierService {

		_, originOK := origin[v.CourierCode]
		_, destinationOK := destination[v.CourierCode]

		if v.CourierTypeCode != shipping_provider.InternalCourier && v.CourierTypeCode != shipping_provider.MerchantCourier {
			continue
		}

		key := global.CourierShippingCodeKey(v.CourierCode, v.ShippingCode)
		value := response.ShippingRateData{
			Distance:         distance,
			TotalPrice:       v.Price,
			InsuranceFee:     v.InsuranceFee,
			InsuranceApplied: v.UseInsurance,
			MustUseInsurance: v.UseInsurance,
			Weight:           req.TotalWeight,
			Volume:           volume,
			VolumeWeight:     volumeWeight,
			FinalWeight:      finalWeight,
			MinDay:           0,
			MaxDay:           0,
			UnitPrice:        0,
			AvailableCode:    200,
			Error:            response.SetShippingRateErrorMessage(message.SuccessMsg),
		}

		if v.CourierTypeCode == shipping_provider.MerchantCourier {
			value.TotalPrice = 0
		}

		// check if origin or destination not available
		if !(originOK && destinationOK) {
			value.UpdateMessage(message.ErrCourierCoverageCodeUidNotExist)
		}

		price.Rate[key] = value
	}

	return price
}

// function to get shipping rate from third party shipping provider
func (s *shippingServiceImpl) getThirdPartyPrice(courier []entity.Courier, input *request.GetShippingRateRequest) *response.ShippingRateCommonResponse {
	var resp = &response.ShippingRateCommonResponse{
		Rate:       make(map[string]response.ShippingRateData),
		Summary:    map[string]response.ShippingRateSummary{},
		CourierMsg: map[string]message.Message{},
	}

	var wg sync.WaitGroup
	wg.Add(len(courier))
	for _, v := range courier {

		go func(c entity.Courier) {
			defer wg.Done()
			var (
				courierPrice *response.ShippingRateCommonResponse
				err          error
			)

			if c.CourierType != shipping_provider.ThirPartyCourier {
				return
			}

			// try to get price data from cache
			baseKey := viper.GetString("cache.redis.base-key")
			key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s:%s:%s:%f",
				baseKey,
				c.Code,
				input.Origin.PostalCode,
				input.Destination.PostalCode,
				input.Origin.Subdistrict,
				input.Destination.Subdistrict,
				input.Origin.Latitude,
				input.Origin.Longitude,
				input.Destination.Latitude,
				input.Destination.Longitude,
				input.TotalWeight,
			)

			_ = s.redis.GetJsonStruct(key, &courierPrice)
			// if cache doesn't exist
			if courierPrice == nil {
				switch c.Code {
				case shipping_provider.ShipperCode:
					courierPrice, err = s.shipper.GetShippingRate(&c.ID, input)
				case shipping_provider.GrabCode:
					courierPrice, err = s.grab.GetShippingRate(input)
				default:
					resp.CourierMsg[c.Code] = message.InvalidCourierCodeMsg
					return
				}

				if err == nil {
					// save price to redis cache
					s.redis.SetJsonStruct(key, courierPrice, viper.GetInt("cache.redis.expired-in-minute.shipping-rate"))
				}

			}

			if courierPrice != nil {
				resp.Add(courierPrice)
			}
		}(v)
	}
	wg.Wait()

	return resp
}

// function to generate ShippingRateResponseList
func toGetShippingRateResponseList(req *request.GetShippingRateRequest, courierServices []entity.ChannelCourierServiceForShippingRate, price *response.ShippingRateCommonResponse) []response.GetShippingRateResponse {
	shippingTypeMap := make(map[string][]response.GetShippingRateService)
	var resp []response.GetShippingRateResponse

	estimation := func(courierEtd, dbEtd float64) float64 {
		if courierEtd > 0.0 {
			return util.RoundFloat(courierEtd, 2)
		}
		return dbEtd
	}

	for _, v := range courierServices {
		p := price.FindShippingCode(v.CourierCode, v.ShippingCode)

		if msg := v.Validate(&p.FinalWeight, &req.ContainPrescription); msg != message.SuccessMsg {
			p.UpdateMessage(msg)
		}

		if p.AvailableCode != 200 {
			v.EtdMin = 0
			v.EtdMax = 0
		}

		service := response.GetShippingRateService{
			Courier: response.GetShippingRateCourir{
				CourierUID:      v.CourierUID,
				CourierCode:     v.CourierCode,
				CourierName:     v.CourierName,
				CourierTypeCode: v.CourierTypeCode,
				CourierTypeName: v.CourierTypeName,
			},
			CourierServiceUID:       v.CourierServiceUID,
			ShippingCode:            v.ShippingCode,
			ShippingName:            v.ShippingName,
			ShippingTypeCode:        v.ShippingTypeCode,
			ShippingTypeName:        v.ShippingTypeName,
			ShippingTypeDescription: v.ShippingTypeDescription,
			Logo:                    v.Logo,
			Weight:                  req.TotalWeight,

			AvailableCode:    p.AvailableCode,
			Error:            p.Error,
			Volume:           p.Volume,
			VolumeWeight:     p.VolumeWeight,
			FinalWeight:      p.FinalWeight,
			MinDay:           p.MinDay,
			MaxDay:           p.MaxDay,
			UnitPrice:        p.UnitPrice,
			TotalPrice:       p.TotalPrice,
			InsuranceFee:     p.InsuranceFee,
			MustUseInsurance: p.MustUseInsurance,
			InsuranceApplied: p.InsuranceApplied,
			Distance:         p.Distance,

			// uses date from courier
			// if it does not exist get uses from database
			Etd_Min: estimation(p.Etd_Min, v.EtdMin),
			Etd_Max: estimation(p.Etd_Max, v.EtdMax),
		}

		price.SummaryPerShippingType(v.ShippingTypeCode, service.TotalPrice, service.Etd_Max, service.Etd_Min, service.AvailableCode)

		if _, ok := shippingTypeMap[v.ShippingTypeCode]; !ok {
			shippingTypeMap[v.ShippingTypeCode] = []response.GetShippingRateService{}
		}

		shippingTypeMap[v.ShippingTypeCode] = append(shippingTypeMap[v.ShippingTypeCode], service)
	}

	for k, v := range shippingTypeMap {
		s := price.Summary[k]
		data := response.GetShippingRateResponse{
			ShippingTypeCode:        k,
			ShippingTypeName:        v[0].ShippingTypeName,
			ShippingTypeDescription: v[0].ShippingTypeDescription,
			PriceRange:              s.PriceRange,
			EtdMax:                  s.EtdMax,
			EtdMin:                  s.EtdMin,
			Services:                v,
			AvailableCode:           s.AvailableCode,
			Error:                   s.Error,
		}
		resp = append(resp, data)
	}

	return resp
}

func (s *shippingServiceImpl) populateCreateDelivery(input *request.CreateDelivery) (*entity.CourierService, *entity.OrderShipping, *entity.ShippingCourierStatus, *entity.ShippingCourierStatus, message.Message) {
	logger := log.With(s.logger, "ShippingService", "PopulateCreateDelivery")

	// find Channel By UID
	channel, err := s.channelRepo.FindByUid(&input.ChannelUID)
	if err != nil {
		_ = level.Error(logger).Log("s.channelRepo.FindByUid", err.Error())
		return nil, nil, nil, nil, message.ChannelNotFoundMsg
	}

	if channel == nil {
		return nil, nil, nil, nil, message.ChannelNotFoundMsg
	}

	// find Courier Service By UID
	courierService, err := s.courierServiceRepo.FindCourierService(input.ChannelUID, input.CouirerServiceUID)
	if err != nil {
		_ = level.Error(logger).Log("s.courierServiceRepo.FindCourierService", err.Error())
		return nil, nil, nil, nil, message.CourierServiceNotFoundMsg
	}

	if courierService == nil {
		return nil, nil, nil, nil, message.CourierServiceNotFoundMsg
	}

	if msg := courierService.Validate(input.Package.TotalWeight, input.Package.ContainPrescription > 0); msg != message.SuccessMsg {
		return nil, nil, nil, nil, msg
	}

	// check if order no already exist with status created
	orderShipping, err := s.orderShipping.FindByOrderNo(input.OrderNo)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByOrderNo", err.Error())
		return nil, nil, nil, nil, message.ErrDB
	}

	if orderShipping != nil && orderShipping.Status != shipping_provider.StatusCreated {
		return nil, nil, nil, nil, message.OrderNoAlreadyExistsMsg
	}

	// get shipping status
	createdStatus, _ := s.shippingCourierStatusRepo.FindByCode(channel.ID, courierService.CourierID, shipping_provider.StatusCreated)
	if createdStatus == nil {
		return nil, nil, nil, nil, message.ShippingStatusNotFoundMsg
	}

	requestPickupStatus, _ := s.shippingCourierStatusRepo.FindByCode(channel.ID, courierService.CourierID, shipping_provider.StatusRequestPickup)
	if requestPickupStatus == nil {
		return nil, nil, nil, nil, message.ShippingStatusNotFoundMsg
	}

	// if order no doesn't exist create new one
	if orderShipping == nil {
		orderShipping = &entity.OrderShipping{}
		orderShipping.FromCreateDeliveryRequest(input)
		orderShipping.ChannelID = channel.ID
		orderShipping.CourierID = courierService.CourierID
		orderShipping.CourierServiceID = courierService.ID
	}
	orderShipping.UpdatedBy = input.Username
	return courierService, orderShipping, createdStatus, requestPickupStatus, message.SuccessMsg
}

// swagger:operation POST /shipping/order-shipping Shipping CreateDelivery
// Create Order Shipping
//
// Description :
//
// ---
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/CreateDeliveryResponse'
func (s *shippingServiceImpl) CreateDelivery(input *request.CreateDelivery) (*response.CreateDelivery, message.Message) {
	logger := log.With(s.logger, "ShippingService", "CreateDelivery")

	courierService, orderShipping, created, requestPickup, msg := s.populateCreateDelivery(input)
	if msg != message.SuccessMsg {
		return &response.CreateDelivery{}, msg
	}

	msg = s.createDelivery(orderShipping, courierService, input)
	if msg != message.SuccessMsg {
		return &response.CreateDelivery{}, msg
	}

	orderShipping.AddHistoryStatus(created, fmt.Sprintf("Booking ID [%s]", orderShipping.BookingID))

	if orderShipping.Status == shipping_provider.StatusRequestPickup {
		orderShipping.AddHistoryStatus(requestPickup, fmt.Sprintf("Pickup Code [%s]", *orderShipping.PickupCode))
	}

	orderShipping, err := s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return &response.CreateDelivery{}, message.ErrSaveOrderShipping
	}

	return &response.CreateDelivery{
		OrderNoAPI:       input.OrderNo,
		OrderShippingUID: orderShipping.UID,
	}, message.SuccessMsg
}

func (s *shippingServiceImpl) createDelivery(orderShipping *entity.OrderShipping, courierService *entity.CourierService, input *request.CreateDelivery) message.Message {
	switch courierService.Courier.CourierType {
	case shipping_provider.ThirPartyCourier, shipping_provider.AggregatorCourier:

		orderData, msg := s.createDeliveryThirdParty(orderShipping.BookingID, courierService, input)
		if msg != message.SuccessMsg {
			return msg
		}

		if orderShipping.ID == 0 {
			orderShipping.CreatedBy = input.Username
			orderShipping.Insurance = orderData.Insurance
			orderShipping.InsuranceCost = orderData.InsuranceCost
			orderShipping.ShippingCost = orderData.ShippingCost
			orderShipping.TotalShippingCost = orderData.TotalShippingCost
			orderShipping.ActualShippingCost = orderData.ActualShippingCost
			orderShipping.BookingID = orderData.BookingID
		}
		orderShipping.PickupCode = &orderData.PickUpCode
		orderShipping.Airwaybill = orderData.Airwaybill
		orderShipping.Status = orderData.Status

	default:
		return message.ErrInvalidCourierType
	}

	return message.SuccessMsg
}

func (s *shippingServiceImpl) createDeliveryThirdParty(bookingID string, courierService *entity.CourierService, input *request.CreateDelivery) (*response.CreateDeliveryThirdPartyData, message.Message) {
	switch courierService.Courier.Code {
	case shipping_provider.ShipperCode:
		return s.shipper.CreateDelivery(bookingID, courierService, input)
	case shipping_provider.GrabCode:
		return s.grab.CreateDelivery(courierService, input)
	default:
		return nil, message.ErrInvalidCourierCode
	}
}

// swagger:operation GET /shipping/order-tracking/{uid} Shipping OrderShippingTracking
// Get Order Shipping Tracking
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/GetOrderShippingTrackingResponse'
func (s *shippingServiceImpl) OrderShippingTracking(req *request.GetOrderShippingTracking) ([]response.GetOrderShippingTracking, message.Message) {
	logger := log.With(s.logger, "ShippingService", "OrderShippingTracking")

	if len(req.ChannelUID) == 0 {
		return nil, message.ErrChannelUIDRequired
	}

	orderShipping, err := s.orderShipping.FindByUID(req.UID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByUID", err.Error())
		return nil, message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return nil, message.ErrOrderShippingNotFound
	}

	if orderShipping.Channel.UID != req.ChannelUID {
		return nil, message.ErrOrderBelongToAnotherChannel
	}

	var orderStatus []response.GetOrderShippingTracking
	var msg message.Message
	switch orderShipping.Courier.CourierType {
	case shipping_provider.ThirPartyCourier, shipping_provider.AggregatorCourier:
		orderStatus, msg = s.thridPartyTracking(orderShipping)
	default:
		return []response.GetOrderShippingTracking{}, message.ErrInvalidCourierType
	}

	return response.SortOrderStatusByTimeDesc(orderStatus), msg
}

func (s *shippingServiceImpl) thridPartyTracking(orderShipping *entity.OrderShipping) ([]response.GetOrderShippingTracking, message.Message) {
	switch orderShipping.Courier.Code {
	case shipping_provider.ShipperCode:
		return s.shipper.GetTracking(orderShipping.BookingID)
	case shipping_provider.GrabCode:
		return s.grab.GetTracking(orderShipping.BookingID)
	}

	return nil, message.ErrInvalidCourierCode
}

// swagger:operation POST /public/webhook/shipper Public WebhookUpdateStatusShipper
// Update Status from Shipper Webhook
//
// Description :
//
// ---
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           type: object
func (s *shippingServiceImpl) UpdateStatusShipper(req *request.WebhookUpdateStatusShipper) (*entity.OrderShipping, message.Message) {
	logger := log.With(s.logger, "ShippingService", "UpdateStatusShipper")

	if jsonReq, err := json.Marshal(req); err == nil {
		_ = level.Info(logger).Log("shipper_webhook", string(jsonReq))
	}

	_ = level.Info(logger).Log("check_webhook", shipping_provider.ShipperWebhookAuth())
	_ = level.Info(logger).Log("check_webhook_auth", req.Auth)
	if req.Auth != shipping_provider.ShipperWebhookAuth() {
		return nil, message.ErrUnAuth
	}

	orderShipping, err := s.orderShipping.FindByOrderNo(req.ExternalID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByOrderNo", err.Error())
		return nil, message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return nil, message.ErrOrderShippingNotFound
	}

	if orderShipping.Courier.Code != shipping_provider.ShipperCode {
		return nil, message.ErrOrderShippingNotFound
	}

	statusCode := req.ExternalStatus.Code
	statusDescription := req.ExternalStatus.Description

	shippingStatus, err := s.shippingCourierStatusRepo.FindByCourierStatus(orderShipping.CourierID, fmt.Sprint(statusCode))

	if err != nil {
		_ = level.Error(logger).Log("s.shippingCourierStatusRepo.FindByCourierStatus", err.Error())
		return nil, message.ShippingStatusNotFoundMsg
	}

	if shippingStatus == nil {
		return nil, message.ShippingStatusNotFoundMsg
	}

	orderShipping.Status = shippingStatus.StatusCode
	orderShipping.UpdatedBy = "SHIPPER_WEBHOOK"

	if len(req.Awb) > 0 {
		orderShipping.Airwaybill = req.Awb
	}

	driverInfo := SplitDriverInfo(req.External.Description)
	orderShipping.AddHistoryStatus(shippingStatus, statusDescription, driverInfo.Description())

	orderShipping, err = s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return nil, message.ErrSaveOrderShipping
	}

	topic := viper.GetString("dapr.topic.update-order-shipping")
	topic = strings.ReplaceAll(topic, "{channel-code}", strings.ToLower(orderShipping.Channel.ChannelCode))
	updateOrderRequest := request.UpdateOrderShippingBody{
		ChannelUID:         orderShipping.Channel.UID,
		CourierCode:        orderShipping.Courier.Code,
		CourierServiceUID:  orderShipping.CourierService.UID,
		OrderNo:            orderShipping.OrderNo,
		OrderShippingUID:   orderShipping.UID,
		Airwaybill:         orderShipping.Airwaybill,
		ShippingStatus:     shippingStatus.StatusCode,
		ShippingStatusName: shippingStatus.ShippingStatus.StatusName,
		UpdatedBy:          "shipping_service",
		Timestamp:          time.Now(),
		Details: request.UpdateOrderShippingBodyDetail{
			ExternalStatusCode:        fmt.Sprint(req.ExternalStatus.Code),
			ExternalStatusName:        req.ExternalStatus.Name,
			ExternalStatusDescription: req.ExternalStatus.Description,
		},
		DriverInfo: driverInfo,
	}

	_ = level.Info(logger).Log("PUBLISH_QUEUE, TOPIC", topic)
	s.daprEndpoint.PublishKafka(topic, updateOrderRequest)
	return orderShipping, message.SuccessMsg
}

// swagger:operation GET /shipping/order-shipping Shipping GetOrderShippingList
// Get Order Shipping List
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//           $ref: '#/definitions/MetaPaginationResponse'
//         data:
//           properties:
//             records:
//               type: array
//               items:
//                 $ref: '#/definitions/GetOrderShippingListResponse'
func (s *shippingServiceImpl) GetOrderShippingList(req *request.GetOrderShippingList) ([]response.GetOrderShippingList, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "ShippingService", "GetOrderShippingList")

	if len(req.Filters.OrderShippingDateFrom) > 0 {
		if ok := util.DateValidationYYYYMMDD(req.Filters.OrderShippingDateFrom); !ok {
			return []response.GetOrderShippingList{}, &base.Pagination{}, message.ErrFormatDateYYYYMMDD
		}
	}

	if len(req.Filters.OrderShippingDateTo) > 0 {
		if ok := util.DateValidationYYYYMMDD(req.Filters.OrderShippingDateTo); !ok {
			return []response.GetOrderShippingList{}, &base.Pagination{}, message.ErrFormatDateYYYYMMDD
		}
	}

	filter := make(map[string]interface{})
	filter["order_shipping_uid"] = req.Filters.OrderShippingUID
	filter["order_no"] = req.Filters.OrderNo
	filter["channel_code"] = req.Filters.ChannelCode
	filter["channel_name"] = req.Filters.ChannelName
	filter["courier_name"] = req.Filters.CourierName
	filter["courier_services_name"] = req.Filters.CourierServicesName
	filter["airwaybill"] = req.Filters.Airwaybill
	filter["shipping_status"] = req.Filters.ShippingStatus
	filter["order_shipping_date_from"] = req.Filters.OrderShippingDateFrom
	filter["order_shipping_date_to"] = req.Filters.OrderShippingDateTo
	filter["booking_id"] = req.Filters.BookingID
	filter["merchant_name"] = req.Filters.MerchantName
	filter["customer_name"] = req.Filters.CustomerName

	result, pagination, err := s.orderShipping.FindByParams(req.Limit, req.Page, req.Sort, filter)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByParams", err.Error())
		return result, pagination, message.ErrNoData
	}

	return result, pagination, message.SuccessMsg
}

// swagger:operation GET /shipping/order-shipping/{uid} Shipping GetOrderShippingDetail
// Get Order Shipping Detail By UID
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/GetOrderShippingDetailResponse'
func (s *shippingServiceImpl) GetOrderShippingDetailByUID(uid string) (*response.GetOrderShippingDetail, message.Message) {
	logger := log.With(s.logger, "ShippingService", "GetOrderShippingDetailByUID")
	var resp *response.GetOrderShippingDetail

	orderShipping, err := s.orderShipping.FindByUID(uid)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByUID", err.Error())
		return nil, message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return nil, message.ErrOrderShippingNotFound
	}

	resp = getOrderShippingDetailByUIDResponse(orderShipping)

	shipperStatus, _ := s.shippingCourierStatusRepo.FindByCode(orderShipping.ChannelID, orderShipping.CourierID, orderShipping.Status)

	if shipperStatus != nil {
		resp.ShippingStatusName = shipperStatus.ShippingStatus.StatusName
	}

	return resp, message.SuccessMsg
}

func getOrderShippingDetailByUIDResponse(orderShipping *entity.OrderShipping) *response.GetOrderShippingDetail {
	if orderShipping == nil {
		return nil
	}

	resp := &response.GetOrderShippingDetail{}
	resp.ChannelUID = orderShipping.Channel.UID
	resp.ChannelCode = orderShipping.Channel.ChannelCode
	resp.ChannelName = orderShipping.Channel.ChannelName
	resp.CourierName = orderShipping.Courier.CourierName
	resp.CourierServiceName = orderShipping.CourierService.ShippingName
	resp.OrderShippingUID = orderShipping.UID
	resp.OrderShippingDate = orderShipping.OrderShippingDate
	resp.Airwaybill = orderShipping.Airwaybill
	resp.BookingID = orderShipping.BookingID
	resp.OrderNo = orderShipping.OrderNo
	resp.OrderNoAPI = orderShipping.OrderNoAPI
	resp.ShippingStatus = orderShipping.Status
	resp.TotalProductPrice = orderShipping.TotalProductPrice
	resp.TotalWeight = orderShipping.TotalWeight
	resp.TotalVolume = orderShipping.TotalVolume
	resp.FinalWeight = orderShipping.TotalFinalWeight
	resp.TotalProductPrice = orderShipping.TotalProductPrice
	resp.ShippingCost = orderShipping.ShippingCost
	resp.Insurance = orderShipping.Insurance
	resp.InsuranceCost = orderShipping.InsuranceCost
	resp.TotalShippingCost = orderShipping.TotalShippingCost
	resp.ShippingNotes = orderShipping.ShippingNotes
	resp.MerchantUID = orderShipping.MerchantUID
	resp.MerchantName = orderShipping.MerchantName
	resp.MerchantEmail = orderShipping.MerchantEmail
	resp.MerchantPhone = orderShipping.MerchantPhoneNumber
	resp.MerchantAddress = orderShipping.MerchantAddress
	resp.MerchantSubdistrict = orderShipping.MerchantSubdistrict
	resp.MerchantDistrictName = orderShipping.MerchantDistrictName
	resp.MerchantCityName = orderShipping.MerchantCityName
	resp.MerchantProvinceName = orderShipping.MerchantProvinceName
	resp.MerchantPostalCode = orderShipping.MerchantPostalCode
	resp.CustomerUID = orderShipping.CustomerUID
	resp.CustomerName = orderShipping.CustomerName
	resp.CustomerEmail = orderShipping.CustomerEmail
	resp.CustomerPhone = orderShipping.CustomerPhoneNumber
	resp.CustomerAddress = orderShipping.CustomerAddress
	resp.CustomerSubdistrict = orderShipping.CustomerSubdistrict
	resp.CustomerDistrictName = orderShipping.CustomerDistrictName
	resp.CustomerCityName = orderShipping.CustomerCityName
	resp.CustomerProvinceName = orderShipping.CustomerProvinceName
	resp.CustomerPostalCode = orderShipping.CustomerPostalCode
	resp.CustomerNotes = orderShipping.CustomerNotes
	resp.OrderShippingItem = []response.GetOrderShippingDetailItem{}
	resp.OrderShippingHistory = []response.GetOrderShippingDetailHistory{}

	for _, v := range orderShipping.OrderShippingItem {
		resp.OrderShippingItem = append(resp.OrderShippingItem, response.GetOrderShippingDetailItem{
			ItemName:     v.ItemName,
			ProductUID:   v.ProductUID,
			Qty:          v.Quantity,
			Price:        v.Price,
			Weight:       v.Weight,
			Volume:       v.Volume,
			Prescription: v.Prescription,
		})
	}

	for _, v := range orderShipping.OrderShippingHistory {
		resp.OrderShippingHistory = append(resp.OrderShippingHistory, response.GetOrderShippingDetailHistory{
			CreatedAt:  v.CreatedAt,
			Status:     v.StatusCode,
			Notes:      v.Note,
			CreatedBy:  v.CreatedBy,
			StatusName: v.ShippingCourierStatus.ShippingStatus.StatusName,
		})
	}

	return resp
}

// swagger:operation POST /shipping/cancel-pickup/{uid} Shipping CancelPickup
// Cancel Pickup Order
//
// Description :
//
// ---
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           type: object
func (s *shippingServiceImpl) CancelPickup(req *request.CancelPickup) message.Message {
	logger := log.With(s.logger, "ShippingService", "CancelPickup")
	orderShipping, err := s.orderShipping.FindByUID(req.UID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByUID", err.Error())
		return message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return message.ErrOrderShippingNotFound
	}

	// check if courier service is cancelable
	if orderShipping.CourierService.Cancelable == 0 {
		return message.ErrCantCancelOrderCourierService
	}

	msg := s.cancelPickup(orderShipping)

	if msg != message.SuccessMsg {
		return msg
	}

	//status back to created
	shipperStatus, _ := s.shippingCourierStatusRepo.FindByCode(orderShipping.ChannelID, orderShipping.CourierID, shipping_provider.StatusCreated)

	if shipperStatus == nil {
		return message.ShippingStatusNotFoundMsg
	}

	notes := fmt.Sprintf("request pickup cancelled by merchant [%s]", *orderShipping.PickupCode)
	*orderShipping.PickupCode = ""
	orderShipping.Status = shipping_provider.StatusCreated
	orderShipping.UpdatedBy = req.Body.Username

	orderShipping.AddHistoryStatus(shipperStatus, notes)
	_, err = s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return message.ErrSaveOrderShipping
	}

	return message.SuccessMsg
}

func (s *shippingServiceImpl) cancelPickup(orderShipping *entity.OrderShipping) message.Message {
	switch orderShipping.Courier.CourierType {
	case shipping_provider.ThirPartyCourier, shipping_provider.AggregatorCourier:
		return s.cancelPickupThirdParty(orderShipping)
	}

	return message.ErrInvalidCourierType
}

func (s *shippingServiceImpl) cancelPickupThirdParty(orderShipping *entity.OrderShipping) message.Message {

	// check if order current status is cancelable
	if !shipping_provider.IsPickUpOrderCancelable(orderShipping.Courier.Code, orderShipping.Status) {
		return message.ErrCantCancelOrderShipping
	}

	var err error
	switch orderShipping.Courier.Code {
	case shipping_provider.ShipperCode:
		_, err = s.shipper.CancelPickupRequest(*orderShipping.PickupCode)
	}

	if err != nil {
		return message.ErrCancelPickup
	}

	return message.SuccessMsg
}

// swagger:operation POST /shipping/cancel-order/{uid} Shipping CancelOrder
// Cancel Order Shipping
//
// Description :
//
// ---
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           type: object
func (s *shippingServiceImpl) CancelOrder(req *request.CancelOrder) message.Message {
	logger := log.With(s.logger, "ShippingService", "CancelOrder")

	orderShipping, err := s.orderShipping.FindByUID(req.UID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByUID", err.Error())
		return message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return message.ErrOrderShippingNotFound
	}

	// check if courier service is cancelable
	if orderShipping.CourierService.Cancelable == 0 {
		return message.ErrCantCancelOrderCourierService
	}

	msg := s.cancelOrder(orderShipping, req)

	if msg != message.SuccessMsg {
		return msg
	}

	shipperStatus, _ := s.shippingCourierStatusRepo.FindByCode(orderShipping.ChannelID, orderShipping.CourierID, shipping_provider.StatusCancelled)

	if shipperStatus == nil {
		return message.ShippingStatusNotFoundMsg
	}

	orderShipping.Status = shipping_provider.StatusCancelled
	orderShipping.UpdatedBy = req.Body.Username
	orderShipping.AddHistoryStatus(shipperStatus, req.Body.Reason)

	_, err = s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return message.ErrSaveOrderShipping
	}

	return message.SuccessMsg
}

func (s *shippingServiceImpl) cancelOrder(orderShipping *entity.OrderShipping, req *request.CancelOrder) message.Message {
	switch orderShipping.Courier.CourierType {
	case shipping_provider.ThirPartyCourier, shipping_provider.AggregatorCourier:
		return s.cancelOrderThirdParty(orderShipping, req)
	}

	return message.ErrInvalidCourierType
}

func (s *shippingServiceImpl) cancelOrderThirdParty(orderShipping *entity.OrderShipping, req *request.CancelOrder) message.Message {

	// check if order current status is cancelable
	if !shipping_provider.IsOrderCancelable(orderShipping.Courier.Code, orderShipping.Status) {
		return message.ErrCantCancelOrderShipping
	}

	var err error
	switch orderShipping.Courier.Code {
	case shipping_provider.ShipperCode:
		_, err = s.shipper.CancelOrder(orderShipping.BookingID, req)
	}

	if err != nil {
		return message.ErrCancelPickup
	}

	return message.SuccessMsg
}

// swagger:operation POST /shipping/order-shipping-label/{channel-uid} Shipping GetOrderShippingLabel
// Order Shipping Label
//
// Description :
//
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//           $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             records:
//               type: array
//               items:
//                 $ref: '#/definitions/GetOrderShippingLabelResponse'
func (s *shippingServiceImpl) GetOrderShippingLabel(req *request.GetOrderShippingLabel) ([]response.GetOrderShippingLabelResponse, message.Message) {
	logger := log.With(s.logger, "ShippingService", "GetOrderShippingLabel")

	result, err := s.orderShipping.FindByUIDs(req.ChannelUID, req.Body.OrderShippingUID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByUIDs", err.Error())
		return []response.GetOrderShippingLabelResponse{}, message.ErrOrderShippingNotFound
	}

	return getOrderShippingLabelResponse(result, req.Body.HideProduct), message.SuccessMsg
}

func getOrderShippingLabelResponse(orderShipping []entity.OrderShipping, isHideProduct bool) []response.GetOrderShippingLabelResponse {
	result := []response.GetOrderShippingLabelResponse{}
	for _, v := range orderShipping {
		data := response.GetOrderShippingLabelResponse{
			OrderShippingItems:   []response.GetOrderShippingDetailItem{},
			ChannelCode:          v.Channel.ChannelCode,
			ChannelName:          v.Channel.ChannelName,
			ChannelImage:         v.Channel.ImagePath,
			OrderShippingUID:     v.UID,
			OrderShippingDate:    v.OrderShippingDate,
			OrderNo:              v.OrderNo,
			OrderNoAPI:           v.OrderNoAPI,
			CourierName:          v.Courier.CourierName,
			CourierImage:         v.Courier.ImagePath,
			CourierServiceName:   v.CourierService.ShippingName,
			CourierServiceImage:  v.CourierService.ImagePath,
			Airwaybill:           v.Airwaybill,
			BookingID:            v.BookingID,
			TotalProductPrice:    v.TotalProductPrice,
			TotalWeight:          v.TotalWeight,
			TotalVolume:          v.TotalVolume,
			FinalWeight:          v.TotalFinalWeight,
			ShippingCost:         v.ShippingCost,
			Insurance:            v.Insurance,
			InsuranceCost:        v.InsuranceCost,
			TotalShippingCost:    v.TotalShippingCost,
			ShippingNotes:        v.ShippingNotes,
			MerchantUID:          v.MerchantUID,
			MerchantName:         v.MerchantName,
			MerchantEmail:        v.MerchantEmail,
			MerchantPhone:        v.MerchantPhoneNumber,
			MerchantAddress:      v.MerchantAddress,
			MerchantDistrictName: v.MerchantDistrictName,
			MerchantCityName:     v.MerchantCityName,
			MerchantProvinceName: v.MerchantProvinceName,
			MerchantPostalCode:   v.MerchantPostalCode,
			CustomerUID:          v.CustomerUID,
			CustomerName:         v.CustomerName,
			CustomerEmail:        v.CustomerEmail,
			CustomerPhone:        v.CustomerPhoneNumber,
			CustomerAddress:      v.CustomerAddress,
			CustomerDistrictName: v.CustomerDistrictName,
			CustomerCityName:     v.CustomerCityName,
			CustomerProvinceName: v.CustomerProvinceName,
			CustomerPostalCode:   v.CustomerPostalCode,
			CustomerNotes:        v.CustomerNotes,
		}

		if !isHideProduct {
			items := []response.GetOrderShippingDetailItem{}
			for _, v := range v.OrderShippingItem {
				items = append(items, response.GetOrderShippingDetailItem{
					ItemName:     v.ItemName,
					ProductUID:   v.ProductUID,
					Qty:          v.Quantity,
					Price:        v.Price,
					Weight:       v.Weight,
					Volume:       v.Volume,
					Prescription: v.Prescription,
				})
			}
			data.OrderShippingItems = items
		}

		result = append(result, data)
	}
	return result
}

// swagger:operation POST /shipping/repickup Shipping RepickupOrder
// Repickup Order
//
// Description :
//
// ---
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//           $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//                 $ref: '#/definitions/RepickupOrderResponse'
func (s *shippingServiceImpl) RepickupOrder(req *request.RepickupOrderRequest) (*response.RepickupOrderResponse, message.Message) {
	logger := log.With(s.logger, "ShippingService", "RepickupOrder")
	resp := &response.RepickupOrderResponse{}
	orderShipping, err := s.orderShipping.FindByUID(req.OrderShippingUID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByUIDs", err.Error())
		return resp, message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return resp, message.ErrOrderShippingNotFound
	}

	if orderShipping.Channel.UID != req.ChannelUID {
		return resp, message.ErrOrderBelongToAnotherChannel
	}

	msg := s.repickupOrder(orderShipping)

	if msg != message.SuccessMsg {
		return resp, msg
	}

	//update order status
	orderShipping.Status = shipping_provider.StatusRequestPickup

	shippingStatus, _ := s.shippingCourierStatusRepo.FindByCode(
		orderShipping.ChannelID,
		orderShipping.CourierID,
		shipping_provider.StatusRequestPickup)

	if shippingStatus == nil {
		return resp, message.ShippingStatusNotFoundMsg
	}

	orderShipping.UpdatedBy = req.Username

	//add history
	orderShipping.OrderShippingHistory = append(orderShipping.OrderShippingHistory, entity.OrderShippingHistory{
		OrderShippingID:         orderShipping.ID,
		ShippingCourierStatusID: shippingStatus.ID,
		StatusCode:              shippingStatus.StatusCode,
		Note:                    fmt.Sprintf("(Repickup) Pickup Code [%s]", *orderShipping.PickupCode),
		BaseIDModel: base.BaseIDModel{
			CreatedBy: req.Username,
		},
	})

	orderShipping, err = s.orderShipping.Upsert(orderShipping)

	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return resp, message.ErrSaveOrderShipping
	}

	return &response.RepickupOrderResponse{
		OrderShippingUID: orderShipping.UID,
		OrderNoAPI:       orderShipping.OrderNoAPI,
		PickupCode:       *orderShipping.PickupCode,
	}, message.SuccessMsg
}

func (s *shippingServiceImpl) repickupOrder(orderShipping *entity.OrderShipping) message.Message {

	if orderShipping.Status == shipping_provider.StatusCancelled {
		return message.OrderHasBeenCancelledMsg
	}

	if orderShipping.Status != shipping_provider.StatusCreated {
		return message.RequestPickupHasBeenMadeMsg
	}

	switch orderShipping.Courier.CourierType {
	case shipping_provider.ThirPartyCourier, shipping_provider.AggregatorCourier:
		return s.repickupThirPartyOrder(orderShipping)
	}
	return message.ErrInvalidCourierType
}

func (s *shippingServiceImpl) repickupThirPartyOrder(orderShipping *entity.OrderShipping) message.Message {
	switch orderShipping.Courier.Code {
	case shipping_provider.ShipperCode:
		result, msg := s.shipper.CreatePickUpOrderWithTimeSlots(orderShipping.BookingID)
		if msg != message.SuccessMsg {
			return msg
		}

		//update pickupCode
		orderShipping.PickupCode = &result.Data.OrderActivation[0].PickUpCode
	default:
		return message.ErrInvalidCourierCode
	}

	return message.SuccessMsg
}

// swagger:operation GET /shipping/tracking/{uid} Shipping ShippingTracking
// Get Order Shipping Tracking (No Auth)
//
// Description :
// Get Order Shipping Tracking (No Auth)
// ---
// security:
// - Bearer: []
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           properties:
//             record:
//               $ref: '#/definitions/GetOrderShippingTrackingResponse'
func (s *shippingServiceImpl) ShippingTracking(req *request.GetOrderShippingTracking) ([]response.GetOrderShippingTracking, message.Message) {
	return s.OrderShippingTracking(req)
}

func updateStatusTopic(channelCode string) string {
	topic := viper.GetString("dapr.topic.update-order-shipping")
	return strings.ReplaceAll(topic, "{channel-code}", strings.ToLower(channelCode))
}

// swagger:operation POST /public/webhook/grab Public WebhookUpdateStatusGrab
// Update Status from Grab Webhook
//
// Description :
//
// ---
//
// responses:
//   '200':
//     description: Success Response.
//     schema:
//       properties:
//         meta:
//            $ref: '#/definitions/MetaResponse'
//         data:
//           type: object
func (s *shippingServiceImpl) UpdateStatusGrab(req *request.WebhookUpdateStatusGrabRequest) message.Message {
	logger := log.With(s.logger, "ShippingService", "UpdateStatusGrab")

	if jsonReq, err := json.Marshal(req); err == nil {
		_ = level.Info(logger).Log("grab_webhook", string(jsonReq))
	}

	_ = level.Info(logger).Log("check_webhook_auth", req.Authorization)
	_ = level.Info(logger).Log("check_webhook_auth_id", req.AuthorizationID)
	if !shipping_provider.GrabWebhookAuth(&request.WebhookUpdateStatusGrabHeader{
		AuthorizationID: req.AuthorizationID,
		Authorization:   req.Authorization,
	}) {
		return message.ErrUnAuth
	}

	orderShipping, err := s.orderShipping.FindByOrderNo(req.Body.MerchantOrderID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByOrderNo", err.Error())
		return message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return message.ErrOrderShippingNotFound
	}

	if orderShipping.Courier.Code != shipping_provider.GrabCode {
		return message.ErrOrderShippingNotFound
	}

	statusCode := req.Body.Status
	statusDescription := req.Body.FailedReason

	shippingStatus, err := s.shippingCourierStatusRepo.FindByCourierStatus(orderShipping.CourierID, fmt.Sprint(statusCode))

	if err != nil {
		_ = level.Error(logger).Log("s.shippingCourierStatusRepo.FindByCourierStatus", err.Error())
		return message.ShippingStatusNotFoundMsg
	}

	if shippingStatus == nil {
		return message.ShippingStatusNotFoundMsg
	}

	driverInfo := request.UpdateOrderShippingDriverInfo{
		Name:         req.Body.Driver.Name,
		Phone:        req.Body.Driver.Phone,
		LicencePlate: req.Body.Driver.LicensePlate,
		TrackingURL:  req.Body.TrackURL,
	}

	orderShipping.Status = shippingStatus.StatusCode
	orderShipping.UpdatedBy = "GRAB_WEBHOOK"
	orderShipping.AddHistoryStatus(shippingStatus, statusDescription, driverInfo.Description())

	orderShipping, err = s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return message.ErrSaveOrderShipping
	}
	topic := updateStatusTopic(orderShipping.Channel.ChannelCode)
	updateOrderRequest := request.UpdateOrderShippingBody{
		ChannelUID:         orderShipping.Channel.UID,
		CourierCode:        orderShipping.Courier.Code,
		CourierServiceUID:  orderShipping.CourierService.UID,
		OrderNo:            orderShipping.OrderNo,
		OrderShippingUID:   orderShipping.UID,
		Airwaybill:         orderShipping.Airwaybill,
		ShippingStatus:     shippingStatus.StatusCode,
		ShippingStatusName: shippingStatus.ShippingStatus.StatusName,
		UpdatedBy:          "shipping_service",
		Timestamp:          time.Now(),
		Details: request.UpdateOrderShippingBodyDetail{
			ExternalStatusCode:        req.Body.Status,
			ExternalStatusName:        req.Body.Status,
			ExternalStatusDescription: req.Body.FailedReason,
		},
		DriverInfo: driverInfo,
	}

	_ = level.Info(logger).Log("PUBLISH_QUEUE, TOPIC", topic)
	s.daprEndpoint.PublishKafka(topic, updateOrderRequest)
	return message.SuccessMsg
}

/*
 Example :

 GRAB = "Paket Anda sudah diterima oleh GRAB. Driver Name : Grabu Duraivu, Driver Phone Number : 6287888889999. Live track di sini : <a href='https://express.grab.com/TESTSANDBOX' target='_blank'>https://express.grab.com/TESTSANDBOX</a>"

 GO-SEND = "Paket Anda sudah diterima oleh GO-SEND. Driver Name : Gojeku Duraivu, Driver Phone Number : +6287888889999"
*/
func SplitDriverInfo(driverInfo string) request.UpdateOrderShippingDriverInfo {
	res := request.UpdateOrderShippingDriverInfo{}
	names := strings.Split(driverInfo, "Driver Name : ")
	if len(names) == 2 {
		res.Name = names[1]
		res.Name = strings.Split(res.Name, ", ")[0]
	}

	phones := strings.Split(driverInfo, "Driver Phone Number : ")
	if len(phones) == 2 {
		res.Phone = phones[1]
		res.Phone = strings.Split(res.Phone, ". ")[0]
	}

	url := strings.Split(driverInfo, "<a href='")
	if len(url) == 2 {
		res.TrackingURL = url[1]
		res.TrackingURL = strings.Split(res.TrackingURL, "' ")[0]
	}

	return res
}
