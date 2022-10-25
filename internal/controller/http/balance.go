package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swooosh13/avito-test/internal/domain/balance"
)

func (h *Handler) InitBalanceAPI(api *gin.RouterGroup) {
	balance := api.Group("/balance")
	{
		balance.GET("/:id", h.GetBalanceByID)
		balance.PUT("/", h.UpdateBalance)
		balance.POST("/", h.TransferBalance)
	}
}

type getBalanceByIDResponse struct {
	Balance int64 `json:"balance"`
}

func (h *Handler) GetBalanceByID(ctx *gin.Context) {
	id, err := parseID(ctx, "id")
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid id param", err)
		return
	}

	balance, err := h.services.Balance.GetBalanceByID(ctx.Request.Context(), id)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, "failed to get balance", err)
		return
	}

	ctx.JSON(http.StatusOK, response{
		Data: getBalanceByIDResponse{
			Balance: balance,
		},
	})
}

func (h *Handler) UpdateBalance(ctx *gin.Context) {
	var dto balance.UpdateBalanceDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.services.Balance.UpdateBalance(ctx.Request.Context(), &dto); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, "failed to update balance", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) TransferBalance(ctx *gin.Context) {
	var dto balance.TransferFromToDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.services.Balance.TransferFromTo(ctx.Request.Context(), &dto); err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, "failed to transfer balance", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
