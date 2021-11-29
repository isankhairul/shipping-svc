package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gokit_example/app/model/entity"
	"gokit_example/app/model/request"
	"gokit_example/app/repository"
	"gokit_example/helper/database"
	"gorm.io/gorm"
	"os"
	"testing"
)

var logger log.Logger
var db *gorm.DB
var err error

var baseRepository = &repository.BaseRepositoryMock{Mock: mock.Mock{}}
var productRepository = &repository.ProductRepositoryMock{Mock: mock.Mock{}}
var service = NewproductServiceImpl(logger, baseRepository, productRepository)

func init() {
	db, err = database.NewConnectionDB("sqlite", "", "", "", "", 5432)
	if err != nil {
		fmt.Println(err)
	}
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

}

func TestCreateProduct(t *testing.T) {
	req := request.SaveProductRequest{
		Name:   "Prenagen",
		Weight: 100,
		Sku:    "SKU_X",
		Uom:    "Pcs",
	}

	result := service.CreateProduct(req)

	type responseHttp struct {
		Meta       interface{}    `json:"meta"`
		Pagination *interface{}   `json:"pagination,omitempty"`
		Data       entity.Product `json:"data"`
	}

	var data responseHttp
	jsonData, _ := json.Marshal(result)
	_ = json.Unmarshal(jsonData, &data)
	assert.Equal(t, "Prenagen", data.Data.Name, "Name must be Prenagen")
	assert.Equal(t, int32(100), data.Data.Weight, "Weight must be 100")
	assert.Equal(t, "SKU_X", data.Data.Sku, "SKU must be SKU_X")
	assert.Equal(t, "Pcs", data.Data.Uom, "UOM must be Pcs")
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
	result := service.GetProduct(uid)

	type responseHttp struct {
		Meta       interface{}    `json:"meta"`
		Pagination *interface{}   `json:"pagination,omitempty"`
		Data       entity.Product `json:"data"`
	}

	var data responseHttp
	jsonData, _ := json.Marshal(result)
	_ = json.Unmarshal(jsonData, &data)
	assert.Equal(t, "Prenagen", data.Data.Name, "Name must be Prenagen")
	assert.Equal(t, int32(110), data.Data.Weight, "Weight must be 110")
	assert.Equal(t, "SKU_x", data.Data.Sku, "SKU must be SKU_x")
	assert.Equal(t, "Pcs", data.Data.Uom, "UOM must be Pcs")
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
	result := service.DeleteProduct(uid)

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
