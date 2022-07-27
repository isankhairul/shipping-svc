package initialization

import (
	"go-klikdokter/app/api/transport"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/database"
	"go-klikdokter/pkg/util"
	"net/http"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func DbInit(logger log.Logger) (*gorm.DB, error) {
	// Init DB Connection
	db, err := database.NewConnectionDB(config.GetConfigString(viper.GetString("database.driver")), config.GetConfigString(viper.GetString("database.dbname")),
		config.GetConfigString(viper.GetString("database.host")), config.GetConfigString(viper.GetString("database.username")), config.GetConfigString(viper.GetString("database.password")),
		config.GetConfigInt(viper.GetString("database.port")))
	if err != nil {
		return nil, err
	}

	//Define auto migration here
	_ = db.AutoMigrate(&entity.Courier{})
	_ = db.AutoMigrate(&entity.ShippmentPredefined{})
	_ = db.AutoMigrate(&entity.CourierCoverageCode{})
	_ = db.AutoMigrate(&entity.CourierService{})
	_ = db.AutoMigrate(&entity.Channel{})
	_ = db.AutoMigrate(&entity.ChannelCourier{})
	_ = db.AutoMigrate(&entity.ChannelCourierService{})
	_ = db.AutoMigrate(&entity.ShippingStatus{})
	_ = db.AutoMigrate(&entity.ShippingCourierStatus{})

	if ok := db.Migrator().HasColumn(&entity.ChannelCourierService{}, "channel_id"); ok {
		_ = db.Migrator().DropColumn(&entity.ChannelCourierService{}, "channel_id")
	}

	if ok := db.Migrator().HasColumn(&entity.ChannelCourierService{}, "courier_id"); ok {
		_ = db.Migrator().DropColumn(&entity.ChannelCourierService{}, "courier_id")
	}

	if ok := db.Migrator().HasColumn(&entity.CourierCoverageCode{}, "courier_uid"); ok {
		_ = db.Migrator().DropColumn(&entity.CourierCoverageCode{}, "courier_uid")
	}

	return db, nil
}

func InitRouting(db *gorm.DB, logger log.Logger) *http.ServeMux {
	// Service registry
	courierSvc := registry.RegisterCourierService(db, logger)
	channelCourierSvc := registry.RegisterChannelCourierService(db, logger)
	channelSvc := registry.RegisterChannelService(db, logger)
	shipmentPredefinedService := registry.RegisterShipmentPredefinedService(db, logger)
	courierCoverageCodeSvc := registry.RegisterCourierCoverageCodeService(db, logger)
	channelCourierServiceSvc := registry.RegisterChannelCourierServiceService(db, logger)

	// Transport initialization
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP")) //don't delete or change this !!
	courierHttp := transport.CourierHttpHandler(courierSvc, channelCourierSvc, log.With(logger, "CourierTransportLayer", "HTTP"))
	channelCourierHttp := transport.ChannelCourierHttpHandler(channelCourierSvc, log.With(logger, "ChannelCourierTransportLayer", "HTTP"))
	courierCoverageCodeHttp := transport.CourierCoverageCodeHttpHandler(courierCoverageCodeSvc, log.With(logger, "CourierCoverageCodeTransportLayer", "HTTP"))
	shipmentPredefinedHttp := transport.ShipmentPredefinedHandler(shipmentPredefinedService, log.With(logger, "ShipmentPredefinedTransportLayer", "HTTP"))
	channelHttp := transport.ChannelHttpHandler(channelSvc, log.With(logger, "ChannelTransportLayer", "HTTP"))
	channelCourierServiceHttp := transport.ChannelCourierServiceHttpHandler(channelCourierServiceSvc, log.With(logger, "ChannelCourierServiceTransportLayer", "HTTP"))

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //don't delete or change this!!
	mux.Handle(util.PrefixBase+"/courier/", courierHttp)
	mux.Handle(util.PrefixBase+"/other/", shipmentPredefinedHttp)
	mux.Handle(util.PrefixBase+"/courier/courier-coverage-code/", courierCoverageCodeHttp)
	mux.Handle(util.PrefixBase+"/channel/", channelHttp)
	mux.Handle(util.PrefixBase+"/channel/channel-courier/", channelCourierHttp)
	mux.Handle(util.PrefixBase+"/channel/channel-courier-service/", channelCourierServiceHttp)

	return mux
}
