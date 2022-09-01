package response

import (
	"go-klikdokter/app/model/entity"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util/datatype"
)

type GetShippingRatePriceRange struct {
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
}

type GetShippingRateError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type GetShippingRateCourir struct {
	CourierID       uint64 `json:"-"`
	CourierUID      string `json:"courier_uid"`
	CourierCode     string `json:"courier_code"`
	CourierName     string `json:"courier_name"`
	CourierTypeCode string `json:"courier_type_code"`
	CourierTypeName string `json:"courier_type_name"`
}

type GetShippingRateService struct {
	Courier                 GetShippingRateCourir `json:"courier"`
	CourierServiceUID       string                `json:"courier_service_uid"`
	ShippingCode            string                `json:"shipping_code"`
	ShippingName            string                `json:"shipping_name"`
	ShippingTypeDescription string                `json:"shipping_type_description"`
	Logo                    datatype.JSONB        `json:"logo"`
	ShippingTypeCode        string                `json:"shipping_type_code"`
	ShippingTypeName        string                `json:"shipping_type_name"`
	Etd_Min                 float64               `json:"etd_min"`
	Etd_Max                 float64               `json:"etd_max"`
	AvailableCode           int                   `json:"available_code"`
	Error                   GetShippingRateError  `json:"error"`
	Weight                  float64               `json:"weight"`
	Volume                  float64               `json:"volume"`
	VolumeWeight            float64               `json:"volume_weight"`
	FinalWeight             float64               `json:"final_weight"`
	MinDay                  int                   `json:"min_day"`
	MaxDay                  int                   `json:"max_day"`
	UnitPrice               float64               `json:"unit_price"`
	TotalPrice              float64               `json:"total_price"`
	InsuranceFee            float64               `json:"insurance_fee"`
	MustUseInsurance        bool                  `json:"must_use_insurance"`
	InsuranceApplied        bool                  `json:"insurance_applied"`
	Distance                float64               `json:"distance"`
}

func (g *GetShippingRateService) FromShipper(val PricingsItem) {
	g.Weight = val.Weight
	g.Volume = val.Volume
	g.VolumeWeight = val.VolumeWeight
	g.FinalWeight = val.FinalWeight
	g.MinDay = val.MinDay
	g.MaxDay = val.MaxDay
	g.UnitPrice = val.UnitPrice
	g.TotalPrice = val.TotalPrice
	g.InsuranceFee = val.InsuranceFee
	g.MustUseInsurance = val.MustUseInsurance
	g.InsuranceApplied = val.InsuranceApplied
}

type GetShippingRateResponse struct {
	ShippingTypeCode        string                    `json:"shipping_type_code"`
	ShippingTypeName        string                    `json:"shipping_type_name"`
	ShippingTypeDescription string                    `json:"shipping_type_description"`
	PriceRange              GetShippingRatePriceRange `json:"price_range"`
	EtdMin                  *float64                  `json:"etd_min"`
	EtdMax                  *float64                  `json:"etd_max"`
	AvailableCode           int                       `json:"available_code"`
	Error                   GetShippingRateError      `json:"error"`
	Services                []GetShippingRateService  `json:"services"`
}

func ToGetShippingRateResponseList(input []entity.ChannelCourierServiceForShippingRate, price *ShippingRateCommonResponse) []GetShippingRateResponse {
	shippingTypeMap := make(map[string][]GetShippingRateService)
	var resp []GetShippingRateResponse

	for _, v := range input {
		courierShippingCode := global.CourierShippingCodeKey(v.CourierCode, v.ShippingCode)
		p := price.FindShippingCode(courierShippingCode)
		service := GetShippingRateService{
			Courier: GetShippingRateCourir{
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

			AvailableCode:    p.AvailableCode,
			Error:            p.Error,
			Weight:           p.Weight,
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
		}

		price.SummaryPerShippingType(v.ShippingTypeCode, p.TotalPrice, v.EtdMax, v.EtdMin)

		if _, ok := shippingTypeMap[v.ShippingTypeCode]; !ok {
			shippingTypeMap[v.ShippingTypeCode] = []GetShippingRateService{}
		}

		shippingTypeMap[v.ShippingTypeCode] = append(shippingTypeMap[v.ShippingTypeCode], service)
	}

	for k, v := range shippingTypeMap {
		s := price.Summary[k]
		data := GetShippingRateResponse{
			ShippingTypeCode:        k,
			ShippingTypeName:        v[0].ShippingTypeName,
			ShippingTypeDescription: v[0].ShippingTypeDescription,
			PriceRange:              s.PriceRange,
			EtdMax:                  s.EtdMax,
			EtdMin:                  s.EtdMin,
			Services:                v,
			AvailableCode:           200,
			Error:                   GetShippingRateError{},
		}
		resp = append(resp, data)
	}

	return resp
}

//swagger:response ShippingRate
type GetShippingRateResponseList struct {
	//in: body
	Response []GetShippingRateResponse `json:"response"`
}

type ShippingRateCommonResponse struct {
	Rate    map[string]ShippingRateData
	Summary map[string]ShippingRateSummary
	Msg     message.Message
}

func (s *ShippingRateCommonResponse) Add(data map[string]ShippingRateData) {
	for k, v := range data {
		s.Rate[k] = v
	}
}

func (s *ShippingRateCommonResponse) FindShippingCode(shippingCode string) ShippingRateData {
	data, ok := s.Rate[shippingCode]

	if s.Msg.Code != message.SuccessMsg.Code && s.Msg.Code != 0 {
		return ShippingRateData{
			AvailableCode: 400,
			Error: GetShippingRateError{
				Message: s.Msg.Message,
			},
		}
	}

	if !ok {
		return ShippingRateData{
			AvailableCode: 400,
			Error: GetShippingRateError{
				Message: message.ErrShippingRateNotFound.Message,
			},
		}
	}
	return data
}

func (s *ShippingRateCommonResponse) SummaryPerShippingType(shippingType string, price, etdMax, etdMin float64) {
	summaryData, ok := s.Summary[shippingType]

	if !ok {
		summaryData = ShippingRateSummary{
			PriceRange: GetShippingRatePriceRange{},
		}
	}

	maxPrice := summaryData.PriceRange.MaxPrice
	minPrice := summaryData.PriceRange.MinPrice
	eMax := summaryData.EtdMax
	eMin := summaryData.EtdMin

	if maxPrice == nil || price > *maxPrice {
		summaryData.PriceRange.MaxPrice = &price
	}

	if minPrice == nil || price < *minPrice {
		summaryData.PriceRange.MinPrice = &price
	}

	if eMax == nil || etdMax > *eMax {
		summaryData.EtdMax = &etdMax
	}

	if eMin == nil || etdMin < *eMin {
		summaryData.EtdMin = &etdMin
	}

	s.Summary[shippingType] = summaryData
}

type ShippingRateData struct {
	AvailableCode    int
	Error            GetShippingRateError
	Weight           float64
	Volume           float64
	VolumeWeight     float64
	FinalWeight      float64
	MinDay           int
	MaxDay           int
	UnitPrice        float64
	TotalPrice       float64
	InsuranceFee     float64
	MustUseInsurance bool
	InsuranceApplied bool
	Distance         float64
}

type ShippingRateSummary struct {
	PriceRange GetShippingRatePriceRange
	EtdMin     *float64
	EtdMax     *float64
}
