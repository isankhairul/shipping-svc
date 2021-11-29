package main

import (
	"fmt"
	"gokit_example/app/api/transport"
	"gokit_example/app/registry"
	"gokit_example/helper/database"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//Load configuration
	viper.SetConfigType("yaml")
	var profile string = "dev"
	if os.Getenv("env") != "" {
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
	db, err := database.NewConnectionDB(viper.GetString("database.driver"), viper.GetString("database.database"), viper.GetString("database.host"), viper.GetString("database.username"), viper.GetString("database.password"), viper.GetInt("database.port"))
	if err != nil {
		logger.Log("Err Db connection :", err.Error())
		panic(err.Error())
	}

	logger.Log("message", "Connection Db Success")

	//6. Init Redis
	// var resulRedis redis.Client
	// rds := common.NewConnectionRedis(&redis.Options{
	// 	Addr:     fmt.Sprintf("%s:%d", viper.GetString("cache.redis.hostname"), viper.GetInt("cache.redis.port")),
	// 	Password: viper.GetString("cache.redis.password"), // no password set
	// 	DB:       viper.GetInt("cache.redis.db"),          // use default DB
	// })

	// Register cd Specify the information of an instance.
	host, _ := os.Hostname()
	asr := api.AgentServiceRegistration{
		// Every service instance must have an unique ID.
		ID:   fmt.Sprintf("%v", host),
		Name: viper.GetString("server.service-name"),
		// These two values are the location of an instance.
		Address: host,
		Port:    viper.GetInt("server.port"),
	}
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	sdClient := consul.NewClient(consulClient)
	registar := consul.NewRegistrar(sdClient, &asr, logger)
	registar.Register()
	// According to the official doc of Go kit,
	// it's important to call registar.Deregister() before the program exits.
	defer registar.Deregister()

	// service registry
	prodSvc := registry.NewProductService(db, logger)

	// transport init
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP"))
	prodHttp := transport.ProductHttpHandler(prodSvc, log.With(logger, "ProductTransportLayer", "HTTP"))

	//Path
	mux := http.NewServeMux()

	mux.Handle("/swagger/v1/", swagHttp)
	mux.Handle("/boilerplate/v1/", prodHttp)
	http.Handle("/", accessControl(mux))

	errs := make(chan error)

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
