package mysync

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const userID = "2d53b244-8844-40a6-ab37-e5b89019af0a"

type MocksStore struct {
	LoginData    *mocks.MockILoginDataStore
	TextData     *mocks.MockITextDataStore
	FileData     *mocks.MockIFileDataStore
	BankCardData *mocks.MockIBankCardDataStore
}

func getHandler(t *testing.T) (*SyncServer, *MocksStore) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок объект
	mocksStore := MocksStore{}

	mocksStore.LoginData = mocks.NewMockILoginDataStore(ctrl)
	loginDataApp := logindata.NewLoginData(mocksStore.LoginData)

	mocksStore.TextData = mocks.NewMockITextDataStore(ctrl)
	textDataApp := textdata.NewTextData(mocksStore.TextData)

	mocksStore.FileData = mocks.NewMockIFileDataStore(ctrl)
	fileDataApp := filedata.NewFiletData(mocksStore.FileData)

	mocksStore.BankCardData = mocks.NewMockIBankCardDataStore(ctrl)
	bankCardApp := bankcard.NewBankCardData(mocksStore.BankCardData)

	syncHandler := NewSyncHandler(loginDataApp, textDataApp, fileDataApp, bankCardApp)

	return syncHandler, &mocksStore
}

func TestSyncServer_ListLoginData(t *testing.T) {

	ctx := context.Background()

	syncHandler, mocksStore := getHandler(t)

	respDB := []logindata.LoginData{}
	respDB = append(respDB, logindata.LoginData{
		UUID:     uuid.New(),
		Created:  time.Now().UTC(),
		Login:    []byte("login1"),
		Password: []byte("pass1"),
		Meta:     []byte("meta 1")})
	respDB = append(respDB, logindata.LoginData{
		UUID:     uuid.New(),
		Created:  time.Now().UTC(),
		Login:    []byte("login2"),
		Password: []byte("pass2"),
		Meta:     []byte("meta 2")})

	expacted := pb.LoginDataResponse{}
	for _, loginData := range respDB {
		expacted.LoginData = append(expacted.LoginData, &pb.LoginData{
			Uid:      loginData.UUID.String(),
			Login:    loginData.Login,
			Password: loginData.Password,
			Meta:     loginData.Meta,
			Created:  timestamppb.New(loginData.Created),
		})
	}

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := syncHandler.ListLoginData(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().GetList(ctx, userID).Return(nil, fmt.Errorf("Chren")).Times(1)

		_, err := syncHandler.ListLoginData(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.LoginData.EXPECT().GetList(ctx, userID).Return(respDB, nil).Times(1)

		response, err := syncHandler.ListLoginData(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, response, &expacted)
	})
}

func TestSyncServer_ListTextData(t *testing.T) {

	ctx := context.Background()

	syncHandler, mocksStore := getHandler(t)

	respDB := []textdata.TextData{}
	respDB = append(respDB, textdata.TextData{
		UUID:    uuid.New(),
		Created: time.Now().UTC(),
		Text:    []byte("text 1"),
		Meta:    []byte("meta 1")})
	respDB = append(respDB, textdata.TextData{
		UUID:    uuid.New(),
		Created: time.Now().UTC(),
		Text:    []byte("text 2"),
		Meta:    []byte("meta 2")})

	expacted := pb.TextResponse{}
	for _, data := range respDB {
		expacted.Text = append(expacted.Text, &pb.Text{
			Uid:     data.UUID.String(),
			Text:    data.Text,
			Meta:    data.Meta,
			Created: timestamppb.New(data.Created),
		})
	}

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := syncHandler.ListText(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.TextData.EXPECT().GetList(ctx, userID).Return(nil, fmt.Errorf("Chren")).Times(1)

		_, err := syncHandler.ListText(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.TextData.EXPECT().GetList(ctx, userID).Return(respDB, nil).Times(1)

		response, err := syncHandler.ListText(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, response, &expacted)
	})
}

func TestSyncServer_ListFileData(t *testing.T) {

	ctx := context.Background()

	syncHandler, mocksStore := getHandler(t)

	respDB := []filedata.FileData{}
	respDB = append(respDB, filedata.FileData{
		UUID:    uuid.New(),
		Created: time.Now().UTC(),
		Name:    []byte("name 1"),
		Meta:    []byte("meta 1")})
	respDB = append(respDB, filedata.FileData{
		UUID:    uuid.New(),
		Created: time.Now().UTC(),
		Name:    []byte("name 2"),
		Meta:    []byte("meta 2")})

	expacted := pb.FileResponse{}
	for _, data := range respDB {
		expacted.File = append(expacted.File, &pb.File{
			Uid:     data.UUID.String(),
			Name:    data.Name,
			Meta:    data.Meta,
			Created: timestamppb.New(data.Created),
		})
	}

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := syncHandler.ListFile(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.FileData.EXPECT().GetList(ctx, userID).Return(nil, fmt.Errorf("Chren")).Times(1)

		_, err := syncHandler.ListFile(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.FileData.EXPECT().GetList(ctx, userID).Return(respDB, nil).Times(1)

		response, err := syncHandler.ListFile(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, response, &expacted)
	})
}

func TestSyncServer_ListBankCard(t *testing.T) {

	ctx := context.Background()

	syncHandler, mocksStore := getHandler(t)

	respDB := []bankcard.BankCardData{}
	respDB = append(respDB, bankcard.BankCardData{
		Number:  []byte("number 1"),
		Created: time.Now().UTC(),
		Term:    []byte("term 1"),
		CCV:     []byte("ccv 1"),
		Meta:    []byte("meta 1")})
	respDB = append(respDB, bankcard.BankCardData{
		Number:  []byte("number 2"),
		Created: time.Now().UTC(),
		Term:    []byte("term 2"),
		CCV:     []byte("ccv 2"),
		Meta:    []byte("meta 2")})

	expacted := pb.BankCardResponse{}
	for _, data := range respDB {
		expacted.BankCard = append(expacted.BankCard, &pb.BankCard{
			Number:  data.Number,
			Term:    data.Term,
			Ccv:     data.CCV,
			Meta:    data.Meta,
			Created: timestamppb.New(data.Created),
		})
	}

	t.Run("Unauthenticated", func(t *testing.T) {

		_, err := syncHandler.ListBankCard(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Unauthenticated.String())
	})

	ctx = context.WithValue(ctx, "userID", userID)

	t.Run("Internal", func(t *testing.T) {

		mocksStore.BankCardData.EXPECT().GetList(ctx, userID).Return(nil, fmt.Errorf("Chren")).Times(1)

		_, err := syncHandler.ListBankCard(ctx, &emptypb.Empty{})

		assert.Error(t, err)
		assert.Equal(t, status.Code(err).String(), codes.Internal.String())
	})

	t.Run("Success", func(t *testing.T) {

		mocksStore.BankCardData.EXPECT().GetList(ctx, userID).Return(respDB, nil).Times(1)

		response, err := syncHandler.ListBankCard(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, response, &expacted)
	})
}
