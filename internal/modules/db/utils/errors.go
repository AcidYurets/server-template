package utils

import (
	"server-template/internal/models/err_const"
	"server-template/internal/modules/db/ent"
)

func WrapError(err error) error {
	// Если запись не найдена
	if ent.IsNotFound(err) {
		return err_const.ErrDatabaseRecordNotFound
	}

	return err
}
