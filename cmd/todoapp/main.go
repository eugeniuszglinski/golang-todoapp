package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_postgres_pool "github.com/eugeniuszglinski/golang-todoapp/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/server"
	users_postgres "github.com/eugeniuszglinski/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/eugeniuszglinski/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/eugeniuszglinski/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Initializing postgres connection pool")

	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to initialize postgres connection pool: %v\n", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("Initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	usersTransportHttp := users_transport_http.NewUsersHttpHandler(usersService)

	logger.Debug("Initializing HTTP server")

	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.PanicRecovery(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHttp.Routes()...)

	httpServer.RegisterApiRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
		os.Exit(1)
	}
}
