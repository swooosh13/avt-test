package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swooosh13/avito-test/internal/domain/transactions"
)

func (h *Handler) InitTransactionsAPI(api *gin.RouterGroup) {
	transactions := api.Group("/transactions")
	{
		transactions.GET("/", h.GetReport)
		transactions.POST("/", h.Revenue)
		transactions.GET("/:user_id", h.UserReport)
	}
}

type userTransactionsResponse struct {
	Transactions []transactions.GetUserTransactionDTO `json:"transactions"`
	Range        listRange                            `json:"range"`
}

func (h *Handler) UserReport(ctx *gin.Context) {
	paginationParams := parsePaginationParams(ctx)
	userID, err := parseID(ctx, "user_id")
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid user_id param", err)
		return
	}

	listParams := transactions.NewListParams(paginationParams)
	trs, count, err := h.services.Transactions.GetTransactionsByUser(ctx.Request.Context(), userID, listParams)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, "failed to get transactions", err)
		return
	}

	ctx.JSON(200, response{
		Data: userTransactionsResponse{
			Transactions: trs,
			Range: listRange{
				Count:  count,
				Limit:  int(listParams.Pagination.Limit),
				Offset: int(listParams.Pagination.Offset),
			},
		},
	})
}

func (h *Handler) GetReport(ctx *gin.Context) {
	date, ok := ctx.GetQuery("date")
	if !ok {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid date param in query", nil)
		return
	}

	period, err := time.Parse("2006-01", date)
	if err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid date format in query", err)
		return
	}

	reportDto := transactions.ReportDTO{Period: period}
	report, err := h.services.Transactions.Report(ctx.Request.Context(), &reportDto)
	if err != nil {
		h.errorResponse(ctx, http.StatusInternalServerError, "failed to get report", err)
		return
	}

	ctx.JSON(200, response{
		Data: gin.H{
			"link": report,
		},
	})
}

func (h *Handler) Revenue(ctx *gin.Context) {
	var dto transactions.RevenueDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		h.errorResponse(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	err := h.services.Transactions.Revenue(ctx.Request.Context(), &dto)
	if err != nil {
		// разрезервировать

		// reservationBalanceDto := reservation.ReservationBalanceDTO{
		//	UserID:    dto.UserID,
		//	ServiceID: dto.ServiceID,
		//	Amount:    -dto.Amount,
		//	OrderID:   dto.UserID,
		// }
		// err := h.services.Reservation.ReserveBalance(ctx.Request.Context(), &reservationBalanceDto)
		// if err != nil {
		//	 h.errorResponse(ctx, http.StatusInternalServerError, "failed to de-reserve funds", err)
		//	 return
		// }

		h.errorResponse(ctx, http.StatusInternalServerError, "failed to revenue, de-reserve funds", err)
		return
	}

	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}
