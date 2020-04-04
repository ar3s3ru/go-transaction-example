package app

import (
	"context"
	"time"

	"go-transaction-example/internal/domain/user"
	"go-transaction-example/internal/domain/user/history"
)

type Service struct {
	Adder  user.Adder
	Lister history.Lister
}

func (srv Service) AddUser(ctx context.Context, name string, age uint8) (user.Entity, error) {
	return srv.Adder.Add(ctx, user.State{
		Name: name,
		Age:  age,
	})
}

func (srv Service) ListHistory(ctx context.Context, from, to time.Time) ([]history.Entry, error) {
	return srv.Lister.List(ctx, from, to)
}
