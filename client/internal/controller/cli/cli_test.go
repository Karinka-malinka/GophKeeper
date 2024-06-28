package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GophKeeper/api/proto"
	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/client/cmd/mycripto"
	"github.com/GophKeeper/client/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Mocks struct {
	ManagementServiceClient *mocks.MockManagementServiceClient
	SyncServiceClient       *mocks.MockSyncServiceClient
	UserServiceClient       *mocks.MockUserServiceClient
}

type Request struct {
	loginData *pb.LoginData
	textData  *pb.Text
	fileData  *pb.File
	bankCard  *pb.BankCard
}

func getMosk(t *testing.T) *Mocks {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок объект
	mocksClient := Mocks{}

	mocksClient.ManagementServiceClient = mocks.NewMockManagementServiceClient(ctrl)
	mocksClient.SyncServiceClient = mocks.NewMockSyncServiceClient(ctrl)
	mocksClient.UserServiceClient = mocks.NewMockUserServiceClient(ctrl)

	return &mocksClient
}

func TestHelp(t *testing.T) {
	cli := &CLI{}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cli.Help()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	expectedOutput := `
Возможные команды:
__Доступ__
REGISTER username password - регистрация нового пользователя
LOGIN username password - вход существующего пользователя
EXIT - выход из программы

__Управление парами логин/пароль__
ADDLOGINDATA login password [meta] - добавление пары логин/пароль
LISTLOGINDATA - получить список пар логин/пароль
EDITLOGINDATA id newpassword - изменить пароль у существующей пары логин/пароль
DELETELOGINDATA id - удалить пару логин/пароль

__Управление текстовыми данными__
ADDTEXT text [meta] - добавление текстовых данных
LISTTEXT - получить список текстовых данных
DELETETEXT id - удалить текстовые данные

__Управление бинарными данными__
ADDFILE path [meta] - добавление файла
LISTFILE - получить список бинарных данных
GETFILE id - получить бинарныe данныe
DELETEFILE id- удалить бинарные данные

__Управление данными банковских карт__
ADDBANKCARD number term ccv [meta] - добавление данных о бансковское карте
LISTBANKCARD - получить список банковских карт
DELETEBANKCARD number - удалить данные о банковской карте
`

	if buf.String() != expectedOutput {
		t.Errorf("Expected output does not match actual output")
	}
}

func TestCLI_AddLoginData(t *testing.T) {

	ctx := context.Background()

	mockClient := getMosk(t).ManagementServiceClient

	cli := NewCLI()
	cli.UserID = "2d618858-27bb-474c-9150-8a03126feff7"
	cli.Key = "a2F4c1e8D9d3A7fe"

	words := []string{"addlogindata", "karina", "123", "meta meta meta"}

	uid := uuid.New().String()
	request := Request{}
	request.loginData = &proto.LoginData{
		Uid:      uid,
		Login:    []byte(words[1]),
		Password: []byte(words[2]),
		Meta:     []byte(words[3]),
	}

	t.Run("Input", func(t *testing.T) {

		currentLen := len(cli.MyData.ListLoginData)

		//mockClient := &MockManagementServiceClient{}
		//mockClient.On("AddLoginData", ctx, mock.Anything).Return(request.loginData, nil)

		mockClient.EXPECT().AddLoginData(ctx, gomock.Any()).Return(request.loginData, nil).Times(1)

		// Макетирование ввода с клавиатуры
		mockInput := "testLogin\ntestPassword\n"
		r, w, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = oldStdin
		}()

		w.WriteString(mockInput)
		w.Close()

		words := []string{"add"}

		cli.AddLoginData(ctx, mockClient, words)

		// Проверка, что login и password были добавлены в words
		if len(cli.MyData.ListLoginData) != currentLen+1 {
			t.Error("Incorrect add")
		}

		/*
			expectedLogin := "testLogin"
			if words[1] != expectedLogin {
				t.Errorf("Expected login to be %s, got %s", expectedLogin, words[1])
			}

			expectedPassword := "testPassword"
			if words[2] != expectedPassword {
				t.Errorf("Expected password to be %s, got %s", expectedPassword, words[2])
			}
		*/
	})

	t.Run("Success", func(t *testing.T) {

		mockClient.EXPECT().AddLoginData(ctx, gomock.Any()).Return(request.loginData, nil).Times(1)

		cli.AddLoginData(ctx, mockClient, words)

		// Check if the login data was added to MyData
		if _, ok := cli.MyData.ListLoginData[uid]; !ok {
			t.Error("Login data was not added to MyData")
		}
	})

	t.Run("Unauthenticated", func(t *testing.T) {

		mockClient.EXPECT().AddLoginData(ctx, gomock.Any()).Return(request.loginData, status.Errorf(codes.Unauthenticated, "")).Times(1)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.AddLoginData(ctx, mockClient, words)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		expected := "Неверное имя пользователя или пароль. Выполните команду LOGIN заново\n"

		if buf.String() != expected {
			t.Errorf("Expected output: %s, got: %s", expected, buf.String())
		}
	})

	t.Run("Internal", func(t *testing.T) {

		mockClient.EXPECT().AddLoginData(ctx, gomock.Any()).Return(request.loginData, status.Errorf(codes.Internal, "")).Times(1)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.AddLoginData(ctx, mockClient, words)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		expected := "Ошибка на сервере. Обратитесь в техническую поддержку\n"

		if buf.String() != expected {
			t.Errorf("Expected output: %s, got: %s", expected, buf.String())
		}
	})
}

func TestCLI_ListLoginData(t *testing.T) {

	ctx := context.Background()

	mockClient := getMosk(t).SyncServiceClient

	cli := NewCLI()
	cli.UserID = "2d618858-27bb-474c-9150-8a03126feff7"
	cli.Key = "a2F4c1e8D9d3A7fe"

	resp := pb.LoginDataResponse{}
	resp.LoginData = append(resp.LoginData, &pb.LoginData{
		Uid:      uuid.New().String(),
		Login:    []byte(mycripto.Encrypt([]byte("login1"), []byte(cli.Key))),
		Password: []byte(mycripto.Encrypt([]byte("pass1"), []byte(cli.Key))),
		Meta:     []byte(mycripto.Encrypt([]byte("meta 1"), []byte(cli.Key))),
		Created:  timestamppb.New(time.Now().UTC()),
	})
	resp.LoginData = append(resp.LoginData, &pb.LoginData{
		Uid:      uuid.New().String(),
		Login:    []byte(mycripto.Encrypt([]byte("login2"), []byte(cli.Key))),
		Password: []byte(mycripto.Encrypt([]byte("pass2"), []byte(cli.Key))),
		Meta:     []byte(mycripto.Encrypt([]byte("meta 2"), []byte(cli.Key))),
		Created:  timestamppb.New(time.Now().UTC()),
	})

	t.Run("Success", func(t *testing.T) {

		mockClient.EXPECT().ListLoginData(ctx, gomock.Any()).Return(&resp, nil).Times(1)

		cli.ListLoginData(ctx, mockClient)

		if len(cli.MyData.ListLoginData) != len(resp.LoginData) {
			t.Error("Error sync login data")
		}
	})

	t.Run("Internal", func(t *testing.T) {

		mockClient.EXPECT().ListLoginData(ctx, gomock.Any()).Return(&resp, status.Errorf(codes.Internal, "")).Times(1)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.ListLoginData(ctx, mockClient)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		expected := "OFFLINE режим"

		if !strings.Contains(buf.String(), expected) {
			t.Errorf("Expected output: %s, got: %s", expected, buf.String())
		}
	})
}

func TestCLI_DeleteLoginData(t *testing.T) {

	ctx := context.Background()

	mockClient := getMosk(t).ManagementServiceClient

	uid := uuid.New().String()

	cli := NewCLI()
	cli.UserID = "2d618858-27bb-474c-9150-8a03126feff7"
	cli.Key = "a2F4c1e8D9d3A7fe"

	cli.MyData.ListLoginData[uid] = &pb.LoginData{
		Uid:      uid,
		Login:    []byte("testlogin"),
		Password: []byte("pass"),
	}

	words := []string{"delete", uuid.New().String()}

	request := Request{}
	request.loginData = &proto.LoginData{
		Uid: words[1],
	}

	t.Run("Input", func(t *testing.T) {

		currentLen := len(cli.MyData.ListLoginData)

		//mockClient.On("DeleteLoginData", ctx, mock.Anything).Return(&empty.Empty{}, nil)
		mockClient.EXPECT().DeleteLoginData(ctx, gomock.Any()).Return(&empty.Empty{}, nil).Times(1)

		// Макетирование ввода с клавиатуры
		mockInput := uid + "\n"
		r, w, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = oldStdin
		}()

		w.WriteString(mockInput)
		w.Close()

		words := []string{"add"}

		cli.DeleteLoginData(ctx, mockClient, words)

		// Проверка, что login и password были добавлены в words
		if len(cli.MyData.ListLoginData) != currentLen-1 {
			t.Error("Incorrect delete")
		}
	})

	t.Run("Success", func(t *testing.T) {

		//		mockClient.On("DeleteLoginData", ctx, mock.Anything).Return(&empty.Empty{}, nil)
		mockClient.EXPECT().DeleteLoginData(ctx, gomock.Any()).Return(&empty.Empty{}, nil).Times(1)

		cli.DeleteLoginData(ctx, mockClient, words)

		// Check if the login data was deleted to MyData
		if _, ok := cli.MyData.ListLoginData[uid]; ok {
			t.Error("Login data was not deleted to MyData")
		}
	})

	t.Run("Internal", func(t *testing.T) {

		//mockClient.On("DeleteLoginData", ctx, mock.Anything).Return(&empty.Empty{}, status.Errorf(codes.Internal, ""))
		mockClient.EXPECT().DeleteLoginData(ctx, gomock.Any()).Return(&empty.Empty{}, status.Errorf(codes.Internal, "")).Times(1)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.DeleteLoginData(ctx, mockClient, words)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		expected := "Ошибка на сервере. Обратитесь в техническую поддержку\n"

		if buf.String() != expected {
			t.Errorf("Expected output: %s, got: %s", expected, buf.String())
		}
	})
}

func TestCLI_UpdateLoginData(t *testing.T) {

	ctx := context.Background()

	mockClient := getMosk(t).ManagementServiceClient

	cli := NewCLI()
	cli.UserID = "2d618858-27bb-474c-9150-8a03126feff7"
	cli.Key = "a2F4c1e8D9d3A7fe"

	uid := uuid.New().String()
	words := []string{"edit", uid, "newpass"}

	cli.MyData.ListLoginData[uid] = &pb.LoginData{
		Uid:      uid,
		Login:    []byte("testlogin"),
		Password: []byte("pass"),
	}

	request := Request{}
	request.loginData = &proto.LoginData{
		Uid:      words[1],
		Password: []byte(mycripto.Encrypt([]byte(words[2]), []byte(cli.Key))),
	}

	t.Run("Input", func(t *testing.T) {

		//mockClient.On("EditLoginData", ctx, mock.Anything).Return(&empty.Empty{}, nil)
		mockClient.EXPECT().EditLoginData(ctx, gomock.Any()).Return(&empty.Empty{}, nil).Times(1)

		// Макетирование ввода с клавиатуры
		mockInput := uid + "\nnewpass\n"
		r, w, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = oldStdin
		}()

		w.WriteString(mockInput)
		w.Close()

		words := []string{"edit"}

		cli.EditLoginData(ctx, mockClient, words)

		// Проверка, что login и password были добавлены в words
		if mycripto.Decrypt(cli.MyData.ListLoginData[uid].Password, []byte(cli.Key)) != "newpass" {
			t.Error("Incorrect edit")
		}
	})

	t.Run("Success", func(t *testing.T) {

		//mockClient.On("EditLoginData", ctx, mock.Anything).Return(&empty.Empty{}, nil)
		mockClient.EXPECT().EditLoginData(ctx, gomock.Any()).Return(&empty.Empty{}, nil).Times(1)

		cli.EditLoginData(ctx, mockClient, words)

		// Check if the login data was deleted to MyData
		if _, ok := cli.MyData.ListLoginData[uid]; !ok {
			t.Error("Login data was not deleted to MyData")
		}
		if mycripto.Decrypt(cli.MyData.ListLoginData[uid].Password, []byte(cli.Key)) != words[2] {
			t.Error("Password not update")
		}
	})

	t.Run("Internal", func(t *testing.T) {

		//mockClient.On("EditLoginData", ctx, mock.Anything).Return(&empty.Empty{}, status.Errorf(codes.Internal, ""))
		mockClient.EXPECT().EditLoginData(ctx, gomock.Any()).Return(&empty.Empty{}, status.Errorf(codes.Internal, "")).Times(1)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		cli.EditLoginData(ctx, mockClient, words)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		expected := "Ошибка на сервере. Обратитесь в техническую поддержку\n"

		if buf.String() != expected {
			t.Errorf("Expected output: %s, got: %s", expected, buf.String())
		}
	})
}
