package server

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Srv    *grpc.Server
	listen net.Listener
}

func NewServer(addr string) (*Server, error) {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	srv := grpc.NewServer()

	return &Server{Srv: srv, listen: listen}, nil
}

func (s *Server) Start(ctx context.Context) {

	slog.Info("server started:" + s.listen.Addr().String())

	if err := s.Srv.Serve(s.listen); err != nil {
		slog.Info(err.Error())
		ctx.Done()
	}
}

func (s *Server) Stop() {
	s.Srv.GracefulStop()
	slog.Info("Server is graceful shutdown...")
}
