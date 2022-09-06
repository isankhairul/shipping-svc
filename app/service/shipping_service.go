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
}

type shippingServiceImpl struct {
	logger              log.Logger
	baseRepo            repository.BaseRepository
	channelRepo         repository.ChannelRepository
	courierServiceRepo  repository.CourierServiceRepository
	courierCoverageCode repository.CourierCoverageCodeRepository
	shipper             shipping_provider.Shipper
	redis               cache.RedisCache
}

func NewShippingService(
	l log.Logger,
	br repository.BaseRepository,
	chrp repository.ChannelRepository,
	csrp repository.CourierServiceRepository,
	cccrp repository.CourierCoverageCodeRepository,
	sh shipping_provider.Shipper,
	rc cache.RedisCache,
) ShippingService {
	return &shippingServiceImpl{
		l, br, chrp, csrp, cccrp, sh, rc,
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

	if len(courierIDs) == 0 {
		return price
	}

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
			Error: response.GetShippingRateError{
				Message: message.SuccessMsg.Message,
			},
		}

		if v.CourierTypeCode == shipping_provider.MerchantCourier {
			value.TotalPrice = 0
		}

		if !(originOK && destinationOK) {
			value.AvailableCode = 400
			value.Error.Message = message.ErrCourierCoverageCodeUidNotExist.Message
		}

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
