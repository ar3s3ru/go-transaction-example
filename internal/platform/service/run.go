package service

import (
	"database/sql"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
	"go.uber.org/zap"

	"go-transaction-example/internal/app"
	"go-transaction-example/internal/pkg/must"
	srvhttp "go-transaction-example/internal/platform/http"
	"go-transaction-example/internal/platform/postgres"
)

func Run(config Config) error {
	db, err := sql.Open("postgres", config.DB.DSN())
	must.NotFail(err)

	dbx := sqlx.NewDb(db, "postgres")

	logger := zap.L()

	usersRepository := postgres.UserRepository{}
	userHistoryRepository := postgres.UserHistory{
		DB:    dbx,
		Adder: usersRepository,
	}
	transactional := postgres.Transactional{
		DB:    dbx,
		Adder: userHistoryRepository,
	}

	srv := app.Service{
		Adder:  transactional,
		Lister: userHistoryRepository,
	}

	router := srvhttp.Router(srv, logger)
	return http.ListenAndServe(config.Server.Host(), router)
}
