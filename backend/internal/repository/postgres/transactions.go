package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TransactionRepo struct {
	db *sqlx.DB
}

func NewTransactionRepo(db *sqlx.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

type Transaction interface {
	BeginTx(ctx context.Context) (Tx, error)
	getExec(tx Tx) QueryExecutor
}

type Tx interface {
	TX() *sqlx.Tx
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Результирующий интерфейс для всех SQL-операций
type QueryExecutor interface {
	sqlx.ExtContext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// Можно добавить GetContext, SelectContext и т.д.
}

// Метод для начала транзакции
func (r *TransactionRepo) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction. error: %w", err)
	}
	return &repoTx{Tx: tx}, nil
}

func (r *TransactionRepo) getExec(tx Tx) QueryExecutor {
	if tx != nil {
		return tx.TX()
	}
	return r.db
}

type repoTx struct {
	Tx *sqlx.Tx
}

func (t *repoTx) TX() *sqlx.Tx {
	return t.Tx
}

func (t *repoTx) Commit(ctx context.Context) error {
	if err := t.Tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. error: %w", err)
	}
	return nil
}

func (t *repoTx) Rollback(ctx context.Context) error {
	if err := t.Tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction. error: %w", err)
	}
	return nil
}
