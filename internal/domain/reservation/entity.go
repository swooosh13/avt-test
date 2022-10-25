package reservation

type Reservation struct {
	UserID int64 `json:"user_id"`
	Amount int64 `json:"amount"`
}
