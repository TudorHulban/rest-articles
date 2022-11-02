package db

import (
	"fmt"

	"github.com/TudorHulban/rest-articles/app/configs"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		Name:     configs.GetDatabaseName(),
		User:     "postgres",
		Password: "thepassword",
		Host:     "database",
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

	return gorm.Open(postgres.Open(dsn), &gorm.Config{}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}
