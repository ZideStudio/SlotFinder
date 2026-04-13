package db

import (
	"app/commons/constants"
	model "app/db/models"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func startMigration() (err error) {
	// Migrate types
	if err := ensureEventStatusEnumType(); err != nil {
		log.Error().Err(err).Msg("failed to ensure postgres enum type event_status exists")
		panic(err)
	}

	// Migrate models
	models := []any{
		&model.Account{},
		&model.Event{},
		&model.Availability{},
		&model.Slot{},
		&model.AccountEvent{},
		&model.AccountProvider{},
		&model.RefreshToken{},
	}

	for _, m := range models {
		if err := conn.AutoMigrate(m); err != nil {
			return err
		}
	}
	return nil
}

func ensureEventStatusEnumType() error {
	quoted := make([]string, 0, len(constants.EventStatuses))
	for _, s := range constants.EventStatuses {
		quoted = append(quoted, fmt.Sprintf("'%s'", strings.ReplaceAll(string(s), "'", "''")))
	}
	enumValues := strings.Join(quoted, ", ")

	sql := fmt.Sprintf(`
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_status') THEN
			CREATE TYPE event_status AS ENUM (%s);
		END IF;
	END
	$$;
	`, enumValues)

	return conn.Exec(sql).Error
}
