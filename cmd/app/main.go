package main

import (
	"database/sql"
	"net/http"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"go-transaction-example/internal/app"
	srvhttp "go-transaction-example/internal/platform/http"
	"go-transaction-example/internal/platform/postgres"
	"go-transaction-example/internal/platform/service"
)

func mustNotFail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config, err := service.ParseConfig()
	mustNotFail(err)

	db, err := sql.Open("postgres", config.DB.DSN())
	mustNotFail(err)

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
	logger.Fatal("app: http server closed",
		zap.Error(http.ListenAndServe(config.Server.Host(), router)),
	)
}
