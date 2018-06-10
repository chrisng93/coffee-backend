package db

import "database/sql"

func createTransaction(db *sql.DB, cb func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	var success bool
	defer func() {
		if !success {
			tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}

	err = cb(tx)
	if err != nil {
		success = false
	} else {
		success = true
	}

	err = tx.Commit()
	if err != nil {
		success = false
	}

	return err
}
