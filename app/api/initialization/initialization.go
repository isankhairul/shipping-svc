package initialization

import (
	"go-klikdokter/app/api/transport"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/database"
	"net/http"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func DbInit() (*gorm.DB, error) {
	// Init DB Connection
	db, err := database.NewConnectionDB(viper.GetString("database.driver"), viper.GetString("database.dbname"),
		viper.GetString("database.host"), viper.GetString("database.username"), viper.GetString("database.password"),
		viper.GetInt("database.port"))
	if err != nil {
		return nil, err
	}

	//Define auto migration here
	db.AutoMigrate(&entity.Product{})
	db.AutoMigrate(&entity.Doctor{})

	return db, nil
}

func InitRouting(db *gorm.DB, logger log.Logger) *http.ServeMux {
	// Service registry
	prodSvc := registry.RegisterProductService(db, logger)
	doctorSvc := registry.RegisterDoctorService(db, logger)

	// Transport initialization
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP")) //don't delete or change this !!
	prodHttp := transport.ProductHttpHandler(prodSvc, log.With(logger, "ProductTransportLayer", "HTTP"))
	doctorHttp := transport.DoctorHttpHandler(doctorSvc, log.With(logger, "ProductTransportLayer", "HTTP"))

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //don't delete or change this!!
	mux.Handle("/product/", prodHttp)
	mux.Handle("/doctor/", doctorHttp)

	return mux
}
