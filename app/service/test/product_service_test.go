package test

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/repository/repository_mock"
	"go-klikdokter/app/service"
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

var baseRepository = &repository_mock.BaseRepositoryMock{Mock: mock.Mock{}}
var productRepository = &repository_mock.ProductRepositoryMock{Mock: mock.Mock{}}
var svc = service.NewProductService(logger, baseRepository, productRepository)

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

	result, _, _ := svc.CreateProduct(req)

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
	result, _, _ := svc.GetProduct(uid)

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
	code, msg := svc.DeleteProduct(uid)

	assert.Equal(t, message.CODE_SUCCESS, code, "Code must be 1000")
	assert.Equal(t, "", msg, "Message must be Null")
}

func TestListProduct(t *testing.T) {
	req := request.ProductListRequest{
		Page: 1,
		Sort: "",
		UOM: "",
		Limit: 10,
	}

	product := []entity.Product{
		{
			Uom:    "Pcs",
			Sku:    "SKU_x",
			Name:   "Prenagen",
			Weight: 110,
		},
		{
			Uom:    "Box",
			Sku:    "SKU_34x",
			Name:   "Milna",
			Weight: 110,
		},
		{
			Uom:    "Pcs",
			Sku:    "SKU_4",
			Name:   "Hydro COCO",
			Weight: 199,
		},
	}

	filter := map[string]interface{}{
		"name": "",
		"sku":  "",
		"uom":  "",
	}

	paginationResult := base.Pagination{
		Records: 120,
		Limit: 10,
		Page: 1,
		TotalPage: 12,
	}


	productRepository.Mock.On("FindByParams",10,1,"",filter).Return(product, &paginationResult)
	
	products, pagination, code, msg := svc.GetList(req)

	assert.Equal(t, message.CODE_SUCCESS, code, "Code must be 1000")
	assert.Equal(t, "", msg, "Message must be null")
	assert.Equal(t, 3, len(products), "Count of products must be 3")
	assert.Equal(t, int64(120), pagination.Records, "Total record pagination must be 120")
	
}