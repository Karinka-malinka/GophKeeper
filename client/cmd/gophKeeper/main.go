package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	pb "github.com/GophKeeper/server/cmd/proto"

	"github.com/GophKeeper/client/cmd/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Token string

func main() {

	cfg := config.NewConfig()
	parseFlags(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Добавление таймаута в 200 миллисекунд к контексту
	ctxWithTimeout, cancelTimeout := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancelTimeout()

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(cfg.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	fmt.Println("Соединение с сервером установлено")

	for {
		fmt.Print("Введите команду: ")
		reader := bufio.NewReader(os.Stdin)
		commandstr, _ := reader.ReadString('\n')
		commandstr = strings.TrimSpace(commandstr)
		commandstr = strings.ToLower(commandstr)

		words := strings.Split(commandstr, " ")

		switch words[0] {
		case "register":
			register(ctxWithTimeout, c, words)
		case "login":
			login(ctxWithTimeout, c, words)
		case "exit":
			fmt.Println("Программа завершена.")
			os.Exit(0)
		case "help":
			help()
		default:
			fmt.Println("Неверная команда. Попробуйте снова.")
		}
	}

}
