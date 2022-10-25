package db

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

var _dbConnection *sqlx.DB
var mu sync.RWMutex

func GetDBConnection() (*sqlx.DB, error) {
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
