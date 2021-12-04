//  KD Boilerplate Go kit:
//   version: 1.0
//   title: Boilerplate Go Kit Api
//  Schemes: http
//  Host: localhost:5600
//  BasePath: /
//  Produces:
//    - application/json
//
// swagger:meta
package main

import (
	"fmt"
	"go-klikdokter/app/api/initialization"
	"go-klikdokter/helper/consul"
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
	db, err := initialization.DbInit()
	if err != nil {
		logger.Log("Err Db connection :", err.Error())
		panic(err.Error())
	}
	logger.Log("message", "Connection Db Success")

	// Consul initialization
	registar := consul.ConsulRegisterService(viper.GetString("server.service-name"), viper.GetInt("server.port"), logger)
	registar.Register()
	defer registar.Deregister()

	// Routing initialization
	mux := initialization.InitRouting(db, logger)
	http.Handle("/", accessControl(mux))

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
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("server.port")), nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errs)
}

//func handleSw() http.Handler {
//	pr := mux.NewRouter()
//	pr.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
//	//opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
//	//sh := middleware.SwaggerUI(opts, nil)
//	//pr.Handle("/docs", sh)
//
//	return pr
//}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
