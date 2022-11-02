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

type GrabError struct {
	Message    string `json:"message"`
	DevMessage string `json:"devMessage"`
	Arg        string `json:"arg"`
}

func (g *GrabError) GetReason() string {
	args := strings.Split(g.Arg, "Reason: ")
	if len(args) > 1 {
		return args[1]
	}

	msg := util.ReplaceEmptyString(g.Arg, g.Message)
	msg = util.ReplaceEmptyString(msg, g.DevMessage)

	return msg
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

	//for create order
	Packages    []Packages  `json:"packages"`
	Origin      Origin      `json:"origin"`
	Destination Destination `json:"destination"`
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

	// for create order
	Keywords string `json:"keywords"`
	Extra    Extra  `json:"extra"`
}

type Destination struct {
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`

	// for create order
	Keywords string `json:"keywords"`
	Extra    Extra  `json:"extra"`
}

type CreateDeliveryGrab struct {
	DeliveryID  string              `json:"deliveryID"`
	Quote       Quote               `json:"quote"`
	Sender      GrabSenderRecipient `json:"sender"`
	Recipient   GrabSenderRecipient `json:"recipient"`
	PickupPin   string              `json:"pickupPin"`
	Status      string              `json:"status"`
	Courier     Courier             `json:"courier"`
	Timeline    Timeline            `json:"timeline"`
	TrackingURL string              `json:"trackingURL"`
	AdvanceInfo AdvanceInfo         `json:"advanceInfo"`
}
type Service struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}
type Currency struct {
	Code     string `json:"code"`
	Symbol   string `json:"symbol"`
	Exponent int    `json:"exponent"`
}
type EstimatedTimeline struct {
	Create   string `json:"create"`
	Allocate string `json:"allocate"`
	Pickup   string `json:"pickup"`
	Dropoff  string `json:"dropoff"`
	Cancel   string `json:"cancel"`
	Return   string `json:"return"`
}

type Packages struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Quantity    int        `json:"quantity"`
	Price       int        `json:"price"`
	Dimensions  Dimensions `json:"dimensions"`
}

type Extra struct {
}

type GrabSenderRecipient struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Title       string `json:"title"`
	CompanyName string `json:"companyName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	SmsEnabled  bool   `json:"smsEnabled"`
	Instruction string `json:"instruction"`
}
type Vehicle struct {
	PlateNumber string `json:"plateNumber"`
	Model       string `json:"model"`
}
type Courier struct {
	Coordinates Coordinates `json:"coordinates"`
	Name        string      `json:"name"`
	Phone       string      `json:"phone"`
	PictureURL  string      `json:"pictureURL"`
	Vehicle     Vehicle     `json:"vehicle"`
}
type Timeline struct {
	Create   string `json:"create"`
	Allocate string `json:"allocate"`
	Pickup   string `json:"pickup"`
	Dropoff  string `json:"dropoff"`
	Cancel   string `json:"cancel"`
	Return   string `json:"return"`
}
type AdvanceInfo struct {
	FailedReason string `json:"failedReason"`
}

type GrabDeliveryDetail struct {
	DeliveryID      string               `json:"deliveryID"`
	MerchantOrderID string               `json:"merchantOrderID"`
	PaymentMethod   string               `json:"paymentMethod"`
	Quote           Quote                `json:"quote"`
	Sender          GrabSenderRecipient  `json:"sender"`
	Recipient       GrabSenderRecipient  `json:"recipient"`
	Status          string               `json:"status"`
	TrackingURL     string               `json:"trackingURL"`
	Courier         interface{}          `json:"courier"`
	Timeline        map[string]time.Time `json:"timeline"`
	Schedule        interface{}          `json:"schedule"`
	//CashOnDelivery  CashOnDelivery       `json:"cashOnDelivery"`
	InvoiceNo   string      `json:"invoiceNo"`
	PickupPin   string      `json:"pickupPin"`
	AdvanceInfo AdvanceInfo `json:"advanceInfo"`
}

func (g *GrabDeliveryDetail) ToOrderShippingTracking() []GetOrderShippingTracking {
	resp := []GetOrderShippingTracking{}
	for k, v := range g.Timeline {

		// failed notes
		status := strings.ToUpper(k)
		note := status

		if strings.Contains(status, "FAILED") {
			note = fmt.Sprintf("%s %s", note, g.AdvanceInfo.FailedReason)
		}
		resp = append(resp, GetOrderShippingTracking{
			DateTime: v,
			Status:   status,
			Note:     note,
			Date:     v.In(util.Loc).Format(util.LayoutDateOnly),
			Time:     v.In(util.Loc).Format(util.LayoutTimeOnly),
		})
	}

	// order status still QUEUEING
	if len(resp) == 0 {
		resp = append(resp, GetOrderShippingTracking{
			DateTime: g.Quote.EstimationTimeline.PickUp,
			Status:   strings.ToUpper(g.Status),
			Note:     "",
			Date:     g.Quote.EstimationTimeline.PickUp.In(util.Loc).Format(util.LayoutDateOnly),
			Time:     g.Quote.EstimationTimeline.PickUp.In(util.Loc).Format(util.LayoutTimeOnly),
		})
	}
	return resp
}
