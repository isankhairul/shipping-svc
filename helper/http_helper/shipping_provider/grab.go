package shipping_provider

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/response"
	"go-klikdokter/helper/http_helper"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util"
	"math"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

type Grab interface {
	GetShippingRate(input *request.GetShippingRateRequest) (*response.ShippingRateCommonResponse, error)
}

type grab struct {
	Logger log.Logger
	Base   string
	SafeID string
	Secret string
}

func NewGrab(log log.Logger) Grab {
	return &grab{
		Logger: log,
	}
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
	endpoint := viper.GetString("grab.path.get-delivery-quote")
	url := grabUrl(endpoint)

	grabDate := grabDateTimeFormat()
	reqByte, _ := json.Marshal(req)
	encodedRequest := encodeContent(string(reqByte))
	auth := grabAuthentication("POST", grabDate, "application/json", endpoint, encodedRequest)

	headers := map[string]string{
		"Date":          grabDate,
		"Content-Type":  "application/json",
		"Authorization": auth,
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

	errResp := &response.GrabDeliveryQuotesError{}
	err = json.Unmarshal(respByte, &errResp)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(errResp.GetReason())
}

func grabUrl(path string) string {
	base := viper.GetString("grab.base")
	return base + path
}

func grabDateTimeFormat() string {
	const (
		D  = "Mon"
		DD = "Monday"
		d  = "2"
		dd = "02"
		m  = "1"
		mm = "01"
		M  = "Jan"
		MM = "January"
		y  = "06"
		Y  = "2006"
		h  = "03"
		H  = "15"
		i  = "04"
		s  = "05"
	)

	return time.Now().UTC().Format(
		fmt.Sprintf("%s, %s %s %s %s:%s:%s GMT",
			D, dd, M, Y, H, i, s),
	)
}

func grabAuthentication(method string, date string, contentType string, endpoint string, encoded string) string {
	secret := viper.GetString("grab.auth.secret")
	safeID := viper.GetString("grab.auth.safe-id")

	stringToSign := getStringToSign(method, date, contentType, endpoint, encoded)
	hmacSignature := encodeAuthSignature(secret, stringToSign)

	return fmt.Sprintf("%s:%s", safeID, hmacSignature)
}

func encodeAuthSignature(secret string, stringToSign string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return encoded
}

func getStringToSign(method string, date string, ctype string, ep string, encoded string) string {
	sts := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n",
		method, ctype, date, ep, encoded,
	)

	return sts
}

func encodeContent(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return encoded
}
