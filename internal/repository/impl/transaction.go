package impl

import "gorm.io/gorm"

func Transaction(db *gorm.DB, fns ...func(*gorm.DB) error) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, fn := range fns {
		if err := fn(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
