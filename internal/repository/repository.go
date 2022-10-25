package repository

import (
	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/repository/psql"
	"github.com/swooosh13/avito-test/internal/repository/s3"
	"github.com/swooosh13/avito-test/pkg/postgres"
)

type Repositories struct {
	Balance      *psql.BalanceRepository
	Reservation  *psql.ReservationRepository
	Transactions *psql.TransactionRepository
	S3Client     *s3.Client
}

func NewRepositories(db *postgres.Client, logger *zerolog.Logger, s3Client *s3.Client) *Repositories {
	return &Repositories{
		Balance:      psql.NewBalanceRepository(db, logger),
		Reservation:  psql.NewReservationRepository(db, logger),
		Transactions: psql.NewTransactionRepository(db, logger),
		S3Client:     s3Client,
	}
}
