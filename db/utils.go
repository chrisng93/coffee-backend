package db

import "database/sql"

func createTransaction(db *sql.DB, cb func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	var success bool
	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}

	return cb(tx)
}
