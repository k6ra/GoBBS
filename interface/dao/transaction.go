package dao

import (
	"database/sql"

	"github.com/pkg/errors"
)

// ExecWithTx DB操作をトランザクションで実行する
func ExecWithTx[T any](db *sql.DB, f func(tx *sql.Tx) (T, error)) (result T, err error) {
	var resultZeroValue T

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			result = resultZeroValue
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Wrapf(err, "rollback error(%s)", rbErr.Error())
			}
		}
	}()

	result, err = f(tx)
	if err != nil {
		return
	}

	err = tx.Commit()

	return
}
