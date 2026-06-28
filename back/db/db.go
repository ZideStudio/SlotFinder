package db

import (
	"app/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var conn *gorm.DB

func Init() *gorm.DB {
	var err error
	c := config.GetConfig()

	conn, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  c.Db.GetPostgresConnectionInfo(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		panic(err)
	}

	err = startMigration()
	if err != nil {
		log.Error().Err(err).Msg("failed to migrate database")
		panic(err)
	}

	return conn
}

func TestConnection() bool {
	if conn == nil {
		return false
	}

	db, err := conn.DB()
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return false
	}

	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("failed to ping database")
		return false
	}

	err = db.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close database connection")
		return false
	}

	return true
}

func GetDB() *gorm.DB {
	return conn
}
