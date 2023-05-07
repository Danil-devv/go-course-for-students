package httpgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
	"net/http"
)

type Server struct {
	port string
	app  *http.Server
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	s := Server{port: port, app: &http.Server{Addr: port, Handler: handler}}

	AppRouter(handler, a)

	return s
}

func (s *Server) Listen() error {
	return s.app.ListenAndServe()
}

func (s *Server) GracefulShutdown(ctx context.Context) error {
	return s.app.Shutdown(ctx)
}

func (s *Server) Handler() http.Handler {
	return s.app.Handler
}
