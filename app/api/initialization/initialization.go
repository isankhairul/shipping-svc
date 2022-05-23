package initialization

import (
	"go-klikdokter/app/api/transport"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/database"
	"net/http"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func DbInit() (*gorm.DB, error) {
	// Init DB Connection
	db, err := database.NewConnectionDB(config.GetConfigString(viper.GetString("database.driver")), config.GetConfigString(viper.GetString("database.dbname")),
		config.GetConfigString(viper.GetString("database.host")), config.GetConfigString(viper.GetString("database.username")), config.GetConfigString(viper.GetString("database.password")),
		config.GetConfigInt(viper.GetString("database.port")))
	if err != nil {
		return nil, err
	}

	//Define auto migration here
	_ = db.AutoMigrate(&entity.Courier{})
	_ = db.AutoMigrate(&entity.CourierService{})
	_ = db.AutoMigrate(&entity.Channel{})

	// example Seeder
	// for i := 0; i < 1000; i++ {
	// 	fmt.Println("dijalankan")
	// 	fmt.Println("dijalankan")
	// 	product := entity.Product{}
	// 	err := faker.FakeData(&product)
	// 	db.Create(&product)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	return db, nil
}

func InitRouting(db *gorm.DB, logger log.Logger) *http.ServeMux {
	// Service registry
	courierSvc := registry.RegisterCourierService(db, logger)
	channelSvc := registry.RegisterChannelService(db, logger)

	// Transport initialization
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP")) //don't delete or change this !!
	courierHttp := transport.CourierHttpHandler(courierSvc, log.With(logger, "CourierTransportLayer", "HTTP"))
	channelHttp := transport.ChannelHttpHandler(channelSvc, log.With(logger, "ChannelTransportLayer", "HTTP"))

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //don't delete or change this!!
	mux.Handle("/courier/", courierHttp)
	mux.Handle("/channel/", channelHttp)

	return mux
}
