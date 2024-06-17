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

type UsersServer struct {
	pb.UnimplementedUserServiceServer
	UserApp *user.Users
}

// NewURLHandler создает новый экземпляр URLHandler.
func NewUserHandler(userapp *user.Users) *UsersServer {
	return &UsersServer{UserApp: userapp}
}

func (u *UsersServer) Register(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {

	ca := make(chan string)
	errc := make(chan error)

	go func() {

		token, err := u.UserApp.Register(ctx, user.User{Username: in.Login, Password: in.Password})

		if err != nil {
			errc <- err
			return
		}

		ca <- token
	}()

	select {
	case result := <-ca:
		return &pb.UserResponse{Token: result}, nil
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

func (u *UsersServer) Login(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {

	ca := make(chan string)
	errc := make(chan error)

	go func() {

		token, err := u.UserApp.Login(ctx, user.User{Username: in.Login, Password: in.Password})

		if err != nil {
			errc <- err
			return
		}

		ca <- token
	}()

	select {
	case result := <-ca:
		return &pb.UserResponse{Token: result}, nil
	case err := <-errc:
		if err.Error() == "401" {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}

}
