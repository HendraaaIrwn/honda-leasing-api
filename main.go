package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/HendraaaIrwn/honda-leasing-api/api/routers"
	configs "github.com/HendraaaIrwn/honda-leasing-api/internal/config"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/handler"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/repository"
	"github.com/HendraaaIrwn/honda-leasing-api/internal/services"
	"github.com/HendraaaIrwn/honda-leasing-api/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	gin.SetMode(gin.DebugMode)
	if cfg.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	if err := engine.SetTrustedProxies(cfg.Server.TrustedProxy); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	repos := repository.NewRepositoriesFromDatabase(db)
	svcs := services.NewServices(repos)
	handlers := handler.NewHandlers(svcs)

	routers.RegisterERDRouters(engine, cfg.Server.BasePath, handlers)

	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	log.Printf("Server listening on %s%s", cfg.Server.Address, cfg.Server.BasePath)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting server: %v", err)
	}
}
