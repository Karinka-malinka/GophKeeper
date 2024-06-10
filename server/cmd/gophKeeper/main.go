package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	pb "github.com/GophKeeper/server/cmd/proto"

	"github.com/GophKeeper/server/cmd/config"
	userApp "github.com/GophKeeper/server/internal/app/user"
	userHandler "github.com/GophKeeper/server/internal/controller/handler/user"
	"github.com/GophKeeper/server/internal/controller/server"
	"github.com/GophKeeper/server/internal/database"
	userStore "github.com/GophKeeper/server/internal/database/user"
	"github.com/GophKeeper/server/internal/logger"
)

func main() {

	logger.SetUpLogger()
	cfg := config.NewConfig()
	parseFlags(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := database.NewDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("error in open database. error: %v", err)
	}
	defer db.Close()

	userStore := userStore.NewUserStore(db)
	userApp := userApp.NewUser(userStore)
	userHandler := userHandler.NewUserHandler(userApp)

	appServer := server.NewServer(cfg.RunAddr)

	pb.RegisterUserServiceServer(appServer.Srv, userHandler)
	go appServer.Start()

	<-ctx.Done()
	appServer.Stop()
}
