package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type configSQLiteDB struct {
	Path string
}

func newConfigSQLiteDBRAM() configSQLiteDB {
	return configSQLiteDB{
		Path: "file::memory:",
	}
}

// see in https://gorm.io/docs/connecting_to_the_database.html#SQLite
func connectDBSQLite(cfg configSQLiteDB) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}
