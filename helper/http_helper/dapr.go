package http_helper

import (
	"strings"

	"github.com/go-kit/log/level"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
)

type DaprEndpoint interface {
	PublishKafka(topicName string, req interface{})
}

type dapr struct {
	Logger log.Logger
}

func NewDaprEndpoint(log log.Logger) DaprEndpoint {
	return &dapr{log}
}

func (d *dapr) PublishKafka(topicName string, req interface{}) {
	logger := log.With(d.Logger, "Webhook", "PublishKafka")
	url := viper.GetString("dapr.endpoint.publish-kafka")
	url = strings.ReplaceAll(url, "{topic-name}", topicName)
	header := map[string]string{"Content-Type": "application/json"}

	if _, err := Post(url, header, req, d.Logger); err != nil {
		_ = level.Error(logger).Log("PublishKafka", err.Error())
		return
	}
}
