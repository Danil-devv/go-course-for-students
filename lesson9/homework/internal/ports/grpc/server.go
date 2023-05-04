package grpc

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework9/internal/app"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

func (s *Server) Listen() {
	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	// run grpc server
	eg.Go(func() error {
		log.Printf("starting gRPC server, listening on %s\n", s.lis.Addr())
		defer log.Printf("close gRPC server listening on %s\n", s.lis.Addr())

		errCh := make(chan error)

		defer func() {
			s.srv.GracefulStop()
			_ = s.lis.Close()

			close(errCh)
		}()

		go func() {
			if err := s.srv.Serve(s.lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("gRPC server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the gRPC server: %s\n", err.Error())
	}

	log.Println("gRPC server were successfully shutdown")
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println(info.FullMethod)

	return handler(ctx, req)
}
