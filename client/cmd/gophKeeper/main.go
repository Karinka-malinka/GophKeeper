package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	pb "github.com/GophKeeper/server/cmd/proto"

	"github.com/GophKeeper/client/cmd/config"
	"github.com/GophKeeper/client/internal/controller/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	cService := pb.NewUserServiceClient(conn)
	//cSync := pb.NewSyncServiceClient(conn)
	cManegment := pb.NewManagementServiceClient(conn)

	fmt.Println("Соединение с сервером установлено")

	cli := cli.NewCLI()

	for {
		fmt.Print("\nВведите команду: ")
		reader := bufio.NewReader(os.Stdin)
		commandstr, _ := reader.ReadString('\n')
		commandstr = strings.TrimSpace(commandstr)
		commandstr = strings.ToLower(commandstr)

		words := strings.Split(commandstr, " ")

		switch words[0] {
		case "register":
			cli.Register(ctx, cService, words)
		case "login":
			cli.Login(ctx, cService, words)
		case "addlogindata":
			cli.AddLoginData(ctx, cManegment, words)
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
