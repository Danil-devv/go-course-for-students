package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/adapters/usersrepo"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
)

func TestGRRPCCreateUser(t *testing.T) {
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
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Id: 15, Email: "email@exmple.com"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
}

func TestGRPCCreateAd(t *testing.T) {
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
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "hello", res.Title)
	assert.Equal(t, "world", res.Text)
	assert.Equal(t, int64(111), res.AuthorId)
}

func TestGRPCChangeAdStatus(t *testing.T) {
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
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "hello", res.Title)
	assert.Equal(t, "world", res.Text)
	assert.Equal(t, int64(111), res.AuthorId)
	assert.Equal(t, false, res.Published)

	res, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.AuthorId,
		Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	assert.Equal(t, "hello", res.Title)
	assert.Equal(t, "world", res.Text)
	assert.Equal(t, int64(111), res.AuthorId)
	assert.Equal(t, true, res.Published)
}

func TestGRPCUpdateAd(t *testing.T) {
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
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "hello", res.Title)
	assert.Equal(t, "world", res.Text)
	assert.Equal(t, int64(111), res.AuthorId)
	assert.Equal(t, false, res.Published)

	res, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: res.Id, UserId: res.AuthorId, Title: "good bye",
		Text: res.Text})
	assert.NoError(t, err, "client.UpdateAd")

	assert.Equal(t, "good bye", res.Title)
	assert.Equal(t, "world", res.Text)
	assert.Equal(t, int64(111), res.AuthorId)
	assert.Equal(t, false, res.Published)
}

func TestGRPCListAds(t *testing.T) {
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
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(t, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.AuthorId,
		Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	res, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "test", Text: "pupupu", UserId: 111})
	assert.NoError(t, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.AuthorId,
		Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	res, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "some title", Text: "some text", UserId: 115})
	assert.NoError(t, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.AuthorId,
		Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	ads, err := client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err, "client.ListAds")

	assert.Len(t, ads.List, 3)
}

func TestGRPCGetUser(t *testing.T) {
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
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Id: 15, Email: "email@exmple.com"})
	assert.NoError(t, err, "client.CreateUser")

	assert.Equal(t, "Oleg", res.Name)

	res, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 15})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, int64(15), res.Id)
}

func TestGRPCDeleteUser(t *testing.T) {
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
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Id: 15, Email: "email@exmple.com"})
	assert.NoError(t, err, "client.CreateUser")

	assert.Equal(t, "Oleg", res.Name)

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: 15})
	assert.NoError(t, err, "client.DeleteUser")

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 15})
	assert.Error(t, err, "client.GetUSer")
}

func TestGRPCDeleteAd(t *testing.T) {
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
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "hello", res.Title)
	assert.Equal(t, "world", res.Text)
	assert.Equal(t, int64(111), res.AuthorId)

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: res.Id, AuthorId: res.AuthorId})
	assert.NoError(t, err, "client.DeleteAd")

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: res.Id, UserId: res.AuthorId, Title: "good bye",
		Text: res.Text})
	assert.Error(t, err, "client.GetAd")
}
