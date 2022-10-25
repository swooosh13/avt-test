package psql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/domain/transactions"
	"github.com/swooosh13/avito-test/pkg/postgres"
)

type TransactionRepository struct {
	db     *postgres.Client
	logger *zerolog.Logger
}

func NewTransactionRepository(db *postgres.Client, logger *zerolog.Logger) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TransactionRepository) Revenue(ctx context.Context, dto *transactions.RevenueDTO) (err error) {
	sql := `
call remove_reserve_from_user($1, $2, $3, $4);
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
	_, err = tx.Exec(ctx, sql, dto.UserID, dto.ServiceID, dto.OrderID, dto.Amount)
	if err != nil {
		return fmt.Errorf("revenue from reservation: %w", err)
	}

	return nil
}

func (r *TransactionRepository) Report(ctx context.Context, dto *transactions.ReportDTO) ([]transactions.TransactReport, error) {
	sql := `
select service_id, sum(amount) from transactions where to_info = 'payment' 
and date_trunc('month', created_at) = date_trunc('month', $1::timestamp) group by service_id ;
`
	r.logger.Debug().Msgf("sql: %s", sql)

	rows, err := r.db.Pool.Query(ctx, sql, dto.Period)
	if err != nil {
		return nil, fmt.Errorf("run get all news query: %w", err)
	}
	defer rows.Close()

	var reports []transactions.TransactReport

	for rows.Next() {
		var entity transactions.TransactReport
		if err := rows.Scan(
			&entity.ServiceID,
			&entity.TotalAmount,
		); err != nil {
			return nil, fmt.Errorf("scan get report: %w", err)
		}

		reports = append(reports, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read all reports: %w", err)
	}

	return reports, nil
}

func (r *TransactionRepository) GetTransactionsByUser(ctx context.Context, id int64, params transactions.ListParams) ([]transactions.GetUserTransactionDTO, int, error) { //nolint:lll
	query := r.db.Builder.
		Select("user_id", "service_id", "order_id", "amount", "from_info", "to_info", "created_at", "count(*) over () as total").
		From("transactions").
		Where(sq.Eq{"user_id": id}).
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset)

	for _, sort := range params.Sorts {
		query = sort.UseSelectBuilder(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build query: %w", err)
	}
	r.logger.Debug().Msgf("sql: %s", sql)

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("run get all news query: %w", err)
	}
	defer rows.Close()

	var trs []transactions.GetUserTransactionDTO
	var count int
	for rows.Next() {
		var entity transactions.GetUserTransactionDTO
		err := rows.Scan(
			&entity.UserID,
			&entity.ServiceID,
			&entity.OrderID,
			&entity.Amount,
			&entity.FromInfo,
			&entity.ToInfo,
			&entity.CreatedAt,
			&count)
		if err != nil {
			return nil, 0, fmt.Errorf("scan get transaction: %w", err)
		}
		trs = append(trs, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("read all transactions: %w", err)
	}

	return trs, count, nil
}
