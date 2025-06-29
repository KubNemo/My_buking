package main

import (
	"duking/internal/config"
	"duking/internal/logger"
	"duking/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	lg, err := logger.Init(false)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer lg.Sync()

	pool, err := db.InitDB(cfg)
	if err != nil {
		// TODO ...
	}
	_ = pool

	r := gin.New()
	r.Use(gin.Recovery())
	if err := r.Run(":" + cfg.Port); err != nil {
		return // TODO
	}
}
