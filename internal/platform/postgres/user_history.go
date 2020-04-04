package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"

	"db-transaction-example/internal/domain/user"
	"db-transaction-example/internal/domain/user/history"
)

type UserHistory struct {
	DB *sqlx.DB

	user.Adder
}

func (r UserHistory) List(ctx context.Context, from, to time.Time) ([]history.Entry, error) {
	var entries []history.Entry

	err := r.DB.SelectContext(ctx, &entries, `
		SELECT * FROM users_history
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at DESC`,
		from,
		to,
	)

	if err != nil {
		err = fmt.Errorf("postgres.UserHistory: failed to list history %w", err)
	}

	return entries, err
}

func (r UserHistory) Add(ctx context.Context, state user.State) (user.Entity, error) {
	var entity user.Entity
	var err error

	err = WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		entity, err = r.Adder.Add(ctx, state)
		if err != nil {
			return err
		}

		jsonState, err := json.Marshal(entity.State)
		if err != nil {
			return fmt.Errorf("postgres.UserHistory: failed to json marshal User state %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO users_history (user_id, action, state)
			VALUES ($1, $2, $3);`,
			entity.ID,
			history.WasCreatedAction,
			types.JSONText(jsonState),
		)

		if err != nil {
			err = fmt.Errorf("postgres.UserHistory: failed to apped to history %w", err)
		}

		return err
	})

	return entity, err
}
