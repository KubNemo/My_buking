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

	lg, err := logger.Init(cfg.IsProd)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func() {
		if err := lg.Sync(); err != nil {
			lg.Error("failed to sync logger", zap.Error(err))
		}
	}()

	pool, err := db.InitDB(cfg)
	if err != nil {
		lg.Fatal("failed to initialize DB", zap.Error(err))
	}
	defer pool.Close()

	// Инициализация репозиториев
	hotelRepo := repository.NewHotelRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	roomRepo := repository.NewRooomRepository(pool) // ДОБАВЛЕНО: Репозиторий комнат

	// Инициализация сервисов (use cases)
	hotelService := usecase.NewHotelsService(hotelRepo)
	userService := usecase.NewUserService(userRepo, *cfg, lg) // Передаем логгер в UserService
	roomService := usecase.NewRoomService(roomRepo)           // ДОБАВЛЕНО: Сервис комнат

	// Инициализация HTTP-хендлеров, передавая логгер
	hotelHandler := http.NewHotelHandler(hotelService, lg)
	userHandler := http.NewUserHandler(userService, lg)
	roomHandler := http.NewRoomHandler(roomService, lg) // ДОБАВЛЕНО: Хендлер комнат

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger()) // Можно использовать, но Zap будет более детальным

	// Группировка маршрутов для отелей
	hotelRoutes := r.Group("/hotels")
	{
		hotelRoutes.POST("/", hotelHandler.HotelCreate) // Обновил маршрут
		hotelRoutes.GET("/:id", hotelHandler.HotelGetOne)
		hotelRoutes.GET("/", hotelHandler.HotelGetAll)
		hotelRoutes.PATCH("/:id", hotelHandler.UpdateHotel)
		hotelRoutes.DELETE("/:id", hotelHandler.HotelDelete)
	}

	// Группировка маршрутов для пользователей
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", userHandler.CreateUser)
		userRoutes.POST("/login", userHandler.UserLogin)
		userRoutes.GET("/profile/:id", userHandler.GetProfile)
		userRoutes.PATCH("/profile/:id", userHandler.UpdateProfile)
	}

	// ДОБАВЛЕНО: Группировка маршрутов для комнат
	roomRoutes := r.Group("/rooms")
	{
		roomRoutes.POST("/", roomHandler.CreateRoom)
		roomRoutes.GET("/:id", roomHandler.GetRoom)
		roomRoutes.PATCH("/:id", roomHandler.UpdateRoom)
		roomRoutes.DELETE("/:id", roomHandler.DeleteRoom)
	}

	// ДОБАВЛЕНО: Маршрут для получения комнат по отелю
	r.GET("/hotels/:hotel_id/rooms", roomHandler.GetRoomsByHotel)

	lg.Info("Server starting", zap.String("port", cfg.Port))
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		lg.Fatal("failed to start server", zap.Error(err))
	}
}
