package psql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/domain/balance"
	"github.com/swooosh13/avito-test/pkg/postgres"
)

type BalanceRepository struct {
	db     *postgres.Client
	logger *zerolog.Logger
}

func NewBalanceRepository(db *postgres.Client, logger *zerolog.Logger) *BalanceRepository {
	return &BalanceRepository{
		db:     db,
		logger: logger,
	}
}

func (r *BalanceRepository) GetBalanceByID(ctx context.Context, id int64) (int64, error) {
	sql := `
SELECT amount FROM balance WHERE user_id = $1;
`
	r.logger.Debug().Msgf("sql: %s", sql)

	var amount int64
	if err := r.db.Pool.QueryRow(ctx, sql, id).Scan(&amount); err != nil {
		return 0, fmt.Errorf("get balance by id: %w", err)
	}

	return amount, nil
}

func (r *BalanceRepository) UpdateBalance(ctx context.Context, dto *balance.UpdateBalanceDTO) error {
	sql := `
INSERT INTO balance (user_id, amount) VALUES ($1, $2) 
ON CONFLICT (user_id) DO update SET amount = balance.amount + $2;
	`
	r.logger.Debug().Msgf("sql: %s", sql)

	_, err := r.db.Pool.Exec(ctx, sql, dto.UserID, dto.Amount)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}

func (r *BalanceRepository) TransferFromTo(ctx context.Context, dto *balance.TransferFromToDTO) (err error) {
	sql := `
call transfer($1, $2, $3);
`
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
	_, err = tx.Exec(ctx, sql, dto.FromUserID, dto.ToUserID, dto.Amount)
	if err != nil {
		return fmt.Errorf("transfer from to: %w", err)
	}

	return nil
}
