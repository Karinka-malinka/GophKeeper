package main

import (
	"context"
	"log/slog"
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

	cfgToken, err := cfg.GetConfig(cfg.PathKeys)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, m, err := database.NewDB(cfg.DatabaseDSN)
	if err != nil {
		os.Exit(1)
	}

	userStore := userStore.NewUserStore(db)
	userApp := userApp.NewUser(userStore, cfgToken)
	userHandler := userHandler.NewUserHandler(userApp)

	appServer, err := server.NewServer(cfg.RunAddr, userApp)
	if err != nil {
		os.Exit(1)
	}

	pb.RegisterUserServiceServer(appServer.Srv, userHandler)

	//pb.RegisterSyncServiceServer(appServer.Srv, syncHandler)

	//pb.RegisterManagementServiceServer(appServer.Srv, managementHandler)

	go appServer.Start(ctx)

	<-ctx.Done()
	appServer.Stop()
	m.Close()
	db.Close()
}
