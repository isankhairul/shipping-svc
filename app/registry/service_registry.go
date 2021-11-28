package registry

import (
	rp "gokit_example/app/repository"
	"gokit_example/app/service"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

func NewProductService(db *gorm.DB, logger log.Logger) service.ProductService {
	return service.NewproductServiceImpl(
		logger,
		rp.NewBaseRepository(db),
		rp.NewProductRepository(rp.NewBaseRepository(db)),
	)
}
