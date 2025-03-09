package helper

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			logrus.Errorf("rollback error %v", rollbackError)
			panic(rollbackError)
		}
		logrus.Errorf("error recover %v", err)
		panic(err)
	}
	commitRollback := tx.Commit()
	if commitRollback != nil {
		logrus.Errorf("commit error %v", commitRollback)
		panic(commitRollback)
	}
}
