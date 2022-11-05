package configs

import "os"

func GetDatabaseName() string {
	res := os.Getenv("POSTGRES_DATABASE")
	if len(res) == 0 {
		return "rest"
	}

	return res
}

func GetDatabaseHost() string {
	if len(os.Getenv("POSTGRES_DATABASE")) == 0 {
		return "localhost"
	}

	return "database"
}
