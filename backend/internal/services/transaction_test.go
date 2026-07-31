package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

func TestExecuteInTx_CommitsOnSuccess(t *testing.T) {
	tx := &fakeTx{}
	repo := &fakeRepoTx{tx: tx}
	tm := NewTransactionManager(repo)

	err := tm.ExecuteInTx(context.Background(), func(tx postgres.Tx) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tx.committed {
		t.Fatal("expected transaction to be committed")
	}
	if tx.rolledBack {
		t.Fatal("expected no rollback on success")
	}
}

func TestExecuteInTx_RollsBackOnError(t *testing.T) {
	tx := &fakeTx{}
	repo := &fakeRepoTx{tx: tx}
	tm := NewTransactionManager(repo)

	fnErr := errors.New("query failed")
	err := tm.ExecuteInTx(context.Background(), func(tx postgres.Tx) error {
		return fnErr
	})
	if !errors.Is(err, fnErr) {
		t.Fatalf("expected original error, got %v", err)
	}
	if !tx.rolledBack {
		t.Fatal("expected transaction to be rolled back")
	}
	if tx.committed {
		t.Fatal("expected no commit on error")
	}
}

func TestExecuteInTx_BeginError(t *testing.T) {
	beginErr := errors.New("cannot begin")
	repo := &fakeRepoTx{beginErr: beginErr}
	tm := NewTransactionManager(repo)

	err := tm.ExecuteInTx(context.Background(), func(tx postgres.Tx) error {
		t.Fatal("fn should not be called when BeginTx fails")
		return nil
	})
	if !errors.Is(err, beginErr) {
		t.Fatalf("expected begin error, got %v", err)
	}
}

func TestExecuteInTx_CommitError(t *testing.T) {
	commitErr := errors.New("commit failed")
	tx := &commitFailingTx{fakeTx: &fakeTx{}, err: commitErr}
	repo := &fakeRepoTx{tx: tx}
	tm := NewTransactionManager(repo)

	err := tm.ExecuteInTx(context.Background(), func(tx postgres.Tx) error {
		return nil
	})
	if !errors.Is(err, commitErr) {
		t.Fatalf("expected commit error, got %v", err)
	}
}

type commitFailingTx struct {
	*fakeTx
	err error
}

func (f *commitFailingTx) Commit(context.Context) error { return f.err }
