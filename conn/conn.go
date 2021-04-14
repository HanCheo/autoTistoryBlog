package conn

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// DbConn connect postgres db
func DbConn() (*sql.DB, error) {
	var (
		dbHost     string = os.Getenv("DB_HOST_dev")
		dbPort     string = os.Getenv("DB_PORT_dev")
		dbUserName string = os.Getenv("DB_USERNAME_dev")
		dbPassword string = os.Getenv("DB_PASSWORD_dev")
		// dbSchemas  string = os.Getenv("DB_SCHEMAS_dev")
		dbName string = os.Getenv("DB_NAME_dev")
	)

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable sslmode=disable", dbHost, dbPort, dbUserName, dbPassword, dbName)

	return sql.Open("postgres", dbinfo)

}
