package ggrpc

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

type Server struct {
	grpcServer *grpc.Server
	port       string
}

func NewServer(p string) *Server {
	return &Server{
		grpcServer: grpc.NewServer(),
		port:       p,
	}
}

func (s *Server) Start() error {
	b := strings.Builder{}
	b.WriteString(":")
	b.WriteString(s.port)

	lis, err := net.Listen("tcp", b.String())
	if err != nil {
		return err
	}
	b.Reset()

	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	log.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}
