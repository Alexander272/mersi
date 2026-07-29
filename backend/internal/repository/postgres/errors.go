package postgres

import (
	"database/sql"

	"github.com/Alexander272/mersi/backend/internal/models"
)

func checkRowsAffected(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return models.ErrNoRows
	}
	return nil
}
