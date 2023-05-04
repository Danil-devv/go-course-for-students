package httpgin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"

	"homework9/internal/app"
	"log"
	"net/http"
	"time"
)

type Server struct {
	port string
	app  *http.Server
}

func NewHTTPServer(port string, a *app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	s := Server{port: port, app: &http.Server{Addr: port, Handler: handler}}

	AppRouter(handler, *a)

	return s
}

func (s *Server) Listen() {
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

	eg.Go(func() error {
		log.Printf("starting http server, listening on %s\n", s.app.Addr)
		defer log.Printf("close http server listening on %s\n", s.app.Addr)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := s.app.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", s.app.Addr, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := s.app.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
		log.Printf("gracefully shutting down the http server: %s\n", err.Error())
	}

	log.Println("http server were successfully shutdown")
}

func (s *Server) Handler() http.Handler {
	return s.app.Handler
}
