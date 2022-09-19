package service

import (
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/http_helper/shipping_provider"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/cache"
	"go-klikdokter/pkg/util"
	"math"
	"strconv"

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
	CancelPickup(uid string) message.Message
	CancelOrder(req *request.CancelOrder) message.Message
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
) ShippingService {
	return &shippingServiceImpl{
		l, br, chrp, csrp, cccrp, sh, rc, osr, cr, scs,
	}
}

// swagger:route POST /shipping/shipping-rate/{shipping-type} Shipping ShippingRateByShippingType
// Get Shipping Rate By Shipping Type
//
// responses:
//  200: ShippingRate
func (s *shippingServiceImpl) GetShippingRateByShippingType(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message) {
	if input.ShippingType == "" {
		return nil, message.ErrShippingTypeRequired
	}
	return s.GetShippingRate(input)
}

// swagger:route POST /shipping/shipping-rate Shipping ShippingRate
// Get Shipping Rate
//
// responses:
//  200: ShippingRate
func (s *shippingServiceImpl) GetShippingRate(input request.GetShippingRateRequest) ([]response.GetShippingRateResponse, message.Message) {
	logger := log.With(s.logger, "ShippingService", "GetShippingRate")

	if len(input.CourierServiceUID) == 0 {
		return nil, message.ErrCourierServiceIsRequired
	}

	// find Channel By UID
	channel, err := s.channelRepo.FindByUid(&input.ChannelUID)

	if err != nil {
		_ = level.Error(logger).Log("s.channelRepo.FindByUid", err.Error())
		return nil, message.ErrChannelNotFound
	}

	if channel == nil {
		return nil, message.ErrChannelNotFound
	}

	// find Courier Servies By Channel UID and Courier Servies UID Slice
	courierServices, err := s.courierServiceRepo.FindCourierServiceByChannelAndUIDs(input.ChannelUID, input.CourierServiceUID, input.ContainPrescription, input.ShippingType)

	if err != nil {
		_ = level.Error(logger).Log("s.courierServiceRepo.FindCourierServiceByChannelAndUIDs", err.Error())
		return nil, message.ErrCourierServiceNotFound
	}

	if len(courierServices) == 0 {
		return nil, message.ErrCourierServiceNotFound
	}

	price := s.getAllCourierPrice(courierServices, &input)

	return toGetShippingRateResponseList(courierServices, price), message.SuccessMsg
}

// function to populate price data
func (s *shippingServiceImpl) getAllCourierPrice(courierServices []entity.ChannelCourierServiceForShippingRate, req *request.GetShippingRateRequest) *response.ShippingRateCommonResponse {
	var resp = &response.ShippingRateCommonResponse{
		Rate:    make(map[string]response.ShippingRateData),
		Summary: map[string]response.ShippingRateSummary{},
	}

	// Populate requested courier list distinct
	var courierList []entity.Courier
	courier := make(map[uint64]string)
	for _, v := range courierServices {
		_, ok := courier[v.CourierID]
		if !ok {
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

	resp.Add(internalAndMerchantPrice.Rate)
	resp.Add(thirdPartyPrice.Rate)
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
		}

		if v.CourierTypeCode == shipping_provider.MerchantCourier {
			value.TotalPrice = 0
		}

		value.SetMessage(false, message.SuccessMsg)

		// check if origin or destination not available
		isErr := !(originOK && destinationOK)
		value.SetMessage(isErr, message.ErrCourierCoverageCodeUidNotExist)

		// check if weight exceeds the maximum weight allowed
		isErr = finalWeight > v.MaxWeight && v.MaxWeight > 0
		value.SetMessage(isErr, message.ErrWeightExceeds)

		price.Rate[key] = value
	}

	return price
}

// function to get shipping rate from third party shipping provider
func (s *shippingServiceImpl) getThirdPartyPrice(courier []entity.Courier, input *request.GetShippingRateRequest) *response.ShippingRateCommonResponse {
	logger := log.With(s.logger, "ShippingService", "getAllCourierPrice")
	var resp = &response.ShippingRateCommonResponse{
		Rate:    make(map[string]response.ShippingRateData),
		Summary: map[string]response.ShippingRateSummary{},
	}

	for _, v := range courier {
		var (
			couriePrice *response.ShippingRateCommonResponse
			err         error
			log         string
		)

		if v.CourierType != shipping_provider.ThirPartyCourier {
			continue
		}

		// try to get price data from cache
		key := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%f",
			v.Code,
			input.Origin.PostalCode,
			input.Destination.PostalCode,
			input.Origin.Latitude,
			input.Origin.Longitude,
			input.Destination.Latitude,
			input.Destination.Longitude,
			input.TotalHeight,
		)

		_ = s.redis.GetJsonStruct(key, &couriePrice)
		// if cache doesn't exist
		if couriePrice == nil {
			switch v.Code {
			case shipping_provider.ShipperCode:
				couriePrice, err = s.shipper.GetShippingRate(&v.ID, input)
			default:
				continue
			}

			if err != nil {
				_ = level.Error(logger).Log(log, err.Error())
				continue
			}

			// save price to redis cache
			s.redis.SetJsonStruct(key, couriePrice, viper.GetInt("cache.redis.expired-in-minute.shipping-rate"))
		}

		if couriePrice != nil {
			resp.Add(couriePrice.Rate)
		}
	}

	return resp
}

// function to generate ShippingRateResponseList
func toGetShippingRateResponseList(input []entity.ChannelCourierServiceForShippingRate, price *response.ShippingRateCommonResponse) []response.GetShippingRateResponse {
	shippingTypeMap := make(map[string][]response.GetShippingRateService)
	var resp []response.GetShippingRateResponse

	for _, v := range input {
		courierShippingCode := global.CourierShippingCodeKey(v.CourierCode, v.ShippingCode)
		p := price.FindShippingCode(courierShippingCode)
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
		return nil, nil, nil, message.ErrChannelNotFound
	}

	if channel == nil {
		return nil, nil, nil, message.ErrChannelNotFound
	}

	// find Courier Service By UID
	courierService, err := s.courierServiceRepo.FindCourierService(input.ChannelUID, input.CouirerServiceUID)
	if err != nil {
		_ = level.Error(logger).Log("s.courierServiceRepo.FindCourierService", err.Error())
		return nil, nil, nil, message.ErrCourierServiceNotFound
	}

	if courierService == nil {
		return nil, nil, nil, message.ErrCourierServiceNotFound
	}

	// check if order no already exist with status created
	orderShipping, err := s.orderShipping.FindByOrderNo(input.OrderNo)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByOrderNo", err.Error())
		return nil, nil, nil, message.ErrDB
	}

	if orderShipping != nil && orderShipping.Status != shipping_provider.StatusCreated {
		return nil, nil, nil, message.ErrOrderNoAlreadyExists
	}

	// get shipping status
	shippingStatus, _ := s.shippingCourierStatusRepo.FindByCode(channel.ID, courierService.CourierID, shipping_provider.StatusRequestPickup)
	if shippingStatus == nil {
		return nil, nil, nil, message.ErrShippingStatus
	}

	// if order no doesn't exist create new one
	if orderShipping == nil {
		orderShipping = &entity.OrderShipping{}
		orderShipping.FromCreateDeliveryRequest(input)
		orderShipping.ChannelID = channel.ID
		orderShipping.CourierID = courierService.CourierID
		orderShipping.CourierServiceID = courierService.ID
	}

	return courierService, orderShipping, shippingStatus, message.SuccessMsg
}

// swagger:route POST /shipping/order-shipping Shipping CreateDelivery
// Create Order Shipping
//
// responses:
//  200: CreateDelivery
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
			return nil, msg
		}

		if orderShipping.ID == 0 {
			orderShipping.Insurance = orderData.Insurance
			orderShipping.InsuranceCost = orderData.InsuranceCost
			orderShipping.ShippingCost = orderData.ShippingCost
			orderShipping.TotalShippingCost = orderData.TotalShippingCost
			orderShipping.ActualShippingCost = orderData.ActualShippingCost
			orderShipping.BookingID = orderData.BookingID
		}

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

// swagger:route GET /shipping/order-tracking/{uid} Shipping OrderShippingTracking
// Get Order Shipping Tracking
//
// responses:
//  200: OrderShippingTracking
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

	switch orderShipping.Courier.CourierType {
	case shipping_provider.ThirPartyCourier:
		return s.thridPartyTracking(orderShipping)
	}

	return nil, message.ErrInvalidCourierType
}

func (s *shippingServiceImpl) thridPartyTracking(orderShipping *entity.OrderShipping) ([]response.GetOrderShippingTracking, message.Message) {
	switch orderShipping.Courier.Code {
	case shipping_provider.ShipperCode:
		return s.shipper.GetTracking(orderShipping.BookingID)
	}

	return nil, message.ErrInvalidCourierCode
}

// swagger:route POST /shipping/webhook/shipper Shipping WebhookUpdateStatusShipper
// Update Status Shipper
//
// responses:
//  200: SuccessResponse
func (s *shippingServiceImpl) UpdateStatusShipper(req *request.WebhookUpdateStatusShipper) (*entity.OrderShipping, message.Message) {
	logger := log.With(s.logger, "ShippingService", "UpdateStatusShipper")

	orderShipping, err := s.orderShipping.FindByOrderNo(req.ExternalID)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.FindByOrderNo", err.Error())
		return nil, message.ErrOrderShippingNotFound
	}

	if orderShipping == nil {
		return nil, message.ErrOrderShippingNotFound
	}

	shippingStatus, err := s.shippingCourierStatusRepo.FindByCourierStatus(orderShipping.CourierID, fmt.Sprint(req.ExternalStatus.Code))

	if err != nil {
		_ = level.Error(logger).Log("s.shippingCourierStatusRepo.FindByCourierStatus", err.Error())
		return nil, message.ErrShippingStatus
	}

	if shippingStatus == nil {
		return nil, message.ErrShippingStatus
	}

	orderShipping.Status = shippingStatus.StatusCode

	if len(req.Awb) > 0 {
		orderShipping.Airwaybill = req.Awb
	}

	orderShipping.AddHistoryStatus(shippingStatus, req.ExternalStatus.Description)

	orderShipping, err = s.orderShipping.Upsert(orderShipping)
	if err != nil {
		_ = level.Error(logger).Log("s.orderShipping.Upsert", err.Error())
		return nil, message.ErrSaveOrderShipping
	}

	return orderShipping, message.SuccessMsg
}

// swagger:route GET /shipping/order-shipping Shipping GetOrderShippingList
// Get Order Shipping List
//
// responses:
//  200: GetOrderShippingList
func (s *shippingServiceImpl) GetOrderShippingList(req *request.GetOrderShippingList) ([]response.GetOrderShippingList, *base.Pagination, message.Message) {
	logger := log.With(s.logger, "ShippingService", "GetOrderShippingList")

	if len(req.Filters.OrderShippingDateFrom) > 0 {
		if ok := util.DateValidationYYYYMMDD(req.Filters.OrderShippingDateFrom); !ok {
			return nil, nil, message.ErrFormatDateYYYYMMDD
		}
	}

	if len(req.Filters.OrderShippingDateTo) > 0 {
		if ok := util.DateValidationYYYYMMDD(req.Filters.OrderShippingDateTo); !ok {
			return nil, nil, message.ErrFormatDateYYYYMMDD
		}
	}

	filter := make(map[string]interface{})
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

// swagger:route GET /shipping/order-shipping/{uid} Shipping GetOrderShippingDetail
// Get Order Shipping Detail By UID
//
// responses:
//  200: GetOrderShippingDetail

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
	resp.ChannelCode = orderShipping.Channel.ChannelCode
	resp.ChannelName = orderShipping.Channel.ChannelName
	resp.CourierName = orderShipping.Courier.CourierName
	resp.CourierServiceName = orderShipping.CourierService.ShippingName
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
			CreatedAt: v.CreatedAt,
			Status:    v.StatusCode,
			Notes:     v.Note,
		})
	}

	return resp
}

// swagger:route POST /shipping/cancel-pickup/{uid} Shipping CancelPickup
// Cancel Pickup Order
//
// responses:
//  200: SuccessResponse
func (s *shippingServiceImpl) CancelPickup(uid string) message.Message {
	logger := log.With(s.logger, "ShippingService", "CancelPickup")
	orderShipping, err := s.orderShipping.FindByUID(uid)
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
		return message.ErrShippingStatus
	}

	orderShipping.AddHistoryStatus(shipperStatus, "")
	orderShipping.Status = shipping_provider.StatusCancelled
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

// swagger:route POST /shipping/cancel-order/{uid} Shipping CancelOrder
// Cancel Order Shipping
//
// responses:
//  200: SuccessResponse
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
		return message.ErrShippingStatus
	}

	orderShipping.AddHistoryStatus(shipperStatus, req.Body.Reason)
	orderShipping.Status = shipping_provider.StatusCancelled
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
