package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	pb "github.com/GophKeeper/api/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Mock для ManagementServiceClient
type MockManagementServiceClient struct {
	mock.Mock
}

func (m *MockManagementServiceClient) AddLoginData(ctx context.Context, in *pb.LoginData, _ ...grpc.CallOption) (*pb.LoginData, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.LoginData), args.Error(1)
}

func (m *MockManagementServiceClient) AddText(ctx context.Context, in *pb.Text, _ ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockManagementServiceClient) AddFile(ctx context.Context, in *pb.File, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockManagementServiceClient) AddBankCard(ctx context.Context, in *pb.BankCard, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockManagementServiceClient) EditLoginData(ctx context.Context, in *pb.LoginData, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockManagementServiceClient) GetFile(ctx context.Context, in *pb.File, opts ...grpc.CallOption) (*pb.File, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.File), args.Error(1)
}

func (m *MockManagementServiceClient) DeleteLoginData(ctx context.Context, in *pb.LoginData, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockManagementServiceClient) DeleteText(ctx context.Context, in *pb.Text, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockManagementServiceClient) DeleteFile(ctx context.Context, in *pb.File, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockManagementServiceClient) DeleteBankCard(ctx context.Context, in *pb.BankCard, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func TestCLI_AddLoginData(t *testing.T) {

	ctx := context.Background()

	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cManegment := pb.NewManagementServiceClient(conn)

	fmt.Println("Соединение с сервером установлено")

	cli := &CLI{UserID: "2d618858-27bb-474c-9150-8a03126feff7", Key: "a2F4c1e8D9d3A7fe"}

	words := []string{"addlogindata", "karina", "123", "meta meta meta"}

	mockClient := &MockManagementServiceClient{}
	mockClient.On("AddLoginData", ctx, mock.Anything).Return(new(empty.Empty), nil)

	tests := []struct {
		name  string
		cli   *CLI
		words []string
	}{
		{
			name:  "Test case 1",
			cli:   &CLI{UserID: "2d618858-27bb-474c-9150-8a03126feff7", Key: "a2F4c1e8D9d3A7fe"},
			words: []string{"addlogindata", "karina", "123", "meta meta meta"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.AddLoginData(ctx, mockClient, tt.words)
		})
	}

	t.Run("Success", func(t *testing.T) {

		mockClient := &MockManagementServiceClient{}
		mockClient.On("AddLoginData", ctx, mock.Anything).Return(new(empty.Empty), nil)

		cli.AddLoginData(ctx, mockClient, words)

		assert.Equal(t, http.StatusUnauthorized, cManegment)

	})
}
