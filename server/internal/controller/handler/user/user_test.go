package user

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/cmd/config"
	userApp "github.com/GophKeeper/server/internal/app/user"
	"github.com/GophKeeper/server/internal/database"
	"github.com/GophKeeper/server/internal/database/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//const userID = "2d53b244-8844-40a6-ab37-e5b89019af0a"

func TestRegister(t *testing.T) {

	ctx := context.Background()

	request := &proto.UserRequest{
		Login:    "testuser",
		Password: "testpassword",
	}

	cfg := config.NewConfig()

	cfg.PathKeys = "../../../../cmd/config/config_example.json"
	cfgToken, err := cfg.GetConfig(cfg.PathKeys)
	if err != nil {
		log.Fatalf("cnf is not set. Error: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок объект
	mockUserStore := mocks.NewMockIUserStore(ctrl)

	userapp := userApp.NewUser(mockUserStore, cfgToken)
	userHandler := NewUserHandler(userapp)

	t.Run("Success", func(t *testing.T) {

		mockUserStore.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

		response, err := userHandler.Register(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	t.Run("Conflict", func(t *testing.T) {

		mockUserStore.EXPECT().Create(ctx, gomock.Any()).Return(&database.ConflictError{}).Times(1)

		_, err := userHandler.Register(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err), codes.AlreadyExists)
	})

	t.Run("Internal", func(t *testing.T) {

		mockUserStore.EXPECT().Create(ctx, gomock.Any()).Return(fmt.Errorf("Chren")).Times(1)

		_, err := userHandler.Register(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err), codes.Internal)
	})
}

func TestLogin(t *testing.T) {

	ctx := context.Background()

	request := &proto.UserRequest{
		Login:    "testuser",
		Password: "hashpassword",
	}

	expacted := &userApp.User{
		Username: "testuser",
		Password: "hashpassword",
	}

	cfg := config.NewConfig()

	cfg.PathKeys = "../../../../cmd/config/config_example.json"
	cfgToken, err := cfg.GetConfig(cfg.PathKeys)
	if err != nil {
		log.Fatalf("cnf is not set. Error: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок объект
	mockUserStore := mocks.NewMockIUserStore(ctrl)

	userapp := userApp.NewUser(mockUserStore, cfgToken)
	userHandler := NewUserHandler(userapp)

	t.Run("Success", func(t *testing.T) {

		mockUserStore.EXPECT().GetUser(ctx, request.Login).Return(expacted, nil).Times(1)

		response, err := userHandler.Login(ctx, request)

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})

	t.Run("Unauthenticated", func(t *testing.T) {

		expactedUnauthenticated := &userApp.User{
			Username: "testuser",
			Password: "hashpass",
		}

		mockUserStore.EXPECT().GetUser(ctx, request.Login).Return(expactedUnauthenticated, nil).Times(1)

		_, err := userHandler.Login(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err), codes.Unauthenticated)
	})

	t.Run("Internal", func(t *testing.T) {

		mockUserStore.EXPECT().GetUser(ctx, request.Login).Return(nil, fmt.Errorf("Chren")).Times(1)

		_, err := userHandler.Login(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, status.Code(err), codes.Internal)
	})
}
