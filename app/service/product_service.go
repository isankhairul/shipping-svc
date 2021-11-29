package service

import (
	"gokit_example/app/model/entity"
	"gokit_example/app/model/request"
	"gokit_example/app/model/response"
	"gokit_example/app/repository"
	"gokit_example/helper/message"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type ProductService interface {
	CreateProduct(input request.SaveProductRequest) interface{}
	GetList(input request.ProductListRequest) interface{}
	GetProduct(uid string) interface{}
	UpdateProduct(uid *string, input *request.SaveProductRequest) interface{}
	DeleteProduct(uid string) interface{}
}

type productServiceImpl struct {
	logger      log.Logger
	baseRepo    repository.BaseRepository
	productRepo repository.ProductRepository
}

func NewproductServiceImpl(
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
func (s *productServiceImpl) CreateProduct(input request.SaveProductRequest) interface{} {
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
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_ERR_SAVE_DATA, "", nil, nil)
	}
	s.baseRepo.CommitTx()

	return response.SetResponse(message.CODE_SUCCESS, message.MSG_SUCCESS, "", result, nil)
}

// swagger:route GET /product/  get one product
// Get product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) GetProduct(uid string) interface{} {
	logger := log.With(s.logger, "ProductService", "GetProduct")

	result, err := s.productRepo.FindByUid(&uid)
	if err != nil {
		level.Error(logger).Log(err)
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_ERR_DB, "", nil, nil)
	}

	if result == nil {
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_NO_DATA, "", nil, nil)
	}

	return response.SetResponse(message.CODE_SUCCESS, message.MSG_SUCCESS, "", result, nil)
}

// swagger:route GET /product/list  productList
// Get products
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) GetList(input request.ProductListRequest) interface{} {
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
		"uom":  input.Uom,
	}

	result, pagination, err := s.productRepo.FindByParams(input.Limit, input.Page, input.Sort, filter)
	if err != nil {
		level.Error(logger).Log(err)
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_ERR_DB, "", nil, nil)
	}

	if result == nil {
		level.Warn(logger).Log(message.MSG_NO_DATA)
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_NO_DATA, "", nil, nil)
	}

	return response.SetResponse(message.CODE_SUCCESS, message.MSG_SUCCESS, "", result, pagination)
}

// swagger:route PUT /prescription/product/{id} prescription update_product
// Get product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) UpdateProduct(uid *string, input *request.SaveProductRequest) interface{} {
	logger := log.With(s.logger, "ProductService", "UpdateProduct")

	_, err := s.productRepo.FindByUid(uid)
	if err != nil {
		level.Error(logger).Log(err)
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_INVALID_REQUEST, "", nil, nil)
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
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_ERR_SAVE_DATA, "", nil, nil)
	}

	return response.SetResponse(message.CODE_SUCCESS, message.MSG_SUCCESS, "", nil, nil)
}

// swagger:route DELETE /product/
// Delete product
//
// security:
// - apiKey: []
// responses:
//  401: ErrorResponse
//  200: SuccessResponse
func (s *productServiceImpl) DeleteProduct(uid string) interface{} {
	logger := log.With(s.logger, "ProductService", "DeleteProduct")

	_, err := s.productRepo.FindByUid(&uid)
	if err != nil {
		level.Error(logger).Log(err)
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_INVALID_REQUEST, "", nil, nil)
	}

	err = s.productRepo.Delete(uid)
	if err != nil {
		level.Error(logger).Log(err)
		return response.SetResponse(message.CODE_ERR_DB, message.MSG_ERR_SAVE_DATA, "", nil, nil)
	}

	return response.SetResponse(message.CODE_SUCCESS, message.MSG_SUCCESS, "", nil, nil)
}
