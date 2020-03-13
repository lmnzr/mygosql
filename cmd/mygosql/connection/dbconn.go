package database

import "database/sql"

//DBConn : Interface For sql.DB and sql.Tx
type DBConn interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...interface{}) (sql.Result, error)
}
