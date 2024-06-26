// Package server предоставляет реализацию сервера gRPC и REST API.
package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server представляет сервер с gRPC и REST API.
type Server struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedManagementServiceServer
	pb.UnimplementedSyncServiceServer
	Srv     *grpc.Server
	listen  net.Listener
	SrvREST *http.Server
}

// NewServer создает новый сервер с gRPC и REST API.
func NewServer(ctx context.Context, addrgRPS, addrRESR string, userApp *user.Users) (*Server, error) {

	//gRPC
	listen, err := net.Listen("tcp", addrgRPS)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(userApp.TokenAuthMiddlewareGRPS))

	//RESR
	mux := runtime.NewServeMux(runtime.WithMetadata(userApp.TokenAuthMiddlewareREST))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addrgRPS, opts)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	err = pb.RegisterSyncServiceHandlerFromEndpoint(ctx, mux, addrgRPS, opts)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	err = pb.RegisterManagementServiceHandlerFromEndpoint(ctx, mux, addrgRPS, opts)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	srvREST := http.Server{
		Addr:              addrRESR,
		Handler:           mux,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}

	return &Server{Srv: srv, listen: listen, SrvREST: &srvREST}, nil
}

// Start запускает сервер gRPC.
func (s *Server) Start(ctx context.Context) {

	slog.Info("gRPS server started:" + s.listen.Addr().String())

	if err := s.Srv.Serve(s.listen); err != nil {
		slog.Info(err.Error())
		ctx.Done()
	}
}

// Stop останавливает сервер gRPC.
func (s *Server) Stop() {
	s.Srv.GracefulStop()
	slog.Info("Server gRPS is graceful shutdown...")
}

// StartRest запускает сервер REST API.
func (s *Server) StartRest(ctx context.Context) {

	slog.Info("RESR server started:" + s.SrvREST.Addr)

	if err := s.SrvREST.ListenAndServe(); err != nil {
		slog.Info(err.Error())
		ctx.Done()
	}
}

// StopREST останавливает сервер REST API.
func (s *Server) StopREST(ctx context.Context) {

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*2))
	defer cancel()

	err := s.SrvREST.Shutdown(timeoutCtx)
	if err != nil {
		slog.Info(err.Error())
	}
	slog.Info("Server RESR is graceful shutdown...")

}
