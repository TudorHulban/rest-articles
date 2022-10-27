package db

import (
	"database/sql"
	"sync"
)

var _dbConnection *sql.DB
var mu sync.RWMutex

func GetDBConnection() (*sql.DB, error) {
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
