package config

import (
	"time"

	"github.com/codepnw/gopher-social/internal/utils/env"
)

type Config struct {
	App  AppConfig
	DB   DBConfig
	Mail MailConfig
}

type AppConfig struct {
	Addr        string
	AppVersion  string
	ApiURL      string
	ApiVersion  string
	Env         string
	FrontendURL string
}

type MailConfig struct {
	Exp       time.Duration
	ApiKey    string
	FromEmail string
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func InitConfig() Config {
	app := AppConfig{
		Addr:        env.GetString("APP_ADDR", ":8080"),
		ApiURL:      env.GetString("APP_API_URL", "localhost:8080"),
		ApiVersion:  env.GetString("APP_API_VERSION", "v1"),
		AppVersion:  env.GetString("APP_VERSION", "0.0.1"),
		Env:         env.GetString("APP_ENV", "development"),
		FrontendURL: env.GetString("APP_FRONTEND_URL", "http://localhost:3000"),
	}

	db := DBConfig{
		Addr:         env.GetString("DB_ADDR", ""),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	mail := MailConfig{
		Exp:       time.Hour * 24 * 3, // 3 days,
		ApiKey:    env.GetString("MAILTRAP_API_KEY", ""),
		FromEmail: env.GetString("FROM_EMAIL", ""),
	}

	return Config{
		App:  app,
		DB:   db,
		Mail: mail,
	}
}
