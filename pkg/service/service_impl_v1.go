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

// swagger:route POST /prescription/product/ prescription product
// Create product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
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

// swagger:route GET /prescription/products/ prescription productList
// Get product lists
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *serviceImplV1) Show(ctx context.Context) (entity.ResponseHttp, error) {
	products, err := s.repostory.FindAll(ctx)
	if err != nil {
		s.logger.Log(err)
		return entity.ResponseHttp{
			Code:    common.ERRCODE_DB,
			Message: "Internal Server Error",
			Data:    "",
		}, nil
	}

	return entity.ResponseHttp{
		Code:    common.ERRCODE_SUCCESS,
		Message: "Display Products",
		Data:    products,
	}, nil
}

// swagger:route GET /prescription/product/{id} prescription getProduct
// Get product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *serviceImplV1) GetProduct(ctx context.Context, uid string) (entity.ResponseHttp, error) {
	products, err := s.repostory.FindById(ctx, uid)
	if err != nil {
		s.logger.Log(err)
		return entity.ResponseHttp{
			Code:    common.ERRCODE_DB,
			Message: "Internal Server Error",
			Data:    "",
		}, nil
	}

	return entity.ResponseHttp{
		Code:    common.ERRCODE_SUCCESS,
		Message: "Display Products",
		Data:    products,
	}, nil
}

// swagger:route PUT /prescription/product/{id} prescription update_product
// Get product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *serviceImplV1) UpdateProduct(ctx context.Context, req entity.JSONRequestUpdateProduct) (entity.ResponseHttp, error) {
	_, err := s.repostory.FindById(ctx, req.Id)
	if err != nil {
		s.logger.Log(err)
		return entity.ResponseHttp{
			Code:    common.ERRCODE_VALIDATE,
			Message: "Product Not Found",
			Data:    "",
		}, nil
	}

	product := entity.Product{
		Name:   req.Name,
		Sku:    req.Sku,
		Uom:    req.Uom,
		Weight: req.Weight,
	}

	update, err := s.repostory.Update(ctx, req.Id, &product)

	result := entity.JSONResponseProduct{
		Id:     update.UID,
		Name:   update.Name,
		Sku:    update.Sku,
		Uom:    update.Uom,
		Weight: update.Weight,
	}

	return entity.ResponseHttp{
		Code:    common.ERRCODE_SUCCESS,
		Message: "Display Products",
		Data:    result,
	}, nil
}
