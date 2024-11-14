package management

import (
	"context"
	"errors"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AddLoginData добавляет новую пару логин/пароль и возвращает ответ с добавленными данными типа pb.LoginData.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) AddLoginData(ctx context.Context, in *pb.LoginData) (*pb.LoginData, error) {

	ca := make(chan *pb.LoginData)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		loginData := logindata.LoginData{Login: in.Login, Password: in.Password, Meta: in.Meta, UserID: userID}
		nloginData, err := m.LoginDataApp.Add(ctx, loginData)

		if err != nil {
			errc <- err
			return
		}

		resp := &pb.LoginData{
			Uid:      nloginData.UUID.String(),
			Login:    nloginData.Login,
			Password: nloginData.Password,
			Meta:     nloginData.Meta,
			Created:  timestamppb.New(nloginData.Created),
		}

		ca <- resp
	}()

	select {
	case result := <-ca:
		return result, nil
	case err := <-errc:
		var UnauthorizedError *user.UnauthorizedError
		if errors.As(err, &UnauthorizedError) {
			return nil, status.Errorf(codes.Unauthenticated, "")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}

// EditLoginData изменяет пару логин/пароль и возвращает пустой ответ типа emptypb.Empty.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) EditLoginData(ctx context.Context, in *pb.LoginData) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		_, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		if err := m.LoginDataApp.Edit(ctx, in.Uid, in.Password); err != nil {
			errc <- err
			return
		}

		ca <- true
	}()

	select {
	case <-ca:
		return &emptypb.Empty{}, nil
	case err := <-errc:
		var UnauthorizedError *user.UnauthorizedError
		if errors.As(err, &UnauthorizedError) {
			return nil, status.Errorf(codes.Unauthenticated, "")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}

// DeleteLoginData удаляет пару логин/пароль и возвращает пустой ответ типа emptypb.Empty.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) DeleteLoginData(ctx context.Context, in *pb.LoginData) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		_, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		if err := m.LoginDataApp.Delete(ctx, in.Uid); err != nil {
			errc <- err
			return
		}

		ca <- true
	}()

	select {
	case <-ca:
		return &emptypb.Empty{}, nil
	case err := <-errc:
		var UnauthorizedError *user.UnauthorizedError
		if errors.As(err, &UnauthorizedError) {
			return nil, status.Errorf(codes.Unauthenticated, "")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}
