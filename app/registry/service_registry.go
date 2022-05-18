package registry

import (
	rp "go-klikdokter/app/repository"
	"go-klikdokter/app/service"

	"github.com/go-kit/log"
	"gorm.io/gorm"
)

func RegisterCourierService(db *gorm.DB, logger log.Logger) service.CourierService {
	return service.NewCourierService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewCourierRepository(rp.NewBaseRepository(db)),
	)
}

func RegisterChannelService(db *gorm.DB, logger log.Logger) service.ChannelService {
	return service.NewChannelService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewChannelRepository(rp.NewBaseRepository(db)),
	)
}
