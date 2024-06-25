// Package user предоставляет реализацию методов для работы с пользователями
package user

import (
	"context"
	"errors"

	pb "github.com/GophKeeper/server/cmd/proto"
	"github.com/GophKeeper/server/internal/app/user"
	"github.com/GophKeeper/server/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UsersServer представляет обработчик для работы с пользователями.
type UsersServer struct {
	pb.UnimplementedUserServiceServer
	UserApp *user.Users
}

// NewUserHandler создает новый экземпляр UsersServer.
func NewUserHandler(userapp *user.Users) *UsersServer {
	return &UsersServer{UserApp: userapp}
}

// Register регистрирует нового пользователя и возвращает ответ с токеном и идентификатором.
// Возможные коды ошибок:
//   - codes.AlreadyExists: Если пользователь с таким именем уже существует.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (u *UsersServer) Register(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {

	ca := make(chan *user.User)
	errc := make(chan error)

	go func() {

		user, err := u.UserApp.Register(ctx, user.User{Username: in.Login, Password: in.Password})

		if err != nil {
			errc <- err
			//return
		}

		ca <- user
	}()

	select {
	case result := <-ca:
		return &pb.UserResponse{Token: result.Token, Uid: result.UUID.String(), Key: u.UserApp.Cfg.Key}, nil
	case err := <-errc:
		var errConflict *database.ConflictError
		if errors.As(err, &errConflict) {
			return nil, status.Errorf(codes.AlreadyExists, "")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}

// Login выполняет процесс аутентификации пользователя и возвращает ответ с токеном и идентификатором.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (u *UsersServer) Login(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {

	ca := make(chan *user.User)
	errc := make(chan error)

	go func() {

		user, err := u.UserApp.Login(ctx, user.User{Username: in.Login, Password: in.Password})

		if err != nil {
			errc <- err
			return
		}

		ca <- user
	}()

	select {
	case result := <-ca:
		return &pb.UserResponse{Token: result.Token, Uid: result.UUID.String(), Key: u.UserApp.Cfg.Key}, nil
	case err := <-errc:
		if err.Error() == "401" {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}

}
