package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/config"
	"github.com/swooosh13/avito-test/internal/controller/http"
	"github.com/swooosh13/avito-test/internal/domain"
)

func NewRouter(services *domain.Services, logger *zerolog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	if config.Get().Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	v1 := rest.NewHandler(services, logger)
	api := router.Group("/api")
	{
		v1.InitRoutes(api)
	}

	return router
}
