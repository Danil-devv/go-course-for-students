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
	"homework10/internal/users"
	"net"
	"testing"
	"time"
)

type TestUser struct {
	Name string
	In   users.User
	Out  any
}

func TestGRPCUsers(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcPort.UnaryServerInterceptor))
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)

	tests := [...]TestUser{
		{Name: "create",
			In: users.User{ID: 15,
				Nickname: "Oleg",
				Email:    "email@example.com"},
			Out: users.User{ID: 15,
				Nickname: "Oleg",
				Email:    "email@example.com"}},
		{Name: "get",
			In: users.User{ID: 15},
			Out: users.User{ID: 15,
				Nickname: "Oleg",
				Email:    "email@example.com"}},
		{Name: "delete",
			In: users.User{ID: 15}},
	}

	for _, test := range tests {
		switch test.Name {
		case "create":
			res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{
				Name:  test.In.Nickname,
				Id:    test.In.ID,
				Email: test.In.Email,
			})
			assert.NoError(t, err, "client.CreateUser")
			assert.Equal(t, test.Out.(users.User).Nickname, res.Name)
			assert.Equal(t, test.Out.(users.User).ID, res.Id)
		case "get":
			res, err := client.GetUser(ctx, &grpcPort.GetUserRequest{
				Id: test.In.ID,
			})
			assert.NoError(t, err, "client.GetUser")
			assert.Equal(t, test.Out.(users.User).Nickname, res.Name)
			assert.Equal(t, test.Out.(users.User).ID, res.Id)
		case "delete":
			_, err := client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{
				Id: test.In.ID,
			})
			assert.NoError(t, err, "client.DeleteUser")
		}
	}

}
