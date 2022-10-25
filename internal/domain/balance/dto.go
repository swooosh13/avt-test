package balance

type UpdateBalanceDTO struct {
	UserID int64 `json:"user_id" binding:"required"`
	Amount int64 `json:"amount" binding:"required"`
}

type TransferFromToDTO struct {
	FromUserID int64 `json:"from_user_id" binding:"required"`
	ToUserID   int64 `json:"to_user_id" binding:"required"`
	Amount     int64 `json:"amount" binding:"required"`
}
