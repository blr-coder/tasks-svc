package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ITransaction interface {
	Finish(context.Context) error
	Rollback(context.Context) error
	GetTx() *sqlx.Tx
}

type IDBTransactionManager interface {
	StartTx(context.Context) (ITransaction, error)
}

type DbTransaction struct {
	Tx *sqlx.Tx
}

func (tx DbTransaction) Finish(ctx context.Context) (err error) {
	if err = tx.Tx.Commit(); err != nil {
		return fmt.Errorf("can not commit tx: %w", err)
	}
	return
}

func (tx DbTransaction) Rollback(ctx context.Context) (err error) {
	if err = tx.Tx.Rollback(); err != nil {
		return fmt.Errorf("can not rollback tx: %w", err)
	}
	return
}

func (tx DbTransaction) GetTx() *sqlx.Tx {
	return tx.Tx
}

type DBTransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) DBTransactionManager {
	return DBTransactionManager{db: db}
}

func (tm DBTransactionManager) StartTx(ctx context.Context) (ITransaction, error) {
	var tx, err = tm.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("can not start tx: %w", err)
	}
	return DbTransaction{Tx: tx}, nil
}
