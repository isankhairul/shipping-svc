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
	UpdateOrderShipping(req *request.UpdateOrderShipping) (*response.UpdateOrderShippingResponse, message.Message)
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
) ShippingService {
	return &shippingServiceImpl{
		l, br, chrp, csrp, cccrp, sh, rc, osr, cr, scs, de,
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
//               $ref: '#/definitions/ShippmentPredefinedDetail'
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
		volume       = util.CalculateVolume(req.TotalHeight, req.TotalLength, req.TotalLength)
		volumeWeight = util.CalculateVolumeWeightKg(req.TotalHeight, req.TotalLength, req.TotalLength)
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

	for _, v := range courier {
		var (
			courierPrice *response.ShippingRateCommonResponse
			err          error
		)

		if v.CourierType != shipping_provider.ThirPartyCourier {
			continue
		}

		// try to get price data from cache
		key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s:%s:%f",
			v.Code,
			input.Origin.PostalCode,
			input.Destination.PostalCode,
			input.Origin.Subdistrict,
			input.Destination.Subdistrict,
			input.Origin.Latitude,
			input.Origin.Longitude,
			input.Destination.Latitude,
			input.Destination.Longitude,
			input.TotalHeight,
		)

		_ = s.redis.GetJsonStruct(key, &courierPrice)
		// if cache doesn't exist
		if courierPrice == nil {
			switch v.Code {
			case shipping_provider.ShipperCode:
				courierPrice, err = s.shipper.GetShippingRate(&v.ID, input)
			default:
				resp.CourierMsg[v.Code] = message.InvalidCourierCodeMsg
				continue
			}

			if err == nil {
				// save price to redis cache
				s.redis.SetJsonStruct(key, courierPrice, viper.GetInt("cache.redis.expired-in-minute.shipping-rate"))
			}

		}

		if courierPrice != nil {
			resp.Add(courierPrice)
		}
	}

	return resp
}

// function to generate ShippingRateResponseList
func toGetShippingRateResponseList(req *request.GetShippingRateRequest, courierServices []entity.ChannelCourierServiceForShippingRate, price *response.ShippingRateCommonResponse) []response.GetShippingRateResponse {
	shippingTypeMap := make(map[string][]response.GetShippingRateService)
	var resp []response.GetShippingRateResponse

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
			Etd_Min:                 v.EtdMin,
			Etd_Max:                 v.EtdMax,
			AvailableCode:           p.AvailableCode,
			Error:                   p.Error,
			Weight:                  p.Weight,
			Volume:                  p.Volume,
			VolumeWeight:            p.VolumeWeight,
			FinalWeight:             p.FinalWeight,
			MinDay:                  p.MinDay,
			MaxDay:                  p.MaxDay,
			UnitPrice:               p.UnitPrice,
			TotalPrice:              p.TotalPrice,
			InsuranceFee:            p.InsuranceFee,
			MustUseInsurance:        p.MustUseInsurance,
			InsuranceApplied:        p.InsuranceApplied,
			Distance:                p.Distance,
		}

		price.SummaryPerShippingType(v.ShippingTypeCode, p.TotalPrice, v.EtdMax, v.EtdMin)

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
			AvailableCode:           200,
			Error:                   response.GetShippingRateError{},
		}
		resp = append(resp, data)
	}

	return resp
}

func (s *shippingServiceImpl) PopulateCreateDelivery(input *request.CreateDelivery) (*entity.CourierService, *entity.OrderShipping, *entity.ShippingCourierStatus, message.Message) {
	logger := log.With(s.logger, "ShippingService", "PopulateCreateDelivery")

	// find Channel By UID
	channel, err := s.channelRepo.FindByUid(&input.ChannelUID)
	if err != nil {
		_ = level.Error(logger).Log("s.channelRepo.FindByUid", err.Error())
		return nil, nil, nil, message.ChannelNotFoundMsg
	}

	if channel == nil {
		return nil, nil, nil, message.ChannelNotFoundMsg
	}

	// find Courier Service By UID
	courierService, err := s.courierServiceRepo.FindCourierService(input.ChannelUID, input.CouirerServiceUID)
	if err != nil {
		_ = level.Error(logger).Log("s.courierServiceRepo.FindCourierService", err.Error())
		return nil, nil, nil, message.CourierServiceNotFoundMsg
	}

	if courierService == nil {
		return nil, nil, nil, message.CourierServiceNotFoundMsg
	}

	if msg := courierService.Validate(input.Package.TotalWeight, input.Package.ContainPrescription > 0); msg != message.SuccessMsg {
		return nil, nil, nil, msg
	}

	// check if order no already exist with status created
	orderShipping, err := s.orderShipping.FindByOrderNo(input.OrderNo)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByOrderNo", err.Error())
		return nil, nil, nil, message.ErrDB
	}

	if orderShipping != nil && orderShipping.Status != shipping_provider.StatusCreated {
		return nil, nil, nil, message.OrderNoAlreadyExistsMsg
	}

	// get shipping status
	shippingStatus, _ := s.shippingCourierStatusRepo.FindByCode(channel.ID, courierService.CourierID, shipping_provider.StatusRequestPickup)
	if shippingStatus == nil {
		return nil, nil, nil, message.ShippingStatusNotFoundMsg
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
	return courierService, orderShipping, shippingStatus, message.SuccessMsg
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
//               $ref: '#/definitions/ShippmentPredefinedDetail'
func (s *shippingServiceImpl) CreateDelivery(input *request.CreateDelivery) (*response.CreateDelivery, message.Message) {
	logger := log.With(s.logger, "ShippingService", "CreateDelivery")
	var resp *response.CreateDelivery

	courierService, orderShipping, shippingStatus, msg := s.PopulateCreateDelivery(input)
	if msg != message.SuccessMsg {
		return nil, msg
	}

	switch courierService.Courier.CourierType {
	case shipping_provider.ThirPartyCourier:

		orderData, msg := s.createDeliveryThirdParty(orderShipping.BookingID, courierService, input)
		if msg != message.SuccessMsg {
			return &response.CreateDelivery{}, msg
		}

		if orderShipping.ID == 0 {
			orderShipping.Insurance = orderData.Insurance
			orderShipping.InsuranceCost = orderData.InsuranceCost
			orderShipping.ShippingCost = orderData.ShippingCost
			orderShipping.TotalShippingCost = orderData.TotalShippingCost
			orderShipping.ActualShippingCost = orderData.ActualShippingCost
			orderShipping.BookingID = orderData.BookingID
		}

		orderShipping.PickupCode = orderData.PickUpCode
		orderShipping.Status = orderData.Status

	case shipping_provider.InternalCourier:
		return resp, message.ErrInvalidCourierType
	case shipping_provider.MerchantCourier:
		return resp, message.ErrInvalidCourierType
	default:
		return resp, message.ErrInvalidCourierType
	}

	if orderShipping.Status != shipping_provider.StatusCreated {
		orderShipping.AddHistoryStatus(shippingStatus, "")
	}

	orderShipping, err := s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return nil, message.ErrSaveOrderShipping
	}

	return &response.CreateDelivery{
		OrderNoAPI:       input.OrderNo,
		OrderShippingUID: orderShipping.UID,
	}, message.SuccessMsg
}

func (s *shippingServiceImpl) createDeliveryThirdParty(bookingID string, courierService *entity.CourierService, input *request.CreateDelivery) (*response.CreateDeliveryThirdPartyData, message.Message) {
	switch courierService.Courier.Code {
	case shipping_provider.ShipperCode:
		return s.shipper.CreateDelivery(bookingID, courierService, input)
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
//               $ref: '#/definitions/ShippmentPredefinedDetail'
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
	case shipping_provider.ThirPartyCourier:
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
	}

	return nil, message.ErrInvalidCourierCode
}

// swagger:operation POST /shipping/webhook/shipper Shipping WebhookUpdateStatusShipper
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

	orderShipping.AddHistoryStatus(shippingStatus, statusDescription)

	orderShipping, err = s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return nil, message.ErrSaveOrderShipping
	}

	topic := "queueing.shipment.order_shipping_update." + orderShipping.Channel.ChannelCode
	updateOrderRequest := request.UpdateOrderShipping{
		TopicName: topic,
		Body: request.UpdateOrderShippingBody{
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
		},
	}
	s.daprEndpoint.UpdateOrderShipping(&updateOrderRequest)
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
	filter["order_no"] = req.Filters.OrderNo
	filter["channel_code"] = req.Filters.ChannelCode
	filter["channel_name"] = req.Filters.ChannelName
	filter["courier_name"] = req.Filters.CourierName
	filter["shipping_status"] = req.Filters.ShippingStatus
	filter["order_shipping_date_from"] = req.Filters.OrderShippingDateFrom
	filter["order_shipping_date_to"] = req.Filters.OrderShippingDateTo

	result, pagination, err := s.orderShipping.FindByParams(req.Limit, req.Page, req.Sort, req.Dir, filter)
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
	resp.MerchantDistrictName = orderShipping.MerchantDistrictCode
	resp.MerchantCityName = orderShipping.MerchantCityCode
	resp.MerchantProvinceName = orderShipping.MerchantProvinceCode
	resp.MerchantPostalCode = orderShipping.MerchantPostalCode
	resp.CustomerUID = orderShipping.CustomerUID
	resp.CustomerName = orderShipping.CustomerName
	resp.CustomerEmail = orderShipping.CustomerEmail
	resp.CustomerPhone = orderShipping.CustomerPhoneNumber
	resp.CustomerAddress = orderShipping.CustomerAddress
	resp.CustomerDistrictName = orderShipping.CustomerDistrictCode
	resp.CustomerCityName = orderShipping.CustomerCityCode
	resp.CustomerProvinceName = orderShipping.CustomerProvinceCode
	resp.CustomerPostalCode = orderShipping.CustomerPostalCode
	resp.CustomerNotes = orderShipping.CustomerNotes
	resp.OrderShippingItem = []response.GetOrderShippingDetailItem{}
	resp.OrderShippingHistory = []response.GetOrderShippingDetailHistory{}

	for _, v := range orderShipping.OrderShippingItem {
		resp.OrderShippingItem = append(resp.OrderShippingItem, response.GetOrderShippingDetailItem{
			ItemName:    v.ItemName,
			ProductUID:  v.ProductUID,
			Qty:         v.Quantity,
			Price:       v.Price,
			Weight:      v.Weight,
			Volume:      v.Volume,
			Prescrition: v.Prescription,
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

	shipperStatus, _ := s.shippingCourierStatusRepo.FindByCode(orderShipping.ChannelID, orderShipping.CourierID, shipping_provider.StatusCancelled)

	if shipperStatus == nil {
		return message.ShippingStatusNotFoundMsg
	}

	orderShipping.AddHistoryStatus(shipperStatus, "")
	orderShipping.Status = shipping_provider.StatusCancelled
	orderShipping.UpdatedBy = req.Body.Username
	_, err = s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return message.ErrSaveOrderShipping
	}

	return message.SuccessMsg
}

func (s *shippingServiceImpl) cancelPickup(orderShipping *entity.OrderShipping) message.Message {
	switch orderShipping.Courier.CourierType {
	case shipping_provider.ThirPartyCourier:
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
		_, err = s.shipper.CancelPickupRequest(orderShipping.PickupCode)
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
	case shipping_provider.ThirPartyCourier:
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

// swagger:operation POST /shipping/update-order-shipping/{topic-name} Shipping UpdateOrderShipping
// Update Order Shipping
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
//               $ref: '#/definitions/UpdateOrderShipping'
func (s *shippingServiceImpl) UpdateOrderShipping(req *request.UpdateOrderShipping) (*response.UpdateOrderShippingResponse, message.Message) {
	return &response.UpdateOrderShippingResponse{
		OrderShippingUID: req.Body.OrderShippingUID,
		OrderNoAPI:       req.Body.OrderNo,
		ShippingStatus:   req.Body.ShippingStatus,
	}, message.SuccessMsg
}
