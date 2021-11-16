package main

import (
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"gokit_example/pkg/common"
	"gokit_example/pkg/repository"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gokit_example/pkg/service"
	"gokit_example/pkg/transport"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//1. Flag parse
	var (
		httpAddr          = flag.Int("http.addr", 8080, "HTTP Port listen address")
		profile           = flag.String("profile", "dev", "Profile app environment")
		logOutput         = flag.String("log.output", "console", "Output log")
		logOutputFilePath = flag.String("log.output.file.path", ".", "Output log file path")
		configPath        = flag.String("config.path", ".", "Config file path")
		serviceName       = flag.String("service.name", "prescription", "Name of instance")
		host              = flag.String("host", "127.0.0.1:8080", "Host this service")
	)
	flag.Parse()

	//2. Load configuration
	viper.SetConfigType("yaml")
	var configFileName []string
	configFileName = append(configFileName, "config-")
	configFileName = append(configFileName, *profile)
	viper.SetConfigName(strings.Join(configFileName, ""))
	viper.AddConfigPath(*configPath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}

	//3. Logging init
	logfile, err := os.OpenFile(*logOutputFilePath+"/prescription.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer logfile.Close()
	var logger log.Logger
	{
		if *logOutput == "file" {
			w := log.NewSyncWriter(logfile)
			logger = log.NewLogfmtLogger(w)
		} else {
			logger = log.NewLogfmtLogger(os.Stderr)
		}

		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

	//4. Init postgresql db
	db, err := common.NewConnectionDB(viper.GetString("database.driver"), viper.GetString("database.database"), viper.GetString("database.host"), viper.GetString("database.username"), viper.GetString("database.password"), viper.GetInt("database.port"))
	if err != nil {
		logger.Log("Err Db connection :", err.Error())
		panic(err.Error())
	}

	logger.Log("message", "Connection Db Success")

	//6. Init Redis
	// var resulRedis redis.Client
	rds := common.NewConnectionRedis(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("cache.redis.hostname"), viper.GetInt("cache.redis.port")),
		Password: viper.GetString("cache.redis.password"), // no password set
		DB:       viper.GetInt("cache.redis.db"),          // use default DB
	})

	// 7. Register cd Specify the information of an instance.
	asr := api.AgentServiceRegistration{
		// Every service instance must have an unique ID.
		ID:   fmt.Sprintf("%v", *host),
		Name: *serviceName,
		// These two values are the location of an instance.
		Address: *host,
		Port:    *httpAddr,
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

	// init repository
	repo := repository.NewProductRepository(db, logger)

	var s service.Service
	{
		s = service.NewServiceImplV1(logger, repo, rds)
		s = service.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = transport.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)

	// configure hystrix
	var prescriptionEndpoint endpoint.Endpoint
	hystrix.ConfigureCommand("prescription Request", hystrix.CommandConfig{Timeout: 1000})
	prescriptionEndpoint = Hystrix("Prescription Request", "Service currently unavailable", logger)(prescriptionEndpoint)
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		errs <- http.ListenAndServe(net.JoinHostPort("", "9000"), hystrixStreamHandler)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", *httpAddr), h)
	}()

	logger.Log("exit", <-errs)
}
