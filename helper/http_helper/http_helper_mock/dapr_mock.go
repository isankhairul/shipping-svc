package http_helper_mock

import (
	"go-klikdokter/app/model/request"
	"go-klikdokter/helper/message"

	"github.com/stretchr/testify/mock"
)

type DaprEndpointMock struct {
	Mock mock.Mock
}

func (d *DaprEndpointMock) UpdateOrderShipping(req *request.UpdateOrderShipping) message.Message {
	return message.SuccessMsg
}
