package db

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var _dbConnection *gorm.DB
var mu sync.RWMutex

func GetDBConnection() (*gorm.DB, error) {
	mu.RLock()
	if _dbConnection != nil {
		defer mu.RUnlock()

		return _dbConnection, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	var errPo error

	_dbConnection, errPo = connectDBPostgres(newConfigPostgresDB())
	if errPo == nil {
		return _dbConnection, nil
	}

	// TODO: enhance switching and signalization
	fmt.Println("postgres not available, switching to sqlite")

	var errLite error

	_dbConnection, errLite = connectDBSQLite(newConfigSQLiteDBRAM())
	if errLite != nil {
		return nil, errLite
	}

	return _dbConnection, nil
}

// Close releases the database connection.
func CloseDBConnection() error {
	mu.Lock()
	defer mu.Unlock()

	sqlDB, errDB := _dbConnection.DB()
	if errDB != nil {
		return errDB
	}

	return sqlDB.Close()
}
