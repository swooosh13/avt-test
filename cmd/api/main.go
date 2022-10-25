package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/app"
)

var version string

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	logger.Info().Msg(fmt.Sprintf("starting api version %s", version))

	app, err := app.New(ctx, &logger)
	if err != nil {
		return fmt.Errorf("create api: %w", err)
	}

	return app.Run(ctx)
}
