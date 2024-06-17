package main

import (
	"flag"
	"os"

	"github.com/GophKeeper/server/cmd/config"
)

func parseFlags(cfg *config.ConfigData) {

	flag.StringVar(&cfg.RunAddr, "a", "localhost:8888", "address and port to run server")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "adress connect database")

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		cfg.RunAddr = envRunAddr
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		cfg.DatabaseDSN = envDatabaseDSN
	}

}
