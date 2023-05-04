package main

import (
	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/usersrepo"
	"homework9/internal/app"
	"homework9/internal/ports/grpc"
	"homework9/internal/ports/httpgin"
	"sync"
)

func main() {
	a := app.NewApp(adrepo.New(), usersrepo.New())
	httpServer := httpgin.NewHTTPServer(":18080", &a)
	grpcServer := grpc.NewGRPCServer(":18020", &a)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		httpServer.Listen()
	}()

	go func() {
		defer wg.Done()
		grpcServer.Listen()
	}()

	wg.Wait()
}
