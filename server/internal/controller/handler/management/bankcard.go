package management

import (
	"context"
	"errors"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/bankcard"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AddBankCard добавляет новые данные о банковской карте и возвращает ответ с добавленными данными типа pb.Text.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) AddBankCard(ctx context.Context, in *pb.BankCard) (*pb.BankCard, error) {

	ca := make(chan *pb.BankCard)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		data := bankcard.BankCardData{Number: in.Number, Term: in.Term, CCV: in.Ccv, Meta: in.Meta, UserID: userID}
		nData, err := m.BankCardDataApp.Add(ctx, data)

		if err != nil {
			errc <- err
			return
		}

		resp := &pb.BankCard{
			Number:  nData.Number,
			Term:    nData.Term,
			Ccv:     nData.CCV,
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

// DeleteBankCard удаляет данные о банковской карте и возвращает пустой ответ типа emptypb.Empty.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если предоставлены неверные учетные данные.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (m *ManagementServer) DeleteBankCard(ctx context.Context, in *pb.BankCard) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		_, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		if err := m.BankCardDataApp.Delete(ctx, in.Number); err != nil {
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
