package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/client/cmd/mycripto"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
		Login:    mycripto.Encrypt([]byte(words[1]), []byte(cli.Key)),
		Password: mycripto.Encrypt([]byte(words[2]), []byte(cli.Key)),
		Meta:     mycripto.Encrypt([]byte(meta), []byte(cli.Key)),
	}

	newLoginData, err := c.AddLoginData(ctx, &loginData)
	if err != nil {
		checkError(err)
		return
	}

	fmt.Println("SUCCESS. Данные успешно добавлены")

	//сохранение данных на клиенте, если они успешно сохранены на сервере
	cli.MyData.ListLoginData[newLoginData.Uid] = newLoginData
}

func (cli *CLI) ListLoginData(ctx context.Context, c pb.SyncServiceClient) {

	respList, err := c.ListLoginData(ctx, &emptypb.Empty{})
	if err != nil {
		checkError(err)
		fmt.Println("\x1b[1;31m\nOFFLINE режим\x1b[0m")
	} else {
		for _, ld := range respList.LoginData {
			cli.MyData.ListLoginData[ld.Uid] = ld
		}
	}

	fmt.Println("№ | Логин | Пароль | Метаинформация | ID | Дата создания")
	n := 1
	for i, d := range cli.MyData.ListLoginData {

		meta := ""
		if len([]byte(d.Meta)) > 0 {
			meta = mycripto.Decrypt([]byte(d.Meta), []byte(cli.Key))
		}

		fmt.Println(n, "|",
			mycripto.Decrypt([]byte(d.Login), []byte(cli.Key)), "|",
			mycripto.Decrypt([]byte(d.Password), []byte(cli.Key)), "|",
			meta, "|",
			i, "|", d.Created.AsTime().Format("2006-01-02 15:04:05"))
		n++
	}
}

func (cli *CLI) EditLoginData(ctx context.Context, c pb.ManagementServiceClient, words []string) {

	reader := bufio.NewReader(os.Stdin)

	if len(words) < 3 {
		if len(words) < 2 {
			fmt.Print("Введите значение id: ")

			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)
			words = append(words, id)
		}

		fmt.Print("Введите значение newpassword: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		words = append(words, password)
	}

	loginData := pb.LoginData{
		Uid:      words[1],
		Password: mycripto.Encrypt([]byte(words[2]), []byte(cli.Key)),
	}

	_, err := c.EditLoginData(ctx, &loginData)
	if err != nil {
		checkError(err)
		return
	}

	fmt.Println("SUCCESS. Данные успешно изменены")

	//изменнение данных на клиенте, если они успешно изменены на сервере
	cli.MyData.ListLoginData[loginData.Uid].Password = loginData.Password
}

func (cli *CLI) DeleteLoginData(ctx context.Context, c pb.ManagementServiceClient, words []string) {

	reader := bufio.NewReader(os.Stdin)

	if len(words) < 2 {
		fmt.Print("Введите значение id: ")

		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		words = append(words, id)
	}

	loginData := pb.LoginData{
		Uid: words[1],
	}

	_, err := c.DeleteLoginData(ctx, &loginData)
	if err != nil {
		checkError(err)
		return
	}

	fmt.Println("SUCCESS. Данные успешно удалены")

	//удаление данных на клиенте, если они успешно удалены на сервере
	delete(cli.MyData.ListLoginData, loginData.Uid)
}
