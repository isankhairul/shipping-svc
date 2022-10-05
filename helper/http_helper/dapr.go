package http_helper

import (
	"encoding/json"
	"go-klikdokter/app/model/request"
	"go-klikdokter/helper/message"
	"strings"

	"github.com/spf13/viper"
)

type DaprEndpoint interface {
	UpdateOrderShipping(req *request.UpdateOrderShipping) message.Message
}

type dapr struct {
}

func NewDaprEndpoint() DaprEndpoint {
	return &dapr{}
}

func (d *dapr) UpdateOrderShipping(req *request.UpdateOrderShipping) message.Message {
	url := viper.GetString("dapr.endpoint.update-order-shipping")
	url = strings.ReplaceAll(url, "{topic-name}", req.TopicName)

	header := map[string]string{"Content-Type": "application/json"}

	respByte, err := Post(url, header, req.Body)

	if err != nil {
		return message.ErrUpdateOrderShipping
	}

	response := map[string]interface{}{}
	err = json.Unmarshal(respByte, &response)

	if err != nil {
		return message.ErrUpdateOrderShipping
	}

	return message.SuccessMsg
}
