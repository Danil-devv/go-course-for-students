package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/usersrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
	"math"
	"net"
	"strconv"
	"testing"
	"time"
)

func FuzzGRPCCreateUserWithWrongId_Fuzz(f *testing.F) {
	lis := bufconn.Listen(1024 * 1024)
	f.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptor))
	f.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(f, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	f.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(f, err, "grpc.DialContext")

	f.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)

	tests := [...]int64{-1, -5, -100, math.MinInt64}

	for _, test := range tests {
		f.Add(test)
	}

	f.Fuzz(func(t *testing.T, test int64) {
		s := strconv.Itoa(int(test))
		_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{
			Id:    test,
			Name:  "Oleg" + s,
			Email: "email@example.com",
		})

		assert.Error(t, err)
	})
}
