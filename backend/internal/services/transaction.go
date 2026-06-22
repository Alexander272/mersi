package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
	"github.com/Alexander272/mersi/backend/pkg/logger"
)

type TransactionManagerService struct {
	repo repository.Transaction
}

func NewTransactionManager(repo repository.Transaction) *TransactionManagerService {
	return &TransactionManagerService{repo: repo}
}

type TransactionManager interface {
	ExecuteInTx(ctx context.Context, fn func(tx postgres.Tx) error) error
}

func (tm *TransactionManagerService) ExecuteInTx(ctx context.Context, fn func(tx postgres.Tx) error) error {
	tx, err := tm.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				logger.Error("failed to rollback transaction on panic", logger.StringAttr("error", rbErr.Error()))
			}
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				logger.Error("failed to rollback transaction on error", logger.StringAttr("error", rbErr.Error()))
			}
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
