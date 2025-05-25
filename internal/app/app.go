package app

import (
	"log"

	_ "github.com/ferdikurniawan/loan-service/docs"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggofiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"

	"github.com/ferdikurniawan/loan-service/internal/pkg/postgres"
	grace "github.com/ferdikurniawan/loan-service/internal/utils/grace"

	"github.com/ferdikurniawan/loan-service/config"
	v1 "github.com/ferdikurniawan/loan-service/internal/controller/http/v1"
	"github.com/ferdikurniawan/loan-service/internal/repo"
	"github.com/ferdikurniawan/loan-service/internal/services"
)

const (
	ServiceName = "loan-service"
)

func Run(config *config.Config) {

	// DB
	pgCfg := &postgres.Config{
		Dsn:     config.PostgreHost,
		MaxConn: config.DBMaxOpenConnection,
		MaxIdle: config.DBMaxIdleConnection,
	}

	pg, err := postgres.New(pgCfg)
	if err != nil {
		log.Fatalf("error init postgres %s", err.Error())
	}

	// services layer
	loanService := services.NewLoanService(repo.NewLoanRepo(pg))

	// gin
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()

	// middlewares
	handler.Use(gintrace.Middleware(ServiceName))
	handler.Use(gin.Logger())
	handler.Use(gzip.Gzip(gzip.DefaultCompression))
	handler.Use(gin.Recovery())

	// swagger
	handler.GET("/swagger/*any", ginswagger.WrapHandler(swaggofiles.Handler))

	v1.NewRouter(handler, v1.Services{
		Cfg:         config,
		LoanService: loanService,
	})

	grace.Serve(config.Port, handler)
}
