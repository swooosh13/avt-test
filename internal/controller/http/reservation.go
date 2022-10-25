package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swooosh13/avito-test/internal/domain/reservation"
)

func (h *Handler) InitReservationAPI(api *gin.RouterGroup) {
	reservation := api.Group("/reservation")
	{
		reservation.GET("/:id", h.GetReservation)
		reservation.POST("/", h.TransferToReservation)
	}
}

type getReservationByIDResponse struct {
	Balance int64 `json:"balance"`
}

func (h *Handler) GetReservation(ctx *gin.Context) {
	id, err := parseID(ctx, "id")
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid id param", err)
		return
	}

	reservationBalance, err := h.services.Reservation.GetReservationByID(ctx.Request.Context(), id)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, "failed to get reservation balance", err)
		return
	}
	ctx.JSON(http.StatusOK, getReservationByIDResponse{
		Balance: reservationBalance.Amount,
	})
}

func (h *Handler) TransferToReservation(ctx *gin.Context) {
	var dto reservation.BalanceReserveDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.services.Reservation.ReserveBalance(ctx.Request.Context(), &dto); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, "failed to transfer to reservation", err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}
