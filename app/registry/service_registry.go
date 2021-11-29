package registry

import (
	rp "go-klikdokter/app/repository"
	"go-klikdokter/app/service"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

func RegisterProductService(db *gorm.DB, logger log.Logger) service.ProductService {
	return service.NewProductService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewProductRepository(rp.NewBaseRepository(db)),
	)
}
