package service

import (
	"context"

	"gokit_example/pkg/entity"
)

type Service interface {
	Create(ctx context.Context, req entity.JSONRequestProduct) (entity.ResponseHttp, error)
}
