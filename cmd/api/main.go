package main

import (
	"log"

	"github.com/codepnw/gopher-social/cmd/router"
	"github.com/codepnw/gopher-social/internal/database"
	"github.com/codepnw/gopher-social/internal/env"
	"github.com/codepnw/gopher-social/internal/store"
	"github.com/joho/godotenv"
)

const envPath = "dev.env"

func main() {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("failed loading env file")
	}

	cfg := router.Config{
		Addr:    env.GetString("APP_ADDR", ":8080"),
		Version: env.GetString("APP_VERSION", "v1"),
	}

	dbConfig := router.DBConfig{
		Addr:         env.GetString("DB_ADDR", ""),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	// Database
	db, err := database.NewDatabase(dbConfig.Addr, dbConfig.MaxOpenConns, dbConfig.MaxIdleConns, dbConfig.MaxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("database connected...")

	// Storage
	store := store.NewStorage(db)

	app := &router.Application{
		Config:   cfg,
		Store:    store,
		DBConfig: dbConfig,
	}

	log.Fatal(app.Run(app.Routes()))
}
