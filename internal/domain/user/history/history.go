package history

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx/types"

	"db-transaction-example/internal/domain"
	"db-transaction-example/internal/domain/user"
)

type Action string

const (
	WasCreatedAction = Action("wasCreated")
)

type Entry struct {
	domain.Created

	ID     int64          `json:"id" db:"user_history_id"`
	UserID user.ID        `json:"userId" db:"user_id"`
	Action Action         `json:"action" db:"action"`
	State  types.JSONText `json:"state" db:"state"`
}

type Lister interface {
	List(ctx context.Context, from, to time.Time) ([]Entry, error)
}
