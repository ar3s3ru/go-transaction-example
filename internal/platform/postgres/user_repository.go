package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"go-transaction-example/internal/app"
	"go-transaction-example/internal/domain/user"
)

type UserRepository struct{}

func (r UserRepository) Add(ctx context.Context, state user.State) (user.Entity, error) {
	var entity user.Entity
	var err error

	err = WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &entity,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING *;",
			state.Name, state.Age,
		)

		if IsAlreadyExistsError(err) {
			err = app.AlreadyExistsError{Inner: err}
		}

		if err != nil {
			err = fmt.Errorf("postgres.UserRepository: failed to insert user %w", err)
		}

		return err
	})

	return entity, err
}
