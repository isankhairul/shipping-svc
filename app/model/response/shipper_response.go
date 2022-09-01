package response

import (
	"fmt"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
)

type ShipperMetaData struct {
	Path           string `json:"path"`
	HTTPStatusCode int    `json:"http_status_code"`
	HTTPStatus     string `json:"http_status"`
	Timestamp      uint64 `json:"timestamp"`
}

type ShipperPagination struct {
	CurrentPage     int      `json:"current_page"`
	CurrentElements int      `json:"current_elements"`
	TotalPages      int      `json:"total_pages"`
	TotalElements   int      `json:"total_elements"`
	SortBy          []string `json:"sort_by"`
}

type GetPricingDomesticData struct {
	Origin      DataAreaDetail `json:"origin"`
	Destination DataAreaDetail `json:"destination"`
	Pricings    []PricingsItem `json:"pricings"`
}

type DataAreaDetail struct {
	// "area_id": 2482,
	// "area_name": "Sampora",
	// "suburb_id": 240,
	// "suburb_name": "Cisauk",
	// "city_id": 21,
	// "city_name": "Tangerang, Kab.",
	// "province_id": 3,
	// "province_name": "Banten",
	// "country_id": 228,
	// "country_name": "INDONESIA",
	// "lat": -6.3124856,
	// "lng": 106.6497623
	AreaID    int     `json:"area_id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type PricingsItem struct {
	Logistic         PricingLogisticDetail `json:"logistic"`
	Rate             PricingRateDetail     `json:"rate"`
	Weight           float64               `json:"weight"`
	Volume           float64               `json:"volume"`
	VolumeWeight     float64               `json:"volume_weight"`
	FinalWeight      float64               `json:"final_weight"`
	MinDay           int                   `json:"min_day"`
	MaxDay           int                   `json:"max_day"`
	UnitPrice        float64               `json:"unit_price"`
	TotalPrice       float64               `json:"total_price"`
	Discount         float64               `json:"discount"`
	DiscountValue    float64               `json:"discount_value"`
	DiscountedPrice  float64               `json:"discounted_price"`
	InsuranceFee     float64               `json:"insurance_fee"`
	MustUseInsurance bool                  `json:"must_use_insurance"`
	LiabilityValue   float64               `json:"liability_value"`
	FinalPrice       float64               `json:"final_price"`
	Currency         string                `json:"currency"`
	InsuranceApplied bool                  `json:"insurance_applied"`
	BasePrice        float64               `json:"base_price"`
	SurchargeFee     float64               `json:"surcharge_fee"`
}

type PricingLogisticDetail struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	Code        string `json:"code"`
	CompanyName string `json:"company_name"`
}

type PricingRateDetail struct {
	ID              uint64 `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Description     string `json:"description"`
	FullDescription string `json:"full_description"`
	IsHubless       bool   `json:"is_hubless"`
}

type GetPricingDomestic struct {
	Metadata   ShipperMetaData        `json:"metadata"`
	Data       GetPricingDomesticData `json:"data"`
	Pagination ShipperPagination      `json:"pagination"`
}

func (g *GetPricingDomestic) ToShippingRate() *ShippingRateCommonResponse {
	if g == nil {
		return nil
	}

	data := map[string]ShippingRateData{}

	for _, v := range g.Data.Pricings {
		courierShippingCode := global.CourierShippingCodeKey("shipper", fmt.Sprint(v.Rate.ID))
		data[courierShippingCode] = ShippingRateData{
			AvailableCode:    200,
			Error:            GetShippingRateError{},
			Weight:           v.Weight,
			Volume:           v.Volume,
			VolumeWeight:     v.VolumeWeight,
			FinalWeight:      v.FinalWeight,
			MinDay:           v.MinDay,
			MaxDay:           v.MaxDay,
			UnitPrice:        v.UnitPrice,
			TotalPrice:       v.TotalPrice,
			InsuranceFee:     v.InsuranceFee,
			MustUseInsurance: v.MustUseInsurance,
			InsuranceApplied: v.InsuranceApplied,
			Distance: global.DistanceKM(g.Data.Origin.Latitude,
				g.Data.Origin.Longitude,
				g.Data.Destination.Latitude,
				g.Data.Destination.Longitude),
		}
	}

	return &ShippingRateCommonResponse{
		Rate:    data,
		Summary: make(map[string]ShippingRateSummary),
		Msg:     message.SuccessMsg,
	}
}
