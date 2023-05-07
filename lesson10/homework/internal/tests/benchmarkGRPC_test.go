package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/adapters/usersrepo"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
)

func BenchmarkGRPCCreateUser(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)

	for i := 0; i < b.N; i++ {
		res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Id: int64(i), Email: "email@exmple.com"})
		assert.NoError(b, err, "client.GetUser")

		assert.Equal(b, "Oleg", res.Name)
	}

}

func BenchmarkGRPCCreateAd(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)

	for i := 0; i < b.N; i++ {
		res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: int64(i)})
		assert.NoError(b, err, "client.CreateAd")

		assert.Equal(b, "hello", res.Title)
		assert.Equal(b, "world", res.Text)
		assert.Equal(b, int64(i), res.AuthorId)
	}

}

func BenchmarkGRPCChangeAdStatus(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(b, err, "client.CreateAd")

	assert.Equal(b, "hello", res.Title)
	assert.Equal(b, "world", res.Text)
	assert.Equal(b, int64(111), res.AuthorId)
	assert.Equal(b, false, res.Published)
	for i := 0; i < b.N; i++ {
		_, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{
			AdId:      res.Id,
			UserId:    res.AuthorId,
			Published: true,
		})
		assert.NoError(b, err, "client.ChangeAdStatus")
	}
}

func BenchmarkGRPCUpdateAd(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(b, err, "client.CreateAd")

	assert.Equal(b, "hello", res.Title)
	assert.Equal(b, "world", res.Text)
	assert.Equal(b, int64(111), res.AuthorId)
	assert.Equal(b, false, res.Published)

	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i)
		res, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{
			AdId:   res.Id,
			UserId: res.AuthorId,
			Title:  "good bye" + s,
			Text:   res.Text})
		assert.NoError(b, err, "client.UpdateAd")

		assert.Equal(b, "good bye"+s, res.Title)
		assert.Equal(b, "world", res.Text)
		assert.Equal(b, int64(111), res.AuthorId)
		assert.Equal(b, false, res.Published)
	}

}

func BenchmarkGRPCListAds(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 111})
	assert.NoError(b, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.AuthorId,
		Published: true})
	assert.NoError(b, err, "client.ChangeAdStatus")

	res, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "test", Text: "pupupu", UserId: 111})
	assert.NoError(b, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.AuthorId,
		Published: true})
	assert.NoError(b, err, "client.ChangeAdStatus")

	res, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "some title", Text: "some text", UserId: 115})
	assert.NoError(b, err, "client.CreateAd")
	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.AuthorId,
		Published: true})
	assert.NoError(b, err, "client.ChangeAdStatus")

	for i := 0; i < b.N; i++ {
		ads, err := client.ListAds(ctx, &emptypb.Empty{})
		assert.NoError(b, err, "client.ListAds")

		assert.Len(b, ads.List, 3)
	}

}

func BenchmarkGRPCGetUser(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Id: 15, Email: "email@exmple.com"})
	assert.NoError(b, err, "client.CreateUser")

	assert.Equal(b, "Oleg", res.Name)

	for i := 0; i < b.N; i++ {
		res, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 15})
		assert.NoError(b, err, "client.GetUser")

		assert.Equal(b, "Oleg", res.Name)
		assert.Equal(b, int64(15), res.Id)
	}

}

func BenchmarkGRPCDeleteUser(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)

	for i := 0; i < b.N; i++ {
		_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{
			Name:  "Oleg",
			Id:    int64(i),
			Email: "email@exmple.com"})
		assert.NoError(b, err, "client.CreateUser")
	}

	for i := 0; i < b.N; i++ {
		_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: int64(i)})
		assert.NoError(b, err, "client.DeleteUser")
	}

}

func BenchmarkGRPCDeleteAd(b *testing.B) {
	lis := bufconn.Listen(1024 * 1024)
	b.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	b.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewService(app.NewApp(adrepo.New(), usersrepo.New()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(b, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	b.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(b, err, "grpc.DialContext")

	b.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{
		Title:  "hello",
		Text:   "world",
		UserId: int64(111)})
	assert.NoError(b, err, "client.CreateAd")

	assert.Equal(b, "hello", res.Title)
	assert.Equal(b, "world", res.Text)
	assert.Equal(b, int64(111), res.AuthorId)

	for i := 1; i < b.N; i++ {

		_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: int64(i), AuthorId: res.AuthorId})
		assert.Error(b, err, "client.DeleteAd")
	}
}
