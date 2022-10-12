package http_helper

import (
	"encoding/json"
	"github.com/go-kit/log/level"
	"go-klikdokter/app/model/request"
	"go-klikdokter/helper/message"
	"strings"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

type DaprEndpoint interface {
	UpdateOrderShipping(req *request.UpdateOrderShipping) message.Message
}

type dapr struct {
	Logger log.Logger
}

func NewDaprEndpoint(log log.Logger) DaprEndpoint {
	return &dapr{log}
}

func (d *dapr) UpdateOrderShipping(req *request.UpdateOrderShipping) message.Message {
	url := viper.GetString("dapr.endpoint.update-order-shipping")
	url = strings.ReplaceAll(url, "{topic-name}", req.TopicName)
	header := map[string]string{"Content-Type": "application/json"}
	logger := log.With(d.Logger, "Webhook", "UpdateOrderShipping")

	_ = level.Info(logger).Log("d.UpdateOrderShipping", url)
	respByte, err := Post(url, header, req.Body, d.Logger)
	if err != nil {
		return message.ErrUpdateOrderShipping
	}

	_ = level.Info(logger).Log("d.UpdateOrderShipping", url)

	response := map[string]interface{}{}
	err = json.Unmarshal(respByte, &response)

	if err != nil {
		return message.ErrUpdateOrderShipping
	}

	return message.SuccessMsg
}
