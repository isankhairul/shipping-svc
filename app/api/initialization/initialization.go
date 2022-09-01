package initialization

import (
	"fmt"
	"go-klikdokter/app/api/transport"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/database"
	"go-klikdokter/helper/global"
	"go-klikdokter/pkg/cache"
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

	return db, nil
}

func InitRouting(db *gorm.DB, logger log.Logger, redis cache.RedisCache) *http.ServeMux {
	// Service registry
	courierSvc := registry.RegisterCourierService(db, logger)
	channelCourierSvc := registry.RegisterChannelCourierService(db, logger)
	channelSvc := registry.RegisterChannelService(db, logger)
	shipmentPredefinedService := registry.RegisterShipmentPredefinedService(db, logger)
	courierCoverageCodeSvc := registry.RegisterCourierCoverageCodeService(db, logger)
	channelCourierServiceSvc := registry.RegisterChannelCourierServiceService(db, logger)
	shippingService := registry.RegisterShippingService(db, logger, redis)

	// Transport initialization
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP")) //don't delete or change this !!
	courierHttp := transport.CourierHttpHandler(courierSvc, channelCourierSvc, log.With(logger, "CourierTransportLayer", "HTTP"))
	channelCourierHttp := transport.ChannelCourierHttpHandler(channelCourierSvc, log.With(logger, "ChannelCourierTransportLayer", "HTTP"))
	courierCoverageCodeHttp := transport.CourierCoverageCodeHttpHandler(courierCoverageCodeSvc, log.With(logger, "CourierCoverageCodeTransportLayer", "HTTP"))
	shipmentPredefinedHttp := transport.ShipmentPredefinedHandler(shipmentPredefinedService, log.With(logger, "ShipmentPredefinedTransportLayer", "HTTP"))
	channelHttp := transport.ChannelHttpHandler(channelSvc, channelCourierSvc, log.With(logger, "ChannelTransportLayer", "HTTP"))
	channelCourierServiceHttp := transport.ChannelCourierServiceHttpHandler(channelCourierServiceSvc, log.With(logger, "ChannelCourierServiceTransportLayer", "HTTP"))
	shippingHttp := transport.ShippingHttpHandler(shippingService, log.With(logger, "ShippingTransportLayer", "HTTP"))

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //don't delete or change this!!
	mux.Handle(fmt.Sprint(global.PrefixBase, global.PrefixCourier), courierHttp)
	mux.Handle(fmt.Sprint(global.PrefixBase, global.PrefixOther), shipmentPredefinedHttp)
	mux.Handle(fmt.Sprint(global.PrefixBase, global.PrefixCourierCoverageCode), courierCoverageCodeHttp)
	mux.Handle(fmt.Sprint(global.PrefixBase, global.PrefixChannel), channelHttp)
	mux.Handle(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourier), channelCourierHttp)
	mux.Handle(fmt.Sprint(global.PrefixBase, global.PrefixChannelCourierService), channelCourierServiceHttp)
	mux.Handle(fmt.Sprint(global.PrefixBase, global.PrefixShipping), shippingHttp)

	return mux
}

func InitCache(logger log.Logger) (cache.RedisCache, error) {

	redis, err := cache.SetupRedisConnection(
		viper.GetString("cache.redis.host"),
		viper.GetString("cache.redis.port"),
		viper.GetInt("cache.redis.index.primary"),
		viper.GetString("cache.redis.password"),
		viper.GetBool("cache.redis.is-active"),
		viper.GetInt("cache.redis.expired-in-minute.default"),
	)

	return redis, err
}
