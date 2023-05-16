package main

import (
	"GDOservice/internal/app"
	"GDOservice/internal/config"
	"GDOservice/pkg/logging"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	log.Println("Config initializing")
	cfg := config.GetConfig()

	log.Println("Logger initializing")
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}
	a.Run()
}
