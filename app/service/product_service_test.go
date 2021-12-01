package service

import (
	"encoding/json"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository"
	"go-klikdokter/helper/message"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger log.Logger

//var db *gorm.DB
var err error

var baseRepository = &repository.BaseRepositoryMock{Mock: mock.Mock{}}
var productRepository = &repository.ProductRepositoryMock{Mock: mock.Mock{}}
var service = NewProductService(logger, baseRepository, productRepository)

func init() {
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
	//db.AutoMigrate(&entity.Product{})
}

func TestCreateProduct(t *testing.T) {
	req := request.SaveProductRequest{
		Name:   "Prenagen",
		Weight: 100,
		Sku:    "SKU_X",
		Uom:    "Pcs",
	}

	result, _, _ := service.CreateProduct(req)

	assert.NotNil(t, result)
	assert.Equal(t, "Prenagen", result.Name, "Name must be Prenagen")
	assert.Equal(t, int32(100), result.Weight, "Weight must be 100")
	assert.Equal(t, "SKU_X", result.Sku, "SKU must be SKU_X")
	assert.Equal(t, "Pcs", result.Uom, "UOM must be Pcs")
}

func TestGetProduct(t *testing.T) {
	product := entity.Product{
		Uom:    "Pcs",
		Sku:    "SKU_x",
		Name:   "Prenagen",
		Weight: 110,
	}

	uid := "123"
	productRepository.Mock.On("FindByUid", &uid).Return(product)
	result, _, _ := service.GetProduct(uid)

	type responseHttp struct {
		Meta       interface{}    `json:"meta"`
		Pagination *interface{}   `json:"pagination,omitempty"`
		Data       entity.Product `json:"data"`
	}

	assert.NotNil(t, result, "Cannot nil")
	assert.Equal(t, "Prenagen", result.Name, "Name must be Prenagen")
	assert.Equal(t, int32(110), result.Weight, "Weight must be 110")
	assert.Equal(t, "SKU_x", result.Sku, "SKU must be SKU_x")
	assert.Equal(t, "Pcs", result.Uom, "UOM must be Pcs")
}

func TestDeleteProduct(t *testing.T) {
	product := entity.Product{
		Uom:    "Pcs",
		Sku:    "SKU_x",
		Name:   "Prenagen",
		Weight: 110,
	}

	uid := "123"
	productRepository.Mock.On("FindByUid", &uid).Return(product)
	code, msg := service.DeleteProduct(uid)

	assert.Equal(t, message.CODE_SUCCESS, code, "Code must be 1000")
	assert.Equal(t, message.MSG_SUCCESS, msg, "Message must be Success")
}

func TestUpdateProduct(t *testing.T) {
	req := request.SaveProductRequest{
		Uom:    "Pcs",
		Sku:    "SKU_x",
		Name:   "Prenagen",
		Weight: 110,
	}

	product := entity.Product{
		Uom:    "Pcs",
		Sku:    "SKU_x",
		Name:   "Prenagen",
		Weight: 110,
	}

	uid := "123"
	productRepository.Mock.On("FindByUid", &uid).Return(product)
	result := service.UpdateProduct(&uid, &req)

	type responseHttp struct {
		Meta struct {
			Code      int    `json:"code"`
			Message   string `json:"message"`
			RequestId string `json:"request_id"`
		} `json:"meta"`
		Pagination *interface{}   `json:"pagination,omitempty"`
		Data       entity.Product `json:"data"`
	}

	var data responseHttp
	jsonData, _ := json.Marshal(result)
	_ = json.Unmarshal(jsonData, &data)
	assert.Equal(t, 1000, data.Meta.Code, "Code must be 1000")
	assert.Equal(t, "Success", data.Meta.Message, "Message must be Success")
}

func TestGetListProduct(t *testing.T) {
	req := request.ProductListRequest{}

	var productList = []entity.Product{
		{
			Uom:    "Pcs",
			Sku:    "SKU_x",
			Name:   "Prenagen",
			Weight: 110,
		},
		{
			Uom:    "Box",
			Sku:    "SKU_AZX",
			Name:   "HydroCoco",
			Weight: 210,
		},
	}

	resPaginate := response.PaginationResponse{
		Limit:        1,
		Page:         1,
		TotalPage:    10,
		TotalRecords: 100,
	}

	filter := map[string]interface{}{
		"name": "",
		"sku":  "",
		"uom":  "",
	}
	productRepository.Mock.On("FindByParams", 10, 1, "", filter).Return(productList, resPaginate, nil)
	result := service.GetList(req)

	type responseHttp struct {
		Meta struct {
			Code      int    `json:"code"`
			Message   string `json:"message"`
			RequestId string `json:"request_id"`
		} `json:"meta"`
		Pagination *interface{}     `json:"pagination,omitempty"`
		Data       []entity.Product `json:"data"`
	}

	var data responseHttp
	jsonData, _ := json.Marshal(result)
	_ = json.Unmarshal(jsonData, &data)
	assert.Equal(t, 1000, data.Meta.Code, "Code must be 1000")
	assert.Equal(t, "Success", data.Meta.Message, "Message must be Success")
}
