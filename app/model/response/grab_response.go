package response

import (
	"fmt"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"math"
	"strings"
	"time"
)

type GrabDeliveryQuotesError struct {
	Message    string `json:"message"`
	DevMessage string `json:"devMessage"`
	Arg        string `json:"arg"`
}

func (g *GrabDeliveryQuotesError) GetReason() string {
	args := strings.Split(g.Arg, "Reason: ")
	if len(args) > 1 {
		return args[1]
	}

	return util.ReplaceEmptyString(g.Arg, g.Message)
}

type GrabDeliveryQuotes struct {
	Quotes      []Quote     `json:"quotes"`
	Origin      Origin      `json:"origin"`
	Destination Destination `json:"destination"`
	Packages    []Package   `json:"packages"`
}

func (g *GrabDeliveryQuotes) ToShippingRate() *ShippingRateCommonResponse {
	if g == nil {
		return nil
	}

	now := time.Now().In(time.UTC)
	data := map[string]ShippingRateData{}

	volume := util.CalculateVolume(float64(g.Packages[0].Dimensions.Width),
		float64(g.Packages[0].Dimensions.Height),
		float64(g.Packages[0].Dimensions.Depth))

	volumeWeight := util.CalculateVolumeWeightKg(float64(g.Packages[0].Dimensions.Width),
		float64(g.Packages[0].Dimensions.Height),
		float64(g.Packages[0].Dimensions.Depth))

	for _, v := range g.Quotes {
		courierShippingCode := global.CourierShippingCodeKey("grab", fmt.Sprint(v.Service.Type))
		data[courierShippingCode] = ShippingRateData{
			AvailableCode: 200,
			Error:         SetShippingRateErrorMessage(message.SuccessMsg),
			Volume:        volume,
			VolumeWeight:  volumeWeight,

			//gram to kg
			FinalWeight: math.Max(float64(g.Packages[0].Dimensions.Weight/1000), volumeWeight),
			UnitPrice:   v.Amount,
			TotalPrice:  v.Amount,
			Etd_Min:     float64(v.EstimationTimeline.PickUp.Sub(now).Hours()),
			Etd_Max:     float64(v.EstimationTimeline.DropOff.Sub(now).Hours()),

			//m to km
			Distance: v.Distance / 1000,
		}
	}

	return &ShippingRateCommonResponse{
		Rate:       data,
		Summary:    make(map[string]ShippingRateSummary),
		CourierMsg: make(map[string]message.Message),
	}
}

type Quote struct {
	Service            QuoteService            `json:"service"`
	Currency           QuoteCurrency           `json:"currency"`
	EstimationTimeline QuoteEstimationTimeline `json:"estimatedTimeline"`
	Amount             float64                 `json:"amount"`
	Distance           float64                 `json:"distance"`
}

type QuoteService struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type QuoteCurrency struct {
	Code     string `json:"code"`
	Symbol   string `json:"symbol"`
	Exponent int    `json:"exponent"`
}

type QuoteEstimationTimeline struct {
	PickUp  time.Time `json:"pickup"`
	DropOff time.Time `json:"dropoff"`
}

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Depth  int `json:"depth"`
	Weight int `json:"weight"`
}

type Package struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Quantity    int        `json:"quantity"`
	Price       int        `json:"price"`
	Dimensions  Dimensions `json:"dimensions"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Origin struct {
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}

type Destination struct {
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}
