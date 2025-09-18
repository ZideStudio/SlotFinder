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

// @title SlotFinder API
// @version 1.0b
// @description SlotFinder API Doc
//
// @contact.email  contact@zide.fr
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("ENV") == "local" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Error loading .env file")
		panic(err)
	}

	config := config.Init()
	db.Init()

	log.Info().Msg(fmt.Sprintf("Server started on %s:%s", config.Host, config.Port))

	server.Init()
}
