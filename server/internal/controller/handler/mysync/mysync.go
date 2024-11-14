// Package mysync предоставляет реализацию методов для синхронизации данных
package mysync

import (
	"context"
	"errors"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/bankcard"
	"github.com/GophKeeper/server/internal/app/filedata"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/app/textdata"
	"github.com/GophKeeper/server/internal/app/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UsersServer представляет обработчик для синхронизации данных
type SyncServer struct {
	pb.UnimplementedSyncServiceServer
	LoginDataApp    *logindata.LoginDatas
	TextDataApp     *textdata.TextDatas
	FileDataApp     *filedata.FileDatas
	BankCardDataApp *bankcard.BankCardDatas
}

// NewSyncHandler создает новый экземпляр SyncServer.
func NewSyncHandler(loginData *logindata.LoginDatas, textData *textdata.TextDatas, fileData *filedata.FileDatas, bankCard *bankcard.BankCardDatas) *SyncServer {
	return &SyncServer{
		LoginDataApp:    loginData,
		TextDataApp:     textData,
		FileDataApp:     fileData,
		BankCardDataApp: bankCard,
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
		var UnauthorizedError *user.UnauthorizedError
		if errors.As(err, &UnauthorizedError) {
			return nil, status.Errorf(codes.Unauthenticated, "")
		}
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}

// ListText возвращает список приватной текстовой информации для синхронизации в ответе типа pb.TextResponse.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если запрос не авторизован из-за отсутствия прав доступа.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (s *SyncServer) ListText(ctx context.Context, _ *emptypb.Empty) (*pb.TextResponse, error) {

	ca := make(chan *pb.TextResponse)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		var response pb.TextResponse

		sData, err := s.TextDataApp.GetList(ctx, userID)

		if err != nil {
			errc <- err
			return
		}

		for _, data := range sData {
			response.Text = append(response.Text, &pb.Text{
				Uid:     data.UUID.String(),
				Text:    data.Text,
				Meta:    data.Meta,
				Created: timestamppb.New(data.Created),
			})
		}

		ca <- &response
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

// ListFile возвращает список приватной бинарной информации для синхронизации в ответе типа pb.FileResponse.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если запрос не авторизован из-за отсутствия прав доступа.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (s *SyncServer) ListFile(ctx context.Context, _ *emptypb.Empty) (*pb.FileResponse, error) {

	ca := make(chan *pb.FileResponse)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		var response pb.FileResponse

		sData, err := s.FileDataApp.GetList(ctx, userID)

		if err != nil {
			errc <- err
			return
		}

		for _, data := range sData {
			response.File = append(response.File, &pb.File{
				Uid:     data.UUID.String(),
				Name:    data.Name,
				Meta:    data.Meta,
				Created: timestamppb.New(data.Created),
			})
		}

		ca <- &response
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

// ListBankCard возвращает список банковских карт для синхронизации в ответе типа pb.BankCardResponse.
// Возможные коды ошибок:
//   - codes.Unauthenticated: Если запрос не авторизован из-за отсутствия прав доступа.
//   - codes.Internal: В случае внутренней ошибки сервера.
//   - codes.Aborted: Если операция была прервана контекстом.
func (s *SyncServer) ListBankCard(ctx context.Context, _ *emptypb.Empty) (*pb.BankCardResponse, error) {

	ca := make(chan *pb.BankCardResponse)
	errc := make(chan error)

	go func() {

		userID, err := user.GetUserID(ctx)
		if err != nil {
			errc <- err
			return
		}

		var response pb.BankCardResponse

		sData, err := s.BankCardDataApp.GetList(ctx, userID)

		if err != nil {
			errc <- err
			return
		}

		for _, data := range sData {
			response.BankCard = append(response.BankCard, &pb.BankCard{
				Number:  data.Number,
				Term:    data.Term,
				Ccv:     data.CCV,
				Meta:    data.Meta,
				Created: timestamppb.New(data.Created),
			})
		}

		ca <- &response
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
