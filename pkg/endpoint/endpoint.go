package endpoint

import (
	"context"
	"fmt"
	"gokit_example/pkg/entity"
	"gokit_example/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

type Endpoint struct {
	Save          endpoint.Endpoint
	Show          endpoint.Endpoint
	GetProduct    endpoint.Endpoint
	UpdateProduct endpoint.Endpoint
}

func MakeServerEndpoints(s service.Service) Endpoint {
	return Endpoint{
		Save:          makeSaveProduct(s),
		Show:          makeShowProducts(s),
		GetProduct:    makeGetProduct(s),
		UpdateProduct: makeUpdateProduct(s),
	}
}

func makeSaveProduct(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.JSONRequestProduct)
		return s.Create(ctx, req)
	}
}

func makeShowProducts(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.Show(ctx)
	}
}

func makeGetProduct(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetProduct(ctx, fmt.Sprint(request))
	}
}

func makeUpdateProduct(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.JSONRequestUpdateProduct)
		return s.UpdateProduct(ctx, req)
	}
}
