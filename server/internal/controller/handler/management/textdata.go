package management

import (
	"context"
	"errors"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/textdata"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AddText добавляет новую приватную информацию и возвращает ответ с добавленными данными типа pb.Text.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) AddText(ctx context.Context, in *pb.Text) (*pb.Text, error) {

	ca := make(chan *pb.Text)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		data := textdata.TextData{Text: in.Text, Meta: in.Meta, UserID: userID}
		nData, err := m.TextDataApp.Add(ctx, data)

		if err != nil {
			errc <- err
			return
		}

		resp := &pb.Text{
			Uid:     nData.UUID.String(),
			Text:    nData.Text,
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

// DeleteTextData удаляет приватную информацию и возвращает пустой ответ типа emptypb.Empty.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) DeleteText(ctx context.Context, in *pb.Text) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		_, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		if err := m.TextDataApp.Delete(ctx, in.Uid); err != nil {
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
