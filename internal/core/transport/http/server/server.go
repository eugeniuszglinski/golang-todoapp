package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

// HttpServer is a simple HTTP server that can be started and stopped gracefully.
type HttpServer struct {
	mux         *http.ServeMux
	config      Config
	logger      *core_logger.Logger
	middlewares []core_http_middleware.Middleware
}

func NewHttpServer(
	config Config,
	logger *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HttpServer {
	return &HttpServer{mux: http.NewServeMux(), config: config, logger: logger, middlewares: middleware}
}

func (s *HttpServer) RegisterApiRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := fmt.Sprintf("/api/%s", router.apiVersion)

		// In global multiplexer, we register a handler for the prefix of the API version (e.g., /api/v1/).
		// This handler uses http.StripPrefix to remove the API version prefix from the incoming request's URL path
		// before passing it to the ApiVersionRouter. This way, the ApiVersionRouter can define its routes without
		// needing to include the API version in their paths. This is idiomatic because *feature* routers should not be
		// aware of the API versioning scheme, allowing them to focus solely on their specific routes and handlers.
		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddleware()))
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middlewares...)

	server := &http.Server{Handler: mux, Addr: s.config.Addr}

	// buffer 1 to not block
	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.logger.Warn("starting http server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.logger.Info("shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.logger.Warn("HTTP server shutdown completed")
	}

	return nil
}
