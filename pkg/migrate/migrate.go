package migrate

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func connect(path, dsn string) (*migrate.Migrate, error) {
	dsn += "&sslmode=disable"

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://"+path, dsn)
		if err == nil {
			break
		}

		log.Printf("migrate postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		return nil, fmt.Errorf("migrate postgres connect error: %w", err)
	}

	return m, nil
}

func Up(path, dsn string) error {
	m, err := connect(path, dsn)
	if err != nil {
		return err
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up error: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate no change")
		return nil
	}

	return nil
}

func Down(path, dsn string) error {
	m, err := connect(path, dsn)
	if err != nil {
		return err
	}

	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate down error: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate no change")
		return nil
	}

	log.Printf("migrate down success")
	return nil
}

func Steps(path, dsn string, steps int) (err error) {
	m, err := connect(path, dsn)
	if err != nil {
		return
	}

	err = m.Steps(steps)
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate steps error: %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate no change")
		return
	}
	if steps > 0 {
		log.Printf("migrate up success")
		return
	} else {
		log.Printf("migrate down success")
		return
	}
}
