package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/client/cmd/mycripto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CLI struct {
	Token  string
	UserID string
	Key    string
	MyData MyData
}

type MyData struct {
	ListLoginData map[string]*pb.LoginData
	ListText      map[string]*pb.Text
	ListFile      map[string]*pb.File
	ListBankCard  map[string]*pb.BankCard
}

func NewCLI() *CLI {

	listLoginData := make(map[string]*pb.LoginData)
	listText := make(map[string]*pb.Text)
	listFile := make(map[string]*pb.File)
	listBankCard := make(map[string]*pb.BankCard)

	myData := MyData{ListLoginData: listLoginData, ListText: listText, ListFile: listFile, ListBankCard: listBankCard}

	return &CLI{MyData: myData}
}

func (cli *CLI) Help() {
	fmt.Println("\nВозможные команды:")
	fmt.Println("__Доступ__")
	fmt.Println("REGISTER username password - регистрация нового пользователя") //реализовано
	fmt.Println("LOGIN username password - вход существующего пользователя")    //реализовано
	fmt.Println("EXIT - выход из программы")                                    //реализовано

	fmt.Println("\n__Управление парами логин/пароль__")
	fmt.Println("ADDLOGINDATA login password [meta] - добавление пары логин/пароль")               //реализовано
	fmt.Println("LISTLOGINDATA - получить список пар логин/пароль")                                //реализовано
	fmt.Println("EDITLOGINDATA id newpassword - изменить пароль у существующей пары логин/пароль") //реализовано
	fmt.Println("DELETELOGINDATA id - удалить пару логин/пароль")                                  //реализовано

	fmt.Println("\n__Управление текстовыми данными__")
	fmt.Println("ADDTEXT text [meta] - добавление текстовых данных") //будет чуть позже
	fmt.Println("LISTTEXT - получить список текстовых данных")       //будет чуть позже
	fmt.Println("DELETETEXT id - удалить текстовые данные")          //будет чуть позже

	fmt.Println("\n__Управление бинарными данными__")
	fmt.Println("ADDFILE path [meta] - добавление файла")     //реализовано
	fmt.Println("LISTFILE - получить список бинарных данных") //реализовано - возвращает только id, name и meta
	fmt.Println("GETFILE id - получить бинарныe данныe")      //реализовано - получает бинарные данные и сохраняет на диск
	fmt.Println("DELETEFILE id- удалить бинарные данные")     //будет чуть позже

	fmt.Println("\n__Управление данными банковских карт__")
	fmt.Println("ADDBANKCARD number term ccv [meta] - добавление данных о бансковское карте") //будет чуть позже
	fmt.Println("LISTBANKCARD - получить список банковских карт")                             //будет чуть позже
	fmt.Println("DELETEBANKCARD number - удалить данные о банковской карте")                  //будет чуть позже
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

	// Добавление таймаута в 1 секунду к контексту
	ctxWithTimeout, cancelTimeout := context.WithTimeout(ctx, 1*time.Second)
	defer cancelTimeout()

	resp, err := c.Register(ctxWithTimeout, &userReq)
	if err != nil {
		checkError(err)
		return
	}

	cli.Token = resp.Token
	cli.UserID = resp.Uid
	cli.Key = resp.Key

	fmt.Println("SUCCESS. Вы успешно зарегистрированы!")
}

func (cli *CLI) Login(ctx context.Context, c pb.UserServiceClient, m pb.SyncServiceClient, words []string) context.Context {

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

	// Добавление таймаута в 3 секунду к контексту
	//ctx, cancelTimeout := context.WithTimeout(ctx, 3*time.Second)
	//defer cancelTimeout()

	resp, err := c.Login(ctx, &userReq)
	if err != nil {
		checkError(err)
		return ctx
	}

	cli.Token = resp.Token
	cli.UserID = resp.Uid
	cli.Key = resp.Key + mycripto.MirrorUID(resp.Uid)

	fmt.Println("Добро пожаловать в GophKeeper!")

	md := metadata.New(map[string]string{"access_token": cli.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	//первичная синхронизация данных между сервером и клиентом
	//список пар логин/пароль
	respList, err := m.ListLoginData(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("ERROR. Синхронизация данных не выполнена")
		checkError(err)
		return ctx
	}
	for _, ld := range respList.LoginData {
		cli.MyData.ListLoginData[ld.Uid] = ld
	}
	fmt.Println("SUCCESS. Синхронизация данных для пар логин/пароль выполнена")

	//список текстовых данных
	//TODO: здесь будет код синхронизации приватных текстовых данных
	fmt.Println("SUCCESS. Синхронизация приватной текстовой информации выполнена")

	//список файлов
	//TODO: здесь будет код синхронизации приватных файлов
	fmt.Println("SUCCESS. Синхронизация приватных файлов выполнена")

	//список банковских карт
	//TODO: здесь будет код синхронизации банковских данных
	fmt.Println("SUCCESS. Синхронизация данных о банковских картах выполнена")

	return ctx
}

func (cli *CLI) CheckLogin() bool {
	login := true
	if cli.Token == "" {
		fmt.Println("ERROR. Выполните сначала команду LOGIN")
		login = false
	}
	return login
}

func checkError(err error) {

	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Unauthenticated: //неверный логин или пароль или протух токен
			fmt.Println("Неверное имя пользователя или пароль. Выполните команду LOGIN заново")
		case codes.AlreadyExists:
			fmt.Println("Логин занят, придумайте другой и повторите процесс регистрации")
		case codes.DeadlineExceeded, codes.Unavailable:
			fmt.Println("Сервер не отвечает. Попробуйте позже")
		case codes.Aborted:
			fmt.Println("Синхронизация данных не выполнена. Попробуйте позже")
		default:
			fmt.Println("Ошибка на сервере. Обратитесь в техническую поддержку")
		}
	} else {
		fmt.Printf("Не получилось распарсить ошибку %v", err)
	}
}
