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
		rp.NewCourierServiceRepository(rp.NewBaseRepository(db)),
		rp.NewShipmentPredefinedRepository(rp.NewBaseRepository(db)),
	)
}

func RegisterChannelService(db *gorm.DB, logger log.Logger) service.ChannelService {
	return service.NewChannelService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewChannelRepository(rp.NewBaseRepository(db)),
		rp.NewShippingCourierStatusRepository(rp.NewBaseRepository(db)),
	)
}

func RegisterChannelCourierService(db *gorm.DB, logger log.Logger) service.ChannelCourierService {
	repo := rp.NewBaseRepository(db)
	return service.NewChannelCourierService(
		logger, repo,
		rp.NewChannelCourierRepository(repo),
		rp.NewChannelCourierServiceRepository(repo),
		rp.NewCourierServiceRepository(repo))
}

func RegisterCourierCoverageCodeService(db *gorm.DB, logger log.Logger) service.CourierCoverageCodeService {
	return service.NewCourierCoverageCodeService(
		logger,
		rp.NewBaseRepository(db),
		rp.NewCourierCoverageCodeRepository(rp.NewBaseRepository(db)))
}

func RegisterShipmentPredefinedService(db *gorm.DB, logger log.Logger) service.ShipmentPredefinedService {
	repo := rp.NewBaseRepository(db)

	return service.NewShipmentPredefinedService(
		logger, repo, rp.NewShipmentPredefinedRepository(repo),
	)
}

func RegisterChannelCourierServiceService(db *gorm.DB, logger log.Logger) service.ChannelCourierServiceService {
	repo := rp.NewBaseRepository(db)
	return service.NewChannelCourierServiceService(
		logger, repo,
		rp.NewChannelCourierRepository(repo),
		rp.NewChannelCourierServiceRepository(repo),
		rp.NewCourierServiceRepository(repo))
}
