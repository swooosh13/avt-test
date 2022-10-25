package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/config"
	"github.com/swooosh13/avito-test/internal/controller"
	"github.com/swooosh13/avito-test/internal/domain"
	"github.com/swooosh13/avito-test/internal/repository"
	s32 "github.com/swooosh13/avito-test/internal/repository/s3"
	"github.com/swooosh13/avito-test/pkg/postgres"
	"github.com/swooosh13/avito-test/pkg/server"
)

type App struct {
	logger     *zerolog.Logger
	httpServer *server.Server
}

func New(ctx context.Context, logger *zerolog.Logger) (*App, error) {
	cfg := config.Get()
	dsn := cfg.Postgres.DSN
	db, err := postgres.NewClient(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	s3 := s32.New(cfg.S3.AccessKey, cfg.S3.SecKey, cfg.S3.Endpoint, cfg.S3.BucketName)

	repos := repository.NewRepositories(db, logger, s3)
	services := domain.NewServices(repos, logger)

	router := controller.NewRouter(services, logger)

	fmt.Println(cfg)
	return &App{
		logger: logger,
		httpServer: server.New(
			router,
			server.WithHost(cfg.HTTP.Host),
			server.WithPort(cfg.HTTP.Port),
			server.WithMaxHeaderBytes(cfg.HTTP.MaxHeaderBytes),
			server.WithReadTimeout(cfg.HTTP.ReadTimeout),
			server.WithWriteTimeout(cfg.HTTP.WriteTimeout),
		),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	eChan := make(chan error)
	interrupt := make(chan os.Signal, 1)

	a.logger.Info().Msg("starting http server")

	go func() {
		if err := a.httpServer.Start(); err != nil {
			eChan <- fmt.Errorf("listen and serve: %w", err)
		}
	}()

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-eChan:
		return fmt.Errorf("start http server: %w", err)
	case <-interrupt:
	}

	const httpShutdownTimeout = 5 * time.Second
	if err := a.httpServer.Stop(ctx, httpShutdownTimeout); err != nil {
		return fmt.Errorf("stop http server: %w", err)
	}

	return nil
}
