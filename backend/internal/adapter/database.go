package adapter

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConnection struct {
	Client *sql.DB
}

func NewDBConnection(username, password, host, dbName string, port int) (*DBConnection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	return &DBConnection{
		Client: db,
	}, nil
}
