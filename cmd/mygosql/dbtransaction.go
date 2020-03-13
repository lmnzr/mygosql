package database

import (
	"database/sql"
)

//Begin :
func Begin(dbcon *sql.DB) (*sql.Tx,error) {
	return dbcon.Begin()
}

//Commit :
func Commit(tx *sql.Tx) error {
	return tx.Commit()
}

//Rollback : Rollback transaction
func Rollback(tx *sql.Tx) error {
	var txerr error
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		txerr = err
	}

	return txerr
}
