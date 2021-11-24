package main

import (
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

	//2. Load configuration
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

	//3. Logging init
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
		logger.Log("transport", "HTTP", "addr", viper.GetInt("server.port"))
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("server.port")), h)
	}()

	logger.Log("exit", <-errs)
}
