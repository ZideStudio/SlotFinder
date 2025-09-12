package db

import model "app/db/models"

func startMigration() {
	conn.AutoMigrate(&model.Account{})
	conn.AutoMigrate(&model.AccountProvider{})
}
