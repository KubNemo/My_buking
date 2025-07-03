package main

import (
	"duking/internal/config"
	"duking/internal/delivery/http"
	"duking/internal/logger"
	"duking/internal/repository"
	"duking/internal/usecase"
	"duking/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		lg.Error("failed to initialize DB", zap.Error(err))
	}
	defer pool.Close()
	repo := repository.NewRepository(pool)
	svc := usecase.NewService(repo)
	h := http.Newhandler(svc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/create", h.HotelCreate)
	r.GET("/oneHotel/:id", h.HotelGetOne)
	r.GET("/allHotel", h.HotelGetAll)
	r.PATCH("updateHotel/:id", h.UpdateHotel)
	r.DELETE("DeleteHotel/:id", h.HotelDelete)
	lg.Info("Server starting")
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		lg.Fatal("", zap.Error(err))
		return
	}
}
