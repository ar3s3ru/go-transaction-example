package main

import (
	"go-transaction-example/internal/pkg/must"
	"go-transaction-example/internal/platform/service"
)

func main() {
	config, err := service.ParseConfig()
	must.NotFail(err)

	must.NotFail(service.Run(config))
}
