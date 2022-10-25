package transactions

import (
	"time"

	"github.com/swooosh13/avito-test/pkg/pagination"
	"github.com/swooosh13/avito-test/pkg/sort"
)

type Transaction struct {
	UserID    int64     `json:"user_id"`
	ServiceID int64     `json:"service_id"`
	OrderID   int64     `json:"order_id"`
	FromInfo  string    `json:"from_info"`
	ToInfo    string    `json:"to_info"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactReport struct {
	ServiceID   int64 `json:"service_id"`
	TotalAmount int64 `json:"amount"`
}

type ListParams struct {
	Pagination pagination.Params
	Sorts      []sort.Sort
}

func NewListParams(params pagination.Params) ListParams {
	if params.Limit == 0 || params.Limit > pagination.MaxLimit {
		params.Limit = pagination.DefaultLimit
	}

	return ListParams{
		Sorts: []sort.Sort{
			sort.New("amount", sort.Desc),
			sort.New("created_at", sort.Desc),
		},
		Pagination: params,
	}
}
