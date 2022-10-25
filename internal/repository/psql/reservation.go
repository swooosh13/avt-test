package psql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/domain/reservation"
	"github.com/swooosh13/avito-test/pkg/postgres"
)

type ReservationRepository struct {
	db     *postgres.Client
	logger *zerolog.Logger
}

func NewReservationRepository(db *postgres.Client, logger *zerolog.Logger) *ReservationRepository {
	return &ReservationRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ReservationRepository) GetReservationByID(ctx context.Context, id int64) (reservation.Reservation, error) {
	sql := `
select user_id, amount from reservation where user_id = $1;
`
	var res reservation.Reservation
	if err := r.db.Pool.QueryRow(ctx, sql, id).Scan(
		&res.UserID,
		&res.Amount,
	); err != nil {
		return reservation.Reservation{}, fmt.Errorf("get reservation by id: %w", err)
	}

	return res, nil
}

func (r *ReservationRepository) ReserveBalance(ctx context.Context, dto *reservation.BalanceReserveDTO) (err error) {
	sql := `
call reserve_from_user($1, $2, $3, $4);
`

	if dto.Amount < 0 {
		sql = `call reserve_to_user($1, $2, $3, $4);`
		dto.Amount = -dto.Amount
	}

	r.logger.Debug().Msgf("sql: %s", sql)

	tx, err := r.db.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()
	_, err = tx.Exec(ctx, sql, dto.UserID, dto.ServiceID, dto.OrderID, dto.Amount)
	if err != nil {
		return fmt.Errorf("transfer from to: %w", err)
	}

	return nil
}
