package service

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB     DB
	Server Server
}

type Server struct {
	Addr string `default:"0.0.0.0"`
	Port uint16 `default:"8080"`
}

func (srv Server) Host() string {
	return fmt.Sprintf("%s:%d", srv.Addr, srv.Port)
}

type DB struct {
	User     string `default:"app_user"`
	Password string `default:"app_password"`
	Addr     string `default:"postgres"`
	Port     uint16 `default:"5432"`
	Name     string `default:"app_db"`
}

func (db DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User, db.Password,
		db.Addr, db.Port,
		db.Name,
	)
}

func ParseConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
