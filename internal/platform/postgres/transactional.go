package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"go-transaction-example/internal/domain/user"
)

type transactionKey struct{}

func WithTransaction(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) error {
	tx := ctx.Value(transactionKey{}).(*sqlx.Tx)
	return fn(ctx, tx)
}

type Transactional struct {
	DB    *sqlx.DB
	Adder user.Adder
}

func (t Transactional) Add(ctx context.Context, state user.State) (user.Entity, error) {
	var user user.Entity
	var err error

	err = t.exec(ctx, func(ctx context.Context) error {
		user, err = t.Adder.Add(ctx, state)
		return err
	})

	return user, err
}

func (t Transactional) exec(ctx context.Context, fn func(context.Context) error) error {
	tx, err := t.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres.Transactional: failed to start transaction %w", err)
	}

	defer func() {
		if err == nil {
			return
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("postgres.WithTransaction: failed to rollback (%s) %w", rollbackErr.Error(), err)
		}
	}()

	ctx = context.WithValue(ctx, transactionKey{}, tx)

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("postgres.Transactional: failed to execute transaction %w", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("postgres.WithCommit: failed to commit transaction %w", err)
	}

	return err
}
