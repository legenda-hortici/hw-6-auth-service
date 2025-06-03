package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Env      string `envconfig:"ENV" default:"local"`
	LogLevel string `envconfig:"LOG_LEVEL" required:"true"`
	Database DatabaseConfig
	GRPC     GRPCconfig
	TokenJWT TokenJWT
}

type DatabaseConfig struct {
	Host                 string        `envconfig:"DB_HOST" required:"true"`
	Port                 string        `envconfig:"DB_PORT" required:"true"`
	Username             string        `envconfig:"DB_USERNAME" required:"true"`
	Password             string        `envconfig:"DB_PASSWORD" required:"true"`
	DatabaseName         string        `envconfig:"DB_DATABASE" required:"true"`
	SSLMode              string        `envconfig:"DB_SSL_MODE" default:"disable" required:"true"`
	PoolMaxConns         int           `envconfig:"DB_POOL_MAX_CONNS" default:"10" required:"true"`
	PoolMaxConnsLifetime time.Duration `envconfig:"DB_POOL_MAX_CONN_LIFETIME" required:"true"`
	PoolMaxConnsIdletime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" required:"true"`
}

type GRPCconfig struct {
	ListenPort int           `envconfig:"GRPC_PORT" required:"true"`
	Timeout    time.Duration `envconfig:"GRPC_TIMEOUT" required:"true"`
}

type TokenJWT struct {
	Secret string        `envconfig:"JWT_SECRET" required:"true"`
	TTL    time.Duration `envconfig:"JWT_TTL" required:"true"`
}

func NewConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		panic("error loading .env file")
	}

	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		panic(err)
	}

	return &conf
}
