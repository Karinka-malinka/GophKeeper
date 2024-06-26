// Package mysync предоставляет реализацию методов для синхронизации данных
package mysync

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UsersServer представляет обработчик для синхронизации данных
type SyncServer struct {
	pb.UnimplementedSyncServiceServer
	LoginDataApp *logindata.LoginDatas
	//TODO: добавить app на все типы хранимых данных
	//TextDataApp *textdata.TextDatas
}

// NewSyncHandler создает новый экземпляр SyncServer.
// TODO: добавить app на все типы хранимых данных
func NewSyncHandler(loginData *logindata.LoginDatas) *SyncServer {
	return &SyncServer{
		LoginDataApp: loginData,
	}
}

// ListLoginData возвращает список данных для синхронизации в ответе типа pb.LoginDataResponse.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если запрос не авторизован из-за отсутствия прав доступа.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (s *SyncServer) ListLoginData(ctx context.Context, _ *emptypb.Empty) (*pb.LoginDataResponse, error) {

	ca := make(chan *pb.LoginDataResponse)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

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
	case err := <-errc:
		if errors.Is(err, fmt.Errorf("Unauthorized")) {
			return nil, status.Errorf(codes.Unauthenticated, "")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}
