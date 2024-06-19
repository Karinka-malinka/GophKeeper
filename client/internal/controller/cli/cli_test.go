package cli

import (
	"context"
	"testing"

	pb "github.com/GophKeeper/server/cmd/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// Mock для ManagementServiceClient
type MockManagementServiceClient struct {
	mock.Mock
}

func (m *MockManagementServiceClient) AddLoginData(ctx context.Context, in *pb.LoginData, _ ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
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

	mockClient := &MockManagementServiceClient{}
	mockClient.On("AddLoginData", ctx, mock.Anything).Return(new(empty.Empty), nil)

	tests := []struct {
		name  string
		cli   *CLI
		words []string
	}{
		{
			name:  "Test case 1",
			cli:   &CLI{},
			words: []string{"addlogindata", "karina", "123", "meta meta meta"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.AddLoginData(ctx, mockClient, tt.words)
		})
	}
}
