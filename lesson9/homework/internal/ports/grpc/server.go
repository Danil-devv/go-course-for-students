package grpc

import (
	"context"
	"google.golang.org/grpc"
	"homework9/internal/app"
	"log"
	"net"
)

type Server struct {
	srv *grpc.Server
	lis net.Listener
}

func NewGRPCServer(port string, a *app.App) Server {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := NewService(*a)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(UnaryServerInterceptor))

	RegisterAdServiceServer(server, svc)

	return Server{
		srv: server,
		lis: lis,
	}
}

func (s *Server) Listen() error {
	return s.srv.Serve(s.lis)
}

func (s *Server) GracefulShutdown() {
	s.srv.GracefulStop()
	_ = s.lis.Close()
}
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println(info.FullMethod)

	return handler(ctx, req)
}
