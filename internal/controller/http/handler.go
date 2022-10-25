package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/domain"
)

type Handler struct {
	services *domain.Services
	logger   *zerolog.Logger
}

func NewHandler(services *domain.Services, logger *zerolog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1") // authMiddleware, permissions
	{
		h.InitReservationAPI(v1)
		h.InitTransactionsAPI(v1)
		h.InitBalanceAPI(v1)
	}
}
