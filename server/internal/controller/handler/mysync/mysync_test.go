package mysync

import (
	"context"
	"fmt"
	"testing"

	logindataApp "github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/database/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const userID = "2d53b244-8844-40a6-ab37-e5b89019af0a"

func TestSyncServer_ListLoginData(t *testing.T) {

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок объект
	mockLogindataStore := mocks.NewMockILoginDataStore(ctrl)
	loginDataApp := logindataApp.NewLoginData(mockLogindataStore)

	syncHandler := NewSyncHandler(loginDataApp)

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := syncHandler.ListLoginData(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mockLogindataStore.EXPECT().GetList(ctx, userID).Return(nil, fmt.Errorf("Chren")).Times(1)

		_, err := syncHandler.ListLoginData(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mockLogindataStore.EXPECT().GetList(ctx, userID).Return(nil, nil).Times(1)

		response, err := syncHandler.ListLoginData(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}
