package reservation

type BalanceReserveDTO struct {
	UserID    int64 `json:"user_id" binding:"required"`
	ServiceID int64 `json:"service_id" binding:"required"`
	OrderID   int64 `json:"order_id" binding:"required"`
	Amount    int64 `json:"amount" binding:"required"`
}
