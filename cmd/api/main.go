package main

import (
	"log"

	"github.com/codepnw/gopher-social/cmd/config"
	"github.com/codepnw/gopher-social/cmd/router"
	"github.com/codepnw/gopher-social/internal/database"
	"github.com/codepnw/gopher-social/internal/store"
	"github.com/codepnw/gopher-social/internal/store/cache"
	"github.com/codepnw/gopher-social/internal/utils/logger"
	"github.com/codepnw/gopher-social/internal/utils/mailer"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const envPath = "dev.env"

func main() {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("failed loading env file")
	}

	cfg := config.InitConfig()

	// Logger
	logger, err := logger.InitLogger()
	if err != nil {
		log.Panic(err)
	}
	defer logger.Sync()

	// Database
	db, err := database.NewDatabase(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	logger.Info("database connected...")

	// Cache
	var rdb *redis.Client
	if cfg.Redis.Enabled {
		rdb = cache.NewRedisClient(cfg.Redis.Addr, cfg.Redis.Pw, cfg.Redis.DB)
		logger.Info("redis cache connected....")
	}

	cacheStorage := cache.NewRedisStorage(rdb)

	// Mailer
	mailer, err := mailer.NewMailTrapClient(cfg.Mail.ApiKey, cfg.Mail.FromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Storage
	store := store.NewStorage(db, cfg, mailer, cacheStorage)

	app := &router.Application{
		Config: cfg,
		Store:  store,
		Logger: logger,
	}

	logger.Fatal(app.Run(app.Routes(cacheStorage)))
}
