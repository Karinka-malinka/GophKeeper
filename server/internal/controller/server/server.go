package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Srv    *grpc.Server
	listen net.Listener
}

func NewServer(addr string) *Server {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// создаём gRPC-сервер без зарегистрированной службы
	srv := grpc.NewServer()

	return &Server{Srv: srv, listen: listen}
}

func (s *Server) Start() {

	if err := s.Srv.Serve(s.listen); err != nil {
		log.Fatal(err)
	}

	log.Println("server started:", s.listen.Addr())

}

func (s *Server) Stop() {

	s.Srv.GracefulStop()
	log.Println("Server is graceful shutdown...")
}
