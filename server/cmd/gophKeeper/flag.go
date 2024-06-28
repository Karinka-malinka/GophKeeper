package main

import (
	"flag"
	"os"

	"github.com/GophKeeper/server/cmd/config"
)

func parseFlags(cfg *config.ConfigData) {

	flag.StringVar(&cfg.RunAddrgRPS, "a", "localhost:8888", "address and port to run server gRPS")
	flag.StringVar(&cfg.RunAddrREST, "r", "localhost:8080", "address and port to run server RESR")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "adress connect database")
	flag.StringVar(&cfg.PathKeys, "k", "", "path keys for token")

	flag.Parse()

	if envRunAddrgRPS := os.Getenv("RUN_ADDRESS_GRPS"); envRunAddrgRPS != "" {
		cfg.RunAddrgRPS = envRunAddrgRPS
	}

	if envRunAddrRESR := os.Getenv("RUN_ADDRESS_REST"); envRunAddrRESR != "" {
		cfg.RunAddrREST = envRunAddrRESR
	}

	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		cfg.DatabaseDSN = envDatabaseDSN
	}

	if envPathKeys := os.Getenv("PATH_KEYS"); envPathKeys != "" {
		cfg.PathKeys = envPathKeys
	}
}
