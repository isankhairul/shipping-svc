package registry

import (
	rp "go-klikdokter/app/repository"
	"go-klikdokter/app/service"

	"github.com/go-kit/log"
	"gorm.io/gorm"
)

func RegisterProductService(db *gorm.DB, logger log.Logger) service.ProductService {
	return service.NewProductService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewProductRepository(rp.NewBaseRepository(db)),
	)
}

func RegisterDoctorService(db *gorm.DB, logger log.Logger) service.DoctorService {
	return service.NewDoctorService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewDoctorRepository(rp.NewBaseRepository(db)),
	)
}

func RegisterCourierService(db *gorm.DB, logger log.Logger) service.CourierService {
	return service.NewCourierService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewCourierRepository(rp.NewBaseRepository(db)),
	)
}

func RegisterCourierServiceService(db *gorm.DB, logger log.Logger) service.CourierServiceService {
	return service.NewCourierServiceService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewCourierServiceRepository(rp.NewBaseRepository(db)),
	)
}
