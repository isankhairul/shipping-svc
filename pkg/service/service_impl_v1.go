package service

import (
	"context"
	"errors"
	"gokit_example/pkg/common"
	"gokit_example/pkg/entity"
	"gokit_example/pkg/repository"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
)

var (
	ErrAlreadyExists = errors.New(common.ERRMSG_ALREADYEXISTS)
	ErrNotFound      = errors.New(common.ERRMSG_NOTFOUND)
)

type serviceImplV1 struct {
	logger    log.Logger
	repostory repository.ProductRepository
	redis     *redis.Client
}

func NewServiceImplV1(logger log.Logger, repostory repository.ProductRepository, redis *redis.Client) Service {
	return &serviceImplV1{
		repostory: repostory,
		logger:    logger,
		redis:     redis,
	}
}

// Create
func (s *serviceImplV1) Create(ctx context.Context, req entity.JSONRequestProduct) (entity.ResponseHttp, error) {

	product := entity.Product{
		Name:   req.Name,
		Sku:    req.Sku,
		Uom:    req.Uom,
		Weight: req.Weight,
	}
	productUid, err := s.repostory.Save(ctx, &product)
	if err != nil {
		s.logger.Log(err)
		return entity.ResponseHttp{
			Code:    common.ERRCODE_DB,
			Message: "Internal Server Error",
			Data:    "",
		}, nil
	}

	result, err := s.repostory.FindById(ctx, productUid)
	if err != nil {
		s.logger.Log(err)
		return entity.ResponseHttp{
			Code:    common.ERRCODE_DB,
			Message: "Internal Server Error",
			Data:    "",
		}, nil
	}

	s.logger.Log(err)
	return entity.ResponseHttp{
		Code:    common.ERRCODE_SUCCESS,
		Message: "Success Created",
		Data:    result,
	}, nil

}
