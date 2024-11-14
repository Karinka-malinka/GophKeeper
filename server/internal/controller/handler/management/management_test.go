package management

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/GophKeeper/api/proto"
	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/bankcard"
	"github.com/GophKeeper/server/internal/app/filedata"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/app/textdata"
	"github.com/GophKeeper/server/internal/database/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const userID = "2d53b244-8844-40a6-ab37-e5b89019af0a"

type Request struct {
	loginData *pb.LoginData
	textData  *pb.Text
	fileData  *pb.File
	bankCard  *pb.BankCard
}

type MocksStore struct {
	LoginData    *mocks.MockILoginDataStore
	TextData     *mocks.MockITextDataStore
	FileData     *mocks.MockIFileDataStore
	BankCardData *mocks.MockIBankCardDataStore
}

func getHandler(t *testing.T) (context.Context, *ManagementServer, *MocksStore, *Request) {

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := Request{}
	mocksStore := MocksStore{}

	request.loginData = &proto.LoginData{
		Uid:      uuid.New().String(),
		Login:    []byte("testuser"),
		Password: []byte("testpassword"),
		Meta:     []byte("test meta"),
	}

	request.textData = &proto.Text{
		Uid:  uuid.New().String(),
		Text: []byte("test text"),
		Meta: []byte("test meta"),
	}

	request.fileData = &proto.File{
		Uid:  uuid.New().String(),
		File: []byte("test file"),
		Name: []byte("test name"),
		Meta: []byte("test meta"),
	}

	request.bankCard = &proto.BankCard{
		Number: []byte("test number"),
		Term:   []byte("01.2051"),
		Ccv:    []byte("777"),
		Meta:   []byte("test meta"),
	}

	// Создаем мок объект
	mocksStore.LoginData = mocks.NewMockILoginDataStore(ctrl)
	loginDataApp := logindata.NewLoginData(mocksStore.LoginData)

	mocksStore.TextData = mocks.NewMockITextDataStore(ctrl)
	textDataApp := textdata.NewTextData(mocksStore.TextData)

	mocksStore.FileData = mocks.NewMockIFileDataStore(ctrl)
	fileDataApp := filedata.NewFiletData(mocksStore.FileData)

	mocksStore.BankCardData = mocks.NewMockIBankCardDataStore(ctrl)
	bankCardApp := bankcard.NewBankCardData(mocksStore.BankCardData)

	managementHandler := NewManagementHandler(loginDataApp, textDataApp, fileDataApp, bankCardApp)

	return ctx, managementHandler, &mocksStore, &request
}

func TestManagementServer_AddLoginData(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.AddLoginData(ctx, request.loginData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.AddLoginData(ctx, request.loginData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.AddLoginData(ctx, request.loginData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_EditLoginData(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.EditLoginData(ctx, request.loginData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.EditLoginData(ctx, request.loginData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.EditLoginData(ctx, request.loginData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_DeleteLoginData(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.DeleteLoginData(ctx, request.loginData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().Delete(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.DeleteLoginData(ctx, request.loginData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().Delete(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.DeleteLoginData(ctx, request.loginData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_AddText(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.AddText(ctx, request.textData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.TextData.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.AddText(ctx, request.textData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.TextData.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.AddText(ctx, request.textData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_DeleteText(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.DeleteText(ctx, request.textData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.TextData.EXPECT().Delete(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.DeleteText(ctx, request.textData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.TextData.EXPECT().Delete(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.DeleteText(ctx, request.textData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_AddFile(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.AddFile(ctx, request.fileData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.FileData.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.AddFile(ctx, request.fileData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.FileData.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.AddFile(ctx, request.fileData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_DeleteFile(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.DeleteFile(ctx, request.fileData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.FileData.EXPECT().Delete(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.DeleteFile(ctx, request.fileData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.FileData.EXPECT().Delete(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.DeleteFile(ctx, request.fileData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_GetFile(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	respDB := &filedata.FileData{
		UUID:    uuid.New(),
		Created: time.Now().UTC(),
		Name:    []byte("name 1"),
		Meta:    []byte("meta 1")}

	expacted := &pb.File{
		Uid:     respDB.UUID.String(),
		Name:    respDB.Name,
		Meta:    respDB.Meta,
		Created: timestamppb.New(respDB.Created),
	}

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.GetFile(ctx, request.fileData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.FileData.EXPECT().GetFile(ctx, gomock.Any()).Return(respDB, fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.GetFile(ctx, request.fileData)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.FileData.EXPECT().GetFile(ctx, gomock.Any()).Return(respDB, nil).Times(1)

		response, err := managementHandler.GetFile(ctx, request.fileData)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, response, expacted)
	})
}

func TestManagementServer_AddBankCard(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.AddBankCard(ctx, request.bankCard)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.BankCardData.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.AddBankCard(ctx, request.bankCard)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.BankCardData.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.AddBankCard(ctx, request.bankCard)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}

func TestManagementServer_DeleteBankCard(t *testing.T) {

	ctx, managementHandler, mocksStore, request := getHandler(t)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := managementHandler.DeleteBankCard(ctx, request.bankCard)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.BankCardData.EXPECT().Delete(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := managementHandler.DeleteBankCard(ctx, request.bankCard)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.BankCardData.EXPECT().Delete(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := managementHandler.DeleteBankCard(ctx, request.bankCard)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}
