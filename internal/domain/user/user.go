package user

import (
	"context"

	"db-transaction-example/internal/domain"
)

type ID int64

type State struct {
	Name string `json:"name" db:"name"`
	Age  uint8  `json:"age" db:"age"`
}

type Entity struct {
	ID ID `json:"id,omitempty" db:"user_id"`

	domain.Created
	domain.Updated
	State
}

type Adder interface {
	Add(context.Context, State) (Entity, error)
}
