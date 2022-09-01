package request

func NewGetPricingDomesticRequest(origin, destination int, input *GetShippingRateRequest) *GetPricingDomestic {
	req := GetPricingDomestic{
		Height: input.TotalHeight,
		Length: input.TotalLength,
		Weight: input.TotalWeight,
		Width:  input.TotalWidth,
		Origin: AreaDetail{
			AreaID:    origin,
			Latitude:  input.Origin.Latitude,
			Longitude: input.Origin.Longitude,
		},
		Destination: AreaDetail{
			AreaID:    destination,
			Latitude:  input.Destination.Latitude,
			Longitude: input.Destination.Longitude,
		},
		Page:      1,
		COD:       false,
		ForOrder:  false,
		ItemValue: input.TotalProductPrice,
	}
	return &req
}

type GetPricingDomestic struct {
	COD         bool       `json:"cod"`
	ForOrder    bool       `json:"for_order"`
	Height      float64    `json:"height"`
	ItemValue   float64    `json:"item_value"`
	Length      float64    `json:"length"`
	Weight      float64    `json:"weight"`
	Width       float64    `json:"width"`
	Page        int        `json:"page"`
	Origin      AreaDetail `json:"origin"`
	Destination AreaDetail `json:"destination"`
}

type AreaDetail struct {
	AreaID    int    `json:"area_id"`
	Latitude  string `json:"lat"`
	Longitude string `json:"long"`
}

type FindShipperCourierCoverage struct {
	CourierID   uint64
	CountryCode string
	PostalCode  string
}
