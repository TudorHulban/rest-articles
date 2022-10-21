package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ConfigDB struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
}

func NewTestConfigDB() ConfigDB {
	return ConfigDB{
		Name:     "rest",
		User:     "postgres",
		Password: "thepassword",
		Host:     "localhost",
		Port:     5432,
	}
}

func Connect(cfg ConfigDB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	return sqlx.Connect("postgres", dsn)
}
