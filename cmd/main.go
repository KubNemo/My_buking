package main

import (
	hostes "duking/internal/Hostes"
	"duking/internal/config"
	"duking/internal/logger"
	"duking/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	err := godotenv.Load(".env") // или "./internal/config/.env" — смотри по расположению
	if err != nil {
		log.Println("⚠️  .env файл не найден, загружаю переменные из окружения")
	}
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
	repo := hostes.NewRepository(pool)
	svc := hostes.NewService(repo)
	h := hostes.Newhandler(svc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/create", h.HotelCreate)
	r.GET("/oneHotel/:id", h.HotelGetOne)
	r.GET("/allHotel", h.HotelGetAll)
	r.PATCH("updateHotel/:id", h.UpdateHotel)
	r.DELETE("DeleteHotel/:id", h.HotelDelete)
	lg.Info("Server starting")
	port := cfg.Port
	// if port == "" {
	// 	port = "8080"
	// }
	if err := r.Run(":" + port); err != nil {
		lg.Fatal("", zap.Error(err))
		return
	}
}
