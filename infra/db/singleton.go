package db

import (
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

	var errCo error

	_dbConnection, errCo = connect(newTestConfigDB())
	if errCo != nil {
		return nil, errCo
	}

	return _dbConnection, nil
}

// Close releases the database connection.
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	sqlDB, errDB := _dbConnection.DB()
	if errDB != nil {
		return errDB
	}

	return sqlDB.Close()
}
