package transactions

import "time"

type RevenueDTO struct {
	UserID    int64 `json:"user_id" binding:"required"`
	ServiceID int64 `json:"service_id" binding:"required"`
	OrderID   int64 `json:"order_id" binding:"required"`
	Amount    int64 `json:"amount" binding:"required"`
}

type GetUserTransactionDTO struct {
	UserID    int64     `json:"user_id"`
	ServiceID *int64    `json:"service_id"`
	OrderID   *int64    `json:"order_id"`
	FromInfo  string    `json:"from_info"`
	ToInfo    string    `json:"to_info"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type ReportDTO struct {
	Period time.Time `json:"period"`
}
