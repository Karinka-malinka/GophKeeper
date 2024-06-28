package management

import (
	"context"
	"errors"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/filedata"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AddFile добавляет новую бинарную информацию и возвращает ответ с добавленными данными типа pb.File.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) AddFile(ctx context.Context, in *pb.File) (*pb.File, error) {

	ca := make(chan *pb.File)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		data := filedata.FileData{File: in.File, Name: in.Name, Meta: in.Meta, UserID: userID}
		nData, err := m.FileDataApp.Add(ctx, data)

		if err != nil {
			errc <- err
			return
		}

		resp := &pb.File{
			Uid:     nData.UUID.String(),
			File:    nData.File,
			Name:    nData.Name,
			Meta:    nData.Meta,
			Created: timestamppb.New(nData.Created),
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

// AddFile добавляет новую бинарную информацию и возвращает ответ с добавленными данными типа pb.File.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) GetFile(ctx context.Context, in *pb.File) (*pb.File, error) {

	ca := make(chan *pb.File)
	errc := make(chan error)

	go func() {

		_, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		nData, err := m.FileDataApp.GetFile(ctx, in.Uid)

		if err != nil {
			errc <- err
			return
		}

		resp := &pb.File{
			Uid:     nData.UUID.String(),
			File:    nData.File,
			Name:    nData.Name,
			Meta:    nData.Meta,
			Created: timestamppb.New(nData.Created),
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

// DeleteFileData удаляет бинарную информацию и возвращает пустой ответ типа emptypb.Empty.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) DeleteFile(ctx context.Context, in *pb.File) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		_, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		if err := m.FileDataApp.Delete(ctx, in.Uid); err != nil {
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
