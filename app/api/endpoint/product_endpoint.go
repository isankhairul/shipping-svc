package endpoint

import (
	"context"
	"fmt"
	"gokit_example/app/model/request"
	"gokit_example/app/service"

	"github.com/go-kit/kit/endpoint"
)

type ProductEndpoint struct {
	Save   endpoint.Endpoint
	Show   endpoint.Endpoint
	List   endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

func MakeProductEndpoints(s service.ProductService) ProductEndpoint {
	return ProductEndpoint{
		Save:   makeSaveProduct(s),
		Show:   makeShowProducts(s),
		List:   makeGetProducts(s),
		Update: makeUpdateProduct(s),
		Delete: makeDeleteProduct(s),
	}
}

func makeSaveProduct(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (response interface{}, err error) {
		req := rqst.(request.SaveProductRequest)
		return s.CreateProduct(req), nil
	}
}

func makeShowProducts(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (response interface{}, err error) {
		return s.GetProduct(fmt.Sprint(rqst)), nil
	}
}

func makeGetProducts(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (response interface{}, err error) {
		req := rqst.(request.ProductListRequest)
		return s.GetList(req), nil
	}
}

func makeUpdateProduct(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (response interface{}, err error) {
		req := rqst.(request.SaveProductRequest)
		return s.UpdateProduct(req.Uid, req), nil
	}
}

func makeDeleteProduct(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (response interface{}, err error) {
		req := rqst.(request.SaveProductRequest)
		return s.UpdateProduct(req.Uid, req), nil
	}
}
