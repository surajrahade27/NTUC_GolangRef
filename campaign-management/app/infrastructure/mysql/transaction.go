package mysql

import (
	"campaign-mgmt/app/domain/services"
	"context"
	"fmt"

	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ContextKey ...
type ContextKey string

const (
	ctxKeyMYSQLTx ContextKey = "mysql-tx"
)

// TransactionService manages DB transactions
type TransactionService struct {
	pool *gorm.DB
}

// NewTransactionService ..
func NewTransactionService(pool *gorm.DB) *TransactionService {
	return &TransactionService{
		pool: pool,
	}
}

// Begin Will start a tx in database and return a dbtx in context
func (ts *TransactionService) Begin(ctx context.Context) (context.Context, error) {
	tx := ts.pool.Begin()
	ctx = WithDBTransaction(ctx, tx)
	return ctx, nil
}

// Commit will be used to Commit a tx
func (ts *TransactionService) Commit(ctx context.Context) error {
	db := ts.getTransaction(ctx)
	if db == nil {
		return services.ErrTxNotFound
	}

	return db.Commit().Error
}

// getTransaction is function to get a transaction
// its caller responsibility to check whether tx object is nil
func (ts *TransactionService) getTransaction(ctx context.Context) *gorm.DB {
	return DBTransaction(ctx)
}

// Rollback will be used to Rollback a tx
func (ts *TransactionService) Rollback(ctx context.Context) error {
	db := ts.getTransaction(ctx)
	if db == nil {
		return services.ErrTxNotFound
	}

	return db.Rollback().Error
}

// RunWithTransaction will start a tx and will call the function
func (ts *TransactionService) RunWithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, err := ts.Begin(ctx)
	if err != nil {
		return err
	}
	if err := fn(ctx); err != nil {
		logger.Errorf("error occured : %s, transaction rolling back", err.Error())
		rollbackErr := ts.Rollback(ctx)
		if rollbackErr != nil {
			err = fmt.Errorf("%w:%v", rollbackErr, err.Error())
		}
		return err
	}
	return ts.Commit(ctx)
}

// WithDBTransaction returns a context with MYSQL transaction
func WithDBTransaction(ctx context.Context, value *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyMYSQLTx, value)
}

// DBTransaction returns MYSQL transaction
func DBTransaction(ctx context.Context) *gorm.DB {
	value, ok := ctx.Value(ctxKeyMYSQLTx).(*gorm.DB)
	if !ok {
		return nil
	}
	return value
}
