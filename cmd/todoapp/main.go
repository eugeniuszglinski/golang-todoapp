package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/eugeniuszglinski/golang-todoapp/internal/core/config"
	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	"github.com/eugeniuszglinski/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/server"
	statistics_postgres "github.com/eugeniuszglinski/golang-todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/eugeniuszglinski/golang-todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/eugeniuszglinski/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres "github.com/eugeniuszglinski/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/eugeniuszglinski/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/eugeniuszglinski/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres "github.com/eugeniuszglinski/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/eugeniuszglinski/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/eugeniuszglinski/golang-todoapp/internal/features/users/transport/http"
	web_file_system "github.com/eugeniuszglinski/golang-todoapp/internal/features/web/repository/file_system"
	web_service "github.com/eugeniuszglinski/golang-todoapp/internal/features/web/service"
	web_transport_http "github.com/eugeniuszglinski/golang-todoapp/internal/features/web/transport/http"
	"go.uber.org/zap"

	_ "github.com/eugeniuszglinski/golang-todoapp/docs"
)

// @title       Golang Todo API
// @version     1.0
// @description Todo Application REST-API scheme
// @host        127.0.0.1:5050
// @BasePath    /api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("time_zone", time.Local))

	logger.Debug("Initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to initialize postgres connection pool: %v\n", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("Initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHttpHandler(usersService)

	logger.Debug("Initializing feature", zap.String("feature", "tasks"))

	tasksRepository := tasks_postgres.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHttp := tasks_transport_http.NewTasksHttpHandler(tasksService)

	logger.Debug("Initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHttp := statistics_transport_http.NewStatisticsHttpHandler(statisticsService)

	logger.Debug("Initializing feature", zap.String("feature", "web"))

	webRepository := web_file_system.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHttp := web_transport_http.NewWebHttpHandler(webService)

	logger.Debug("Initializing HTTP server")

	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHttpServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.PanicRecovery(),
	)

	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHttp.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHttp.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHttp.Routes()...)

	httpServer.RegisterApiRouters(apiVersionRouterV1)
	httpServer.RegisterRoutes(webTransportHttp.Routes()...)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
		os.Exit(1)
	}
}
