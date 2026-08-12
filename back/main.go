package main

import (
	"app/config"
	"app/db"
	"fmt"
	"os"

	"app/server"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// configureLogging sets zerolog's global time format and level, switching to
// debug verbosity when ENV=local.
func configureLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("ENV") == "local" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// loadDotEnvIfPresent loads environment variables from path if it exists,
// leaving the environment untouched (and returning nil) otherwise.
func loadDotEnvIfPresent(path string) error {
	if _, statErr := os.Stat(path); statErr != nil {
		return nil
	}
	if err := godotenv.Load(path); err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
		return err
	}
	log.Debug().Msg(".env file loaded successfully")
	return nil
}

// @title SlotFinder API
// @version 1.0b
// @description SlotFinder API Doc
//
// @contact.email  contact@zide.fr
//
// @securityDefinitions.apikey AccessTokenCookie
// @in cookie
// @name access_token
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**
func main() {
	configureLogging()

	if err := loadDotEnvIfPresent(".env"); err != nil {
		panic(err)
	}

	config := config.Init()
	db.Init()

	log.Info().Msg(fmt.Sprintf("Server started on %s:%s", config.Host, config.Port))

	server.Init()
}
