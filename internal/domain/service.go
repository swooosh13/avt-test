package domain

import (
	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/domain/balance"
	"github.com/swooosh13/avito-test/internal/domain/reservation"
	"github.com/swooosh13/avito-test/internal/domain/transactions"
	"github.com/swooosh13/avito-test/internal/repository"
)

type Services struct {
	Balance      *balance.Service
	Reservation  *reservation.Service
	Transactions *transactions.Service
}

func NewServices(repositories *repository.Repositories, logger *zerolog.Logger) *Services {
	return &Services{
		Balance:      balance.NewService(repositories.Balance, logger),
		Reservation:  reservation.NewService(repositories.Reservation, logger),
		Transactions: transactions.NewService(repositories.Transactions, repositories.S3Client, logger),
	}
}
