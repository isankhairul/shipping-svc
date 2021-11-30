package service

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type ProductService interface {
	CreateProduct(input request.SaveProductRequest) (*entity.Product, int, string)
	GetProduct(uid string) (*entity.Product, int, string)
	GetList(input request.ProductListRequest) ([]entity.Product, *base.Pagination, int, string)
	UpdateProduct(uid string, input request.SaveProductRequest) (int, string)
	DeleteProduct(uid string) (int, string)
}

type productServiceImpl struct {
	logger      log.Logger
	baseRepo    repository.BaseRepository
	productRepo repository.ProductRepository
}

func NewProductService(
	lg log.Logger,
	br repository.BaseRepository,
	pr repository.ProductRepository,
) ProductService {
	return &productServiceImpl{lg, br, pr}
}

// swagger:route POST /product/  product
// Create product

// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  201: Created
func (s *productServiceImpl) CreateProduct(input request.SaveProductRequest) (*entity.Product, int, string) {
	logger := log.With(s.logger, "ProductService", "CreateProduct")
	s.baseRepo.BeginTx()
	//Set request to entity
	product := entity.Product{
		Name:   input.Name,
		Sku:    input.Sku,
		Uom:    input.Uom,
		Weight: input.Weight,
	}

	result, err := s.productRepo.Create(&product)
	if err != nil {
		level.Error(logger).Log(err)
		s.baseRepo.RollbackTx()
		return nil, message.CODE_ERR_DB, message.MSG_ERR_SAVE_DATA
	}
	s.baseRepo.CommitTx()

	return result, message.CODE_SUCCESS, ""
}

// swagger:route GET /product/  get one product
// Get product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) GetProduct(uid string) (*entity.Product, int, string) {
	logger := log.With(s.logger, "ProductService", "GetProduct")

	result, err := s.productRepo.FindByUid(&uid)
	if err != nil {
		level.Error(logger).Log(err)
		return nil, message.CODE_ERR_DB, message.MSG_ERR_DB
	}

	if result == nil {
		return nil, message.CODE_ERR_DB, message.MSG_NO_DATA
	}

	return result, message.CODE_SUCCESS, ""
}

// swagger:route GET /product/list  productList
// Get products
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) GetList(input request.ProductListRequest) ([]entity.Product, *base.Pagination, int, string) {
	logger := log.With(s.logger, "ProductService", "GetList")

	//Set default value
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}

	filter := map[string]interface{}{
		"name": input.Name,
		"sku":  input.Sku,
	}

	result, pagination, err := s.productRepo.FindByParams(input.Limit, input.Page, input.Sort, filter)
	if err != nil {
		level.Error(logger).Log(err)
		return nil, nil, message.CODE_ERR_DB, message.MSG_ERR_DB
	}

	if result == nil {
		level.Warn(logger).Log(message.MSG_NO_DATA)
		return nil, nil, message.CODE_ERR_DB, message.MSG_NO_DATA
	}

	return result, pagination, message.CODE_SUCCESS, ""
}

// swagger:route PUT /prescription/product/{id} prescription update_product
// Get product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) UpdateProduct(uid string, input request.SaveProductRequest) (int, string) {
	logger := log.With(s.logger, "ProductService", "UpdateProduct")

	_, err := s.productRepo.FindByUid(&uid)
	if err != nil {
		level.Error(logger).Log(err)
		return message.CODE_ERR_DB, message.MSG_INVALID_REQUEST
	}

	data := map[string]interface{}{
		"name":   input.Name,
		"sku":    input.Sku,
		"uom":    input.Uom,
		"weight": input.Weight,
	}

	err = s.productRepo.Update(uid, data)
	if err != nil {
		level.Error(logger).Log(err)
		return message.CODE_ERR_DB, message.MSG_ERR_SAVE_DATA
	}

	return message.CODE_SUCCESS, ""
}

// swagger:route DELETE /product/
// Delete product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) DeleteProduct(uid string) (int, string) {
	logger := log.With(s.logger, "ProductService", "DeleteProduct")

	_, err := s.productRepo.FindByUid(&uid)
	if err != nil {
		level.Error(logger).Log(err)
		return message.CODE_ERR_DB, message.MSG_INVALID_REQUEST
	}

	err = s.productRepo.Delete(uid)
	if err != nil {
		level.Error(logger).Log(err)
		return message.CODE_ERR_DB, message.MSG_ERR_SAVE_DATA
	}

	return message.CODE_SUCCESS, ""
}
