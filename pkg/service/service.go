package service

import (
	"context"

	"gokit_example/pkg/entity"
)

type Service interface {
	Create(ctx context.Context, req entity.JSONRequestProduct) (entity.ResponseHttp, error)
	Show(ctx context.Context) (entity.ResponseHttp, error)
	GetProduct(ctx context.Context, uid string) (entity.ResponseHttp, error)
	UpdateProduct(ctx context.Context, req entity.JSONRequestUpdateProduct) (entity.ResponseHttp, error)
}
