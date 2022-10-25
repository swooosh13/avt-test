package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Debug    bool
		HTTP     HTTP
		Postgres Postgres
		Logger   Logger
		S3       S3
	}

	HTTP struct {
		Host           string        `envconfig:"HTTP_HOST"             required:"true" default:"localhost"`
		Port           string        `envconfig:"HTTP_PORT"             required:"true" default:"8080"`
		MaxHeaderBytes int           `envconfig:"HTTP_MAX_HEADER_BYTES"                 default:"1"`
		ReadTimeout    time.Duration `envconfig:"HTTP_READ_TIMEOUT"                     default:"10s"`
		WriteTimeout   time.Duration `envconfig:"HTTP_WRITE_TIMEOUT"                    default:"10s"`
	}

	S3 struct {
		AccessKey  string `envconfig:"S3_CLOUD_KEY" required:"true"`
		SecKey     string `envconfig:"S3_CLOUD_SECRET" required:"true"`
		Endpoint   string `envconfig:"S3_CLOUD_ENDPOINT" required:"true"`
		BucketName string `envconfig:"S3_CLOUD_BUCKET_NAME" required:"true"`
	}

	Postgres struct {
		DSN string `envconfig:"POSTGRES_DSN" json:"-" default:""`
	}

	Logger struct {
		Level string `envconfig:"LOGGER_LEVEL" default:"info"`
	}
)

var (
	instance Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		if err := envconfig.Process("", &instance); err != nil {
			log.Fatal(err)
		}

		if instance.Debug {
			configBytes, err := json.MarshalIndent(instance, "", " ")
			if err != nil {
				log.Fatal(fmt.Errorf("error marshaling indent config: %w", err))
			}

			fmt.Println("Configuration:", string(configBytes))
		}
	})

	return &instance
}
