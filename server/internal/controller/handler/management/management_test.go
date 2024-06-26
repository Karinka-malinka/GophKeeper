package management

import (
	"context"
	"fmt"
	"testing"

	"github.com/GophKeeper/api/proto"
	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/database/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const userID = "2d53b244-8844-40a6-ab37-e5b89019af0a"

type Request struct {
	loginData *pb.LoginData
	//TODO: здесь добавить еще типы запросов
}

type MocksStore struct {
	LoginData *mocks.MockILoginDataStore
	//TODO: здесь добавить еще моки на другие типы инфорации
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
	//TODO: здесь добавить еще типы запросов

	// Создаем мок объект
	mocksStore.LoginData = mocks.NewMockILoginDataStore(ctrl)
	loginDataApp := logindata.NewLoginData(mocksStore.LoginData)

	managementHandler := NewManagementHandler(loginDataApp)

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
