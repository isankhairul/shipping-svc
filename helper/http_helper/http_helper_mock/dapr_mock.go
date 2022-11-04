package http_helper_mock

import (
	"github.com/stretchr/testify/mock"
)

type DaprEndpointMock struct {
	Mock mock.Mock
}

func (d *DaprEndpointMock) PublishKafka(topicName string, req interface{}) {
	//implemented
}
