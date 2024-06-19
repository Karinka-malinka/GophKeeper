package mysync

import (
	"context"

	pb "github.com/GophKeeper/server/cmd/proto"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SyncServer struct {
	pb.UnimplementedUserServiceServer
	LoginDataApp *logindata.LoginDatas
}

// NewSyncHandler создает новый экземпляр SyncServer.
func NewSyncHandler(loginData *logindata.LoginDatas) *SyncServer {
	return &SyncServer{
		LoginDataApp: loginData,
	}
}

func (s *SyncServer) ListLoginData(ctx context.Context, _ *emptypb.Empty) (*pb.LoginDataResponse, error) {

	ca := make(chan *pb.LoginDataResponse)
	errc := make(chan error)

	go func() {

		userID := user.GetUserID(ctx)

		var response pb.LoginDataResponse

		sLoginData, err := s.LoginDataApp.GetList(ctx, userID)

		if err != nil {
			errc <- err
			return
		}

		for _, loginData := range sLoginData {
			response.LoginData = append(response.LoginData, &pb.LoginData{
				Uid:      loginData.UUID.String(),
				Login:    loginData.Login,
				Password: loginData.Password,
				Meta:     loginData.Meta,
				Created:  timestamppb.New(loginData.Created),
			})
		}

		ca <- &response
	}()

	select {
	case result := <-ca:
		return result, nil
	case <-errc:
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}
