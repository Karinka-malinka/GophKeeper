package management

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
	case <-errc:
		return nil, status.Errorf(codes.Internal, "")
	case <-ctx.Done():
		return nil, status.Errorf(codes.Aborted, "")
	}
}

func (m *ManagementServer) EditLoginData(ctx context.Context, in *pb.LoginData) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		if err := m.LoginDataApp.Edit(ctx, in.Uid, in.Password); err != nil {
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

func (m *ManagementServer) DeleteLoginData(ctx context.Context, in *pb.LoginData) (*emptypb.Empty, error) {

	ca := make(chan bool)
	errc := make(chan error)

	go func() {

		if err := m.LoginDataApp.Delete(ctx, in.Uid); err != nil {
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
