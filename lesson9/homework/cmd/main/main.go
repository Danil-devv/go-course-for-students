package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/usersrepo"
	"homework9/internal/app"
	"homework9/internal/ports/grpc"
	"homework9/internal/ports/httpgin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	grpcPort = ":18020"
	httpPort = ":18080"
)

func main() {
	a := app.NewApp(adrepo.New(), usersrepo.New())

	httpServer := httpgin.NewHTTPServer(httpPort, &a)
	grpcServer := grpc.NewGRPCServer(grpcPort, &a)

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
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
		log.Printf("starting grpc server, listening on %s\n", grpcPort)
		defer log.Printf("close grpc server listening on %s\n", grpcPort)

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulShutdown()

			close(errCh)
		}()

		go func() {
			if err := grpcServer.Listen(); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})

	// run http server
	eg.Go(func() error {
		log.Printf("starting http server, listening on %s\n", httpPort)
		defer log.Printf("close http server listening on %s\n", httpPort)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.GracefulShutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", httpPort, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := httpServer.Listen(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}
