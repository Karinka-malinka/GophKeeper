package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/client/cmd/mycripto"
)

func (cli *CLI) AddFile(ctx context.Context, c pb.ManagementServiceClient, words []string) {

	var meta string

	reader := bufio.NewReader(os.Stdin)

	if len(words) < 2 {
		fmt.Print("Введите значение path: ")

		login, _ := reader.ReadString('\n')
		login = strings.TrimSpace(login)
		words = append(words, login)
	}

	if len(words) > 2 {
		meta = strings.Join(words[2:], " ")
	}

	// Открытие файла для чтения
	file, err := os.Open(words[1])
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Получение информации о файле для проверки размера
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return
	}

	// Проверка размера файла (не больше 5 Мб)
	maxSize := int64(5 * 1024 * 1024) // 5 Мб в байтах
	if fileInfo.Size() > maxSize {
		fmt.Println("Файл слишком большой. Максимальный размер файла - 5 Мб")
		return
	}

	// Чтение данных файла в виде байтов
	data := make([]byte, fileInfo.Size())
	_, err = io.ReadFull(file, data)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Получение имени файла из пути
	fileName := filepath.Base(words[1])

	filedata := pb.File{
		Name: mycripto.Encrypt([]byte(fileName), []byte(cli.Key)),
		File: mycripto.Encrypt(data, []byte(cli.Key)),
		Meta: mycripto.Encrypt([]byte(meta), []byte(cli.Key)),
	}

	newFile, err := c.AddFile(ctx, &filedata)
	if err != nil {
		checkError(err)
		return
	}

	fmt.Printf("SUCCESS. Данные успешно добавлены (id %s)\n", newFile.Uid)

	//сохранение данных на клиенте, если они успешно сохранены на сервере
	cli.MyData.ListFile[newFile.Uid] = newFile
}

func (cli *CLI) GetFile(ctx context.Context, c pb.ManagementServiceClient, words []string) {

	reader := bufio.NewReader(os.Stdin)

	if len(words) < 2 {
		fmt.Print("Введите значение id: ")

		id, _ := reader.ReadString('\n')
		id = strings.TrimSpace(id)
		words = append(words, id)
	}

	filedata := pb.File{
		Uid: words[1],
	}

	newFile, err := c.GetFile(ctx, &filedata)
	if err != nil {
		checkError(err)
		return
	}

	file, err := os.Create(mycripto.Decrypt(newFile.Name, []byte(cli.Key)))
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(mycripto.Decrypt(newFile.File, []byte(cli.Key))))
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
		return
	}

	fmt.Printf("SUCCESS. Файл успешно сохранен (%s)\n", file.Name())
}
