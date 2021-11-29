package endpoint

import (
	"context"
	"errors"
	"fmt"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/service"

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
		Show:   makeShowProduct(s),
		List:   makeGetProducts(s),
		Update: makeUpdateProduct(s),
		Delete: makeDeleteProduct(s),
	}
}

func makeSaveProduct(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveProductRequest)
		result, code, msg := s.CreateProduct(req)
		if msg != "" {
			return base.SetHttpResponse(code, msg, nil, nil), errors.New(fmt.Sprintf("%v", code))
		}

		return base.SetHttpResponse(code, msg, result, nil), nil
	}
}

func makeShowProduct(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		result, code, msg := s.GetProduct(fmt.Sprint(rqst))
		if msg != "" {
			return base.SetHttpResponse(code, msg, nil, nil), errors.New(fmt.Sprintf("%v", code))
		}

		return base.SetHttpResponse(code, msg, result, nil), nil
	}
}

func makeGetProducts(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.ProductListRequest)
		result, pagination, code, msg := s.GetList(req)
		if msg != "" {
			return base.SetHttpResponse(code, msg, nil, nil), errors.New(fmt.Sprintf("%v", code))
		}

		return base.SetHttpResponse(code, msg, result, pagination), nil
	}
}

func makeUpdateProduct(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		req := rqst.(request.SaveProductRequest)
		code, msg := s.UpdateProduct(req.Uid, req)
		if msg != "" {
			return base.SetHttpResponse(code, msg, nil, nil), errors.New(fmt.Sprintf("%v", code))
		}

		return base.SetHttpResponse(code, msg, nil, nil), nil
	}
}

func makeDeleteProduct(s service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
		code, msg := s.DeleteProduct(fmt.Sprint(rqst))
		if msg != "" {
			return base.SetHttpResponse(code, msg, nil, nil), errors.New(fmt.Sprintf("%v", code))
		}

		return base.SetHttpResponse(code, msg, nil, nil), nil
	}
}
