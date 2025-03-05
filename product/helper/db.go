package helper

import "database/sql"

func CommitOrRollback(tx *sql.Tx) error {
	err := recover()
	if err != nil {
		rollbackError := tx.Rollback()
		return rollbackError
	}
	commitRollback := tx.Commit()
	return commitRollback
}
