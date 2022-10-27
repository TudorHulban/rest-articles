package db

import (
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type configDB struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
}

func newTestConfigDB() configDB {
	return configDB{
		Name:     "rest",
		User:     "postgres",
		Password: "thepassword",
		Host:     "localhost",
		Port:     5432,
	}
}

func connect(cfg configDB) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
