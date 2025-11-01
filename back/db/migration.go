package db

import model "app/db/models"

func startMigration() (err error) {
	models := []any{
		&model.Account{},
		&model.Event{},
		&model.AccountEvent{},
		&model.AccountProvider{},
	}

	for _, m := range models {
		if err := conn.AutoMigrate(m); err != nil {
			return err
		}
	}
	return nil
}
