package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"

	pb "github.com/GophKeeper/server/cmd/proto"

	"github.com/GophKeeper/client/cmd/config"
	"github.com/GophKeeper/client/internal/controller/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	version   = "1.0.0"
	buildDate = "2024-06-21"
)

func main() {

	cfg := config.NewConfig()
	parseFlags(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(cfg.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cUser := pb.NewUserServiceClient(conn)
	cSync := pb.NewSyncServiceClient(conn)
	cManegment := pb.NewManagementServiceClient(conn)

	fmt.Println("GophKeeper запущен")

	cli := cli.NewCLI()

	for {
		fmt.Print("\nВведите команду: ")
		reader := bufio.NewReader(os.Stdin)
		commandstr, _ := reader.ReadString('\n')
		commandstr = strings.TrimSpace(commandstr)
		commandstr = strings.ToLower(commandstr)

		words := strings.Split(commandstr, " ")

		switch words[0] {
		case "version":
			fmt.Printf("Version: %s\nBuild Date: %s\nGo Version: %s", version, buildDate, runtime.Version())
		case "register":
			cli.Register(ctx, cUser, words)
		case "login":
			ctx = cli.Login(ctx, cUser, cSync, words)
		case "addlogindata":
			if cli.CheckLogin() {
				cli.AddLoginData(ctx, cManegment, words)
			}
		case "editlogindata":
			if cli.CheckLogin() {
				cli.EditLoginData(ctx, cManegment, words)
			}
		case "deletelogindata":
			if cli.CheckLogin() {
				cli.DeleteLoginData(ctx, cManegment, words)
			}
		case "listlogindata":
			cli.ListLoginData(ctx, cSync)
		case "exit":
			fmt.Println("Программа завершена.")
			os.Exit(0)
		case "help":
			cli.Help()
		default:
			fmt.Println("Неверная команда. Попробуйте снова.")
		}
	}

}
