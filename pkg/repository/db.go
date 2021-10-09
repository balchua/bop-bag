package repository

import (
	"context"
	"database/sql"
)

type Db interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	GetClusterInfo() ([]byte, error)
	RemoveNode(address string) (string, error)
	Leader() (string, error)
}
