package main

import (
	hostes "duking/internal/Hostes"
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
	repo := hostes.NewRepository(pool)
	svc := hostes.NewService(repo)
	h := hostes.Newhandler(svc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/create", h.HotelCreate)
	r.GET("/oneHotel", h.HotelGetOne)
	r.GET("/allHotel", h.HotelGetAll)
	r.PATCH("updateHotel", h.UpdateHotel)
	r.DELETE("DeleteHotel", h.HotelDelete)
	if err := r.Run(":" + cfg.Port); err != nil {
		return // TODO
	}
}
