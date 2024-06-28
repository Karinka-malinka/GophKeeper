package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/cmd/config"
	bankcardApp "github.com/GophKeeper/server/internal/app/bankcard"
	filedataApp "github.com/GophKeeper/server/internal/app/filedata"
	logindataApp "github.com/GophKeeper/server/internal/app/logindata"
	textdataApp "github.com/GophKeeper/server/internal/app/textdata"
	userApp "github.com/GophKeeper/server/internal/app/user"
	"github.com/GophKeeper/server/internal/controller/handler/management"
	"github.com/GophKeeper/server/internal/controller/handler/mysync"
	userHandler "github.com/GophKeeper/server/internal/controller/handler/user"
	"github.com/GophKeeper/server/internal/controller/server"
	"github.com/GophKeeper/server/internal/database"
	bankcardStore "github.com/GophKeeper/server/internal/database/bankcard"
	filedataStore "github.com/GophKeeper/server/internal/database/filedata"
	logindataStore "github.com/GophKeeper/server/internal/database/logindata"
	textdataStore "github.com/GophKeeper/server/internal/database/textdata"
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

	loginDataStore := logindataStore.NewLoginDataStore(db)
	loginDataApp := logindataApp.NewLoginData(loginDataStore)

	textDataStore := textdataStore.NewTextDataStore(db)
	textDataApp := textdataApp.NewTextData(textDataStore)

	fileDataStore := filedataStore.NewFileDataStore(db)
	fileDataApp := filedataApp.NewFiletData(fileDataStore)

	bankCardStore := bankcardStore.NewBankCardDataStore(db)
	bankCardApp := bankcardApp.NewBankCardData(bankCardStore)

	syncHandler := mysync.NewSyncHandler(loginDataApp, textDataApp, fileDataApp, bankCardApp)
	managementHandler := management.NewManagementHandler(loginDataApp, textDataApp, fileDataApp, bankCardApp)

	appServer, err := server.NewServer(ctx, cfg.RunAddrgRPS, cfg.RunAddrREST, userApp)
	if err != nil {
		os.Exit(1)
	}

	pb.RegisterUserServiceServer(appServer.Srv, userHandler)
	pb.RegisterSyncServiceServer(appServer.Srv, syncHandler)
	pb.RegisterManagementServiceServer(appServer.Srv, managementHandler)

	go appServer.Start(ctx)
	go appServer.StartRest(ctx)

	<-ctx.Done()
	appServer.Stop()
	appServer.StopREST(ctx)
	m.Close()
	db.Close()
}
