package services

import (
	"campaign-mgmt/app/domain/valueobjects"
	"context"
)

// TransactionService errors
const (
	// ErrTxNotFound is error returned if transaction is not found in context
	ErrTxNotFound valueobjects.Error = "transaction not found"
	// ErrTxAlreadyStarted if transaction already started
	ErrTxAlreadyStarted valueobjects.Error = "transaction already started"
)

// TransactionService ..
type TransactionService interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	RunWithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
