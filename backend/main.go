package main

import (
	"backend/internal/config"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		logger.Error().Err(err)
		os.Exit(-1)
	}

	logger.Info().Msg("Connected to PostgreSQL")

	connectionDb, err := db.DB()
	if err != nil {
		logger.Error().Err(err)
	}

	defer func() {
		if err := connectionDb.Close(); err != nil {
			logger.Error().Err(err).Msg("")
		}
	}()

	server := http.Server{
		Addr: "localhost:5000",
	}

	fmt.Println("Listening on localhost:5000")

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
