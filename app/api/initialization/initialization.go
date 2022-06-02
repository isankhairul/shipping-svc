package initialization

import (
	"go-klikdokter/app/api/transport"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/database"
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

	// db.Migrator().DropTable(&entity.ShippmentPredefined{})
	//Define auto migration here
	_ = db.AutoMigrate(&entity.Courier{})
	_ = db.AutoMigrate(&entity.ShippmentPredefined{})
	// db.Migrator().DropTable(&entity.CourierCoverageCode{})

	//Define auto migration here
	// db.Migrator().DropTable(&entity.CourierCoverageCode{})

	//Define auto migration here
	_ = db.AutoMigrate(&entity.CourierCoverageCode{})
	_ = db.AutoMigrate(&entity.CourierService{})
	_ = db.AutoMigrate(&entity.Channel{})

	seedingPredefined(db, logger)

	return db, nil
}

func seedingPredefined(db *gorm.DB, logger log.Logger) {
	svc := registry.RegisterShipmentPredefinedService(db, logger)
	req := request.CreateShipmentPredefinedRequest{Type: "courier_type", Code: "third_party", Title: "Third Party"}
	_, _ = svc.CreateShipmentPredefined(req)
	req1 := request.CreateShipmentPredefinedRequest{Type: "courier_type", Code: "merchant", Title: "Merchant Courier"}
	_, _ = svc.CreateShipmentPredefined(req1)
	req2 := request.CreateShipmentPredefinedRequest{Type: "courier_type", Code: "internal", Title: "Internal Courier"}
	_, _ = svc.CreateShipmentPredefined(req2)
	req3 := request.CreateShipmentPredefinedRequest{Type: "shipping_type", Code: "instant", Title: "Instant", Note: "Waktu Pengiriman 3 Jam"}
	_, _ = svc.CreateShipmentPredefined(req3)
	req4 := request.CreateShipmentPredefinedRequest{Type: "shipping_type", Code: "same_day", Title: "Same Day", Note: "Waktu Pengiriman 6-8 Jam"}
	_, _ = svc.CreateShipmentPredefined(req4)
	req5 := request.CreateShipmentPredefinedRequest{Type: "shipping_type", Code: "regular", Title: "Regular", Note: "Waktu Pengiriman (2-4 hari)"}
	_, _ = svc.CreateShipmentPredefined(req5)
	req6 := request.CreateShipmentPredefinedRequest{Type: "shipping_type", Code: "next_day", Title: "Next Day", Note: "Waktu Pengiriman (1 hari)"}
	_, _ = svc.CreateShipmentPredefined(req6)
	req7 := request.CreateShipmentPredefinedRequest{Type: "courier_code", Code: "shipper", Title: "Shipper", Note: ""}
	_, _ = svc.CreateShipmentPredefined(req7)
	req8 := request.CreateShipmentPredefinedRequest{Type: "courier_code", Code: "gojek", Title: "Gojek", Note: ""}
	_, _ = svc.CreateShipmentPredefined(req8)
}

func InitRouting(db *gorm.DB, logger log.Logger) *http.ServeMux {
	// Service registry
	courierSvc := registry.RegisterCourierService(db, logger)
	courierCoverageCodeSvc := registry.RegisterCourierCoverageCodeService(db, logger)
	channelSvc := registry.RegisterChannelService(db, logger)
	shipmentPredefinedService := registry.RegisterShipmentPredefinedService(db, logger)

	// Transport initialization
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP")) //don't delete or change this !!
	courierHttp := transport.CourierHttpHandler(courierSvc, log.With(logger, "CourierTransportLayer", "HTTP"))
	courierCoverageCodeHttp := transport.CourierCoverageCodeHttpHandler(courierCoverageCodeSvc, log.With(logger, "CourierCoverageCodeTransportLayer", "HTTP"))
	channelHttp := transport.ChannelHttpHandler(channelSvc, log.With(logger, "ChannelTransportLayer", "HTTP"))
	shipmentPredefinedHttp := transport.ShipmentPredefinedHandler(shipmentPredefinedService, log.With(logger, "ShipmentPredefinedTransportLayer", "HTTP"))

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //don't delete or change this!!
	mux.Handle("/courier/", courierHttp)
	mux.Handle("/other/", shipmentPredefinedHttp)
	mux.Handle("/courier/courier-coverage-code/", courierCoverageCodeHttp)
	mux.Handle("/channel/", channelHttp)

	return mux
}
