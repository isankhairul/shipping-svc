package endpoint

import (
	"context"

	"gokit_example/pkg/entity"
	"gokit_example/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

type Endpoint struct {
	Check endpoint.Endpoint
}

func MakeServerEndpoints(s service.Service) Endpoint {
	return Endpoint{
		Check: makeCheckDevice(s),
	}
}

func makeCheckDevice(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.JSONRequestProduct)
		return s.Create(ctx, req)
	}
}
