package main

import (
	"fmt"
	"go-klikdokter/app/api/transport"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/consul"
	"go-klikdokter/helper/database"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//Load configuration
	viper.SetConfigType("yaml")
	var profile string = "dev"
	if os.Getenv("KD_ENV") != "" {
		profile = "prd"
	}

	var configFileName []string
	configFileName = append(configFileName, "config-")
	configFileName = append(configFileName, profile)
	viper.SetConfigName(strings.Join(configFileName, ""))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}

	// Logging init
	logfile, err := os.OpenFile(viper.GetString("server.output-file-path"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer logfile.Close()
	var logger log.Logger
	{
		if viper.GetString("server.log-output") == "file" {
			w := log.NewSyncWriter(logfile)
			logger = log.NewLogfmtLogger(w)
		} else {
			logger = log.NewLogfmtLogger(os.Stderr)
		}

		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

	// Init DB Connection
	db, err := database.NewConnectionDB(viper.GetString("database.driver"), viper.GetString("database.dbname"), viper.GetString("database.host"), viper.GetString("database.username"), viper.GetString("database.password"), viper.GetInt("database.port"))
	if err != nil {
		logger.Log("Err Db connection :", err.Error())
		panic(err.Error())
	}
	db.AutoMigrate(&entity.Product{})

	logger.Log("message", "Connection Db Success")

	// Consul initialization
	registar := consul.ConsulRegisterService(viper.GetString("server.service-name"), viper.GetInt("server.port"), logger)
	registar.Register()
	defer registar.Deregister()

	// service registry
	prodSvc := registry.RegisterProductService(db, logger)

	// transport init
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP"))
	prodHttp := transport.ProductHttpHandler(prodSvc, log.With(logger, "ProductTransportLayer", "HTTP"))

	//Routing path
	mux := http.NewServeMux()
	mux.Handle("/swagger/", swagHttp)
	mux.Handle("/kd/v1/", prodHttp)
	//http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)

	// // configure hystrix
	// var prescriptionEndpoint endpoint.Endpoint
	// hystrix.ConfigureCommand("prescription Request", hystrix.CommandConfig{Timeout: 1000})
	// prescriptionEndpoint = Hystrix("Prescription Request", "Service currently unavailable", logger)(prescriptionEndpoint)
	// hystrixStreamHandler := hystrix.NewStreamHandler()
	// hystrixStreamHandler.Start()
	// go func() {
	// 	errs <- http.ListenAndServe(net.JoinHostPort("", "9000"), hystrixStreamHandler)
	// }()

	go func() {
		logger.Log("transport", "HTTP", "addr", viper.GetInt("server.port"))
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("server.port")), mux)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
