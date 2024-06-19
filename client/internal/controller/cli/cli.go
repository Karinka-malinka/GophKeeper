package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GophKeeper/client/cmd/mycripto"
	pb "github.com/GophKeeper/server/cmd/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CLI struct {
	Token string
}

func NewCLI() *CLI {
	return &CLI{}
}

func (cli *CLI) Help() {
	fmt.Println("\nВозможные команды:")
	fmt.Println("__Доступ__")
	fmt.Println("REGISTER username password - регистрация нового пользователя") //+
	fmt.Println("LOGIN username password - вход существующего пользователя")    //+
	fmt.Println("EXIT - выход из программы")                                    //+

	//добавление
	fmt.Println("\n__Управление парами логин/пароль__")
	fmt.Println("ADDLOGINDATA login password [meta] - добавление пары логин/пароль")
	fmt.Println("LISTLOGINDATA - получить список пар логин/пароль")
	fmt.Println("EDITLOGINDATA id newpassword - изменить пароль у существующей пары логин/пароль")
	fmt.Println("DELETELOGINDATA id - удалить пару логин/пароль")

	fmt.Println("\n__Управление текстовыми данными__")
	fmt.Println("ADDTEXT text [meta] - добавление текстовых данных")
	fmt.Println("LISTTEXT - получить список текстовых данных")
	fmt.Println("DELETETEXT id - удалить текстовые данные")

	fmt.Println("\n__Управление бинарными данными__")
	fmt.Println("ADDFILE path [meta] - добавление файла")
	fmt.Println("LISTFILE - получить список бинарных данных")
	fmt.Println("DELETEFILE id- удалить бинарные данные")

	fmt.Println("\n__Управление данными банковских карт__")
	fmt.Println("ADDBANKCARD number term ccv [meta] - добавление данных о бансковское карте")
	fmt.Println("LISTBANKCARD - получить список банковских карт")
	fmt.Println("DELETEBANKCARD number - удалить данные о банковской карте")
}

func (cli *CLI) Register(ctx context.Context, c pb.UserServiceClient, words []string) {

	reader := bufio.NewReader(os.Stdin)

	if len(words) < 3 {

		if len(words) < 2 {
			fmt.Print("Введите имя пользователя: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)
			words = append(words, username)
		}

		fmt.Print("Придумайте пароль: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		fmt.Print("Повторите пароль: ")
		password2, _ := reader.ReadString('\n')
		password2 = strings.TrimSpace(password2)

		if password != password2 {
			fmt.Println("Введенные пароли не совпадают. Повторите процесс регистрации еще раз")
			return
		}
		words = append(words, password)
	}

	userReq := pb.UserRequest{Login: words[1], Password: mycripto.WriteHash(words[2])}

	// Добавление таймаута в 200 миллисекунд к контексту
	ctxWithTimeout, cancelTimeout := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancelTimeout()

	resp, err := c.Register(ctxWithTimeout, &userReq)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				fmt.Println("Логин занят, придумайте другой и повторите процесс регистрации")
			case codes.DeadlineExceeded:
				fmt.Println("Сервер не отвечает. Попробуйте позже")
			default:
				fmt.Println(e.Code(), "Ошибка на сервере. Обратитесь в техническую поддержку")
			}
		} else {
			fmt.Printf("Не получилось распарсить ошибку %v", err)
		}
		return
	}

	cli.Token = resp.Token
	fmt.Println("Вы успешно зарегистрированы!")
}

func (cli *CLI) Login(ctx context.Context, c pb.UserServiceClient, words []string) {

	reader := bufio.NewReader(os.Stdin)

	if len(words) < 3 {
		if len(words) < 2 {
			fmt.Print("Введите имя пользователя: ")

			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)
			words = append(words, username)
		}

		fmt.Print("Введите пароль: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		words = append(words, password)
	}

	userReq := pb.UserRequest{Login: words[1], Password: mycripto.WriteHash(words[2])}

	resp, err := c.Login(ctx, &userReq)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				fmt.Println("Неверное имя пользователя или пароль.")
			case codes.DeadlineExceeded:
				fmt.Println("Сервер не отвечает. Попробуйте позже")
			default:
				fmt.Println(e.Code(), e.Message())
			}
		} else {
			fmt.Printf("Не получилось распарсить ошибку %v", err)
		}
		return
	}

	cli.Token = resp.Token
	fmt.Println("Добро пожаловать в GophKeeper!")

	//sync
	//resp, err := c.(ctx, &userReq)
}

func (cli *CLI) AddLoginData(ctx context.Context, c pb.ManagementServiceClient, words []string) {

	var meta string

	reader := bufio.NewReader(os.Stdin)

	if len(words) < 3 {
		if len(words) < 2 {
			fmt.Print("Введите значение login: ")

			login, _ := reader.ReadString('\n')
			login = strings.TrimSpace(login)
			words = append(words, login)
		}

		fmt.Print("Введите значение password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		words = append(words, password)
	}

	if len(words) > 3 {
		meta = strings.Join(words[3:], " ")
	}

	loginData := pb.LoginData{
		Login:    mycripto.Encrypt([]byte(words[1]), []byte("keykeykeykeykey1")),
		Password: mycripto.Encrypt([]byte(words[2]), []byte("keykeykeykeykey1")),
		Meta:     mycripto.Encrypt([]byte(meta), []byte("keykeykeykeykey1")),
	}

	_, err := c.AddLoginData(ctx, &loginData)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				fmt.Println("Неверное имя пользователя или пароль.")
			case codes.DeadlineExceeded:
				fmt.Println("Сервер не отвечает. Попробуйте позже")
			default:
				fmt.Println(e.Code(), e.Message())
			}
		} else {
			fmt.Printf("Не получилось распарсить ошибку %v", err)
		}
		return
	}

	fmt.Println("Данные успешно добавлены")
	//sync

}
