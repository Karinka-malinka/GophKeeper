package management

import (
	"context"

	pb "github.com/GophKeeper/server/cmd/proto"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ManagementServer struct {
	pb.UnimplementedUserServiceServer
	LoginDataApp *logindata.LoginDatas
}

// NewManagementHandler создает новый экземпляр ManagementServer.
func NewManagementHandler(loginData *logindata.LoginDatas) *ManagementServer {
	return &ManagementServer{
		LoginDataApp: loginData,
	}
}

func (m *ManagementServer) AddLoginData(ctx context.Context, in *pb.LoginData) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		userID := user.GetUserID(ctx)

		loginData := logindata.LoginData{Login: in.Login, Password: in.Password, UserID: userID}

		if err := m.LoginDataApp.Add(ctx, loginData); err != nil {
			errc <- err
			return
		}

		ca <- true
	}()

	select {
	case <-ca:
		return &emptypb.Empty{}, nil
	case <-errc:
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}
