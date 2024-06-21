package main

import (
	"flag"
	"os"

	"github.com/GophKeeper/client/cmd/config"
)

func parseFlags(cfg *config.ConfigData) {

	flag.StringVar(&cfg.ServerAddr, "s", "localhost:8888", "address and port to server")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if envServerAddr := os.Getenv("RUN_ADDRESS"); envServerAddr != "" {
		cfg.ServerAddr = envServerAddr
	}

}
