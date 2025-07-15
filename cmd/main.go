package main

import (
	"duking/internal/config"
	"duking/internal/delivery/http"
	"duking/internal/logger"
	"duking/internal/repository"
	"duking/internal/usecase"
	"duking/pkg/db"
	"log" // Используем только для фатальных ошибок инициализации логгера

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()

	// Инициализация логгера в зависимости от режима (prod/dev)
	lg, err := logger.Init(cfg.IsProd) // Используем cfg.IsProd
	if err != nil {
		log.Fatalf("failed to init logger: %v", err) // Используем стандартный log для этой критической ошибки
	}
	defer func() {
		if err := lg.Sync(); err != nil {
			// Обработка ошибки синхронизации логгера, если она возникнет
			// Например, при закрытии приложения
		}
	}()

	pool, err := db.InitDB(cfg)
	if err != nil {
		lg.Fatal("failed to initialize DB", zap.Error(err)) // Используем Zap логгер
		return                                              // Добавляем return, так как Fatal не всегда вызывает os.Exit
	}
	defer pool.Close()

	// Инициализация репозиториев
	hotelRepo := repository.NewHotelRepository(pool) // Переименовал для ясности
	userRepo := repository.NewUserRepository(pool)

	// Инициализация сервисов (use cases)
	hotelService := usecase.NewHotelService(hotelRepo)    // Переименовал для ясности
	userService := usecase.NewUserService(userRepo, *cfg) // Передаем весь конфиг, так как SecretKey нужен для JWT

	// Инициализация HTTP-хендлеров, передавая логгер
	hotelHandler := http.NewHotelHandler(hotelService, lg) // Переименовал и добавил логгер
	userHandler := http.NewUserHandler(userService, lg)    // Добавил и передал логгер

	r := gin.New()
	// Добавляем Gin Recovery middleware для восстановления после паник
	r.Use(gin.Recovery())
	// Добавляем Gin Logger middleware для логирования HTTP-запросов (опционально, Zap уже логирует)
	// r.Use(gin.Logger()) // Можно использовать, но Zap будет более детальным

	// Группировка маршрутов для отелей
	hotelRoutes := r.Group("/hotels")
	{
		hotelRoutes.POST("/create", hotelHandler.HotelCreate)
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
		// TODO: Добавить маршруты для GetProfile и UpdateProfile после их реализации
		// userRoutes.GET("/profile/:id", userHandler.GetProfile)
		// userRoutes.PATCH("/profile/:id", userHandler.UpdateProfile)
	}

	lg.Info("Server starting", zap.String("port", cfg.Port)) // Более информативное сообщение
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		lg.Fatal("failed to start server", zap.Error(err)) // Более информативное сообщение
	}
}
