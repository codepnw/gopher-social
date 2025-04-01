package main

import (
	"log"

	"github.com/codepnw/gopher-social/cmd/router"
	"github.com/codepnw/gopher-social/internal/env"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("dev.env"); err != nil {
		log.Fatal("failed loading env file")
	}

	app := &router.Application{
		Config: router.Config{
			Addr:    env.GetString("APP_ADDR", ":8080"),
			Version: env.GetString("APP_VERSION", "v1"),
		},
	}

	log.Fatal(app.Run(app.Routes()))
}
