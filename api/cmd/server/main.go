package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/prachotx/real-time-chat/api/config"
	"github.com/prachotx/real-time-chat/api/internal/handler"
	"github.com/prachotx/real-time-chat/api/internal/middleware"
	"github.com/prachotx/real-time-chat/api/internal/model"
	"github.com/prachotx/real-time-chat/api/internal/repository"
	"github.com/prachotx/real-time-chat/api/internal/service"
	internalws "github.com/prachotx/real-time-chat/api/internal/ws"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{}, &model.Room{}, &model.Message{})

	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)
	messageHandler := handler.NewMessageHandler(messageService)

	hub := internalws.NewHub()
	go hub.Run()
	wsHandler := handler.NewWsHandler(hub, messageService)

	api := app.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.Post("/login", authHandler.Login)
			auth.Post("/register", authHandler.Register)
			auth.Get("/me", middleware.AuthMiddleware, authHandler.Profile)
		}
		room := api.Group("/rooms")
		{
			room.Post("/", middleware.AuthMiddleware, roomHandler.Create)
			room.Get("/", middleware.AuthMiddleware, roomHandler.FindAll)
			room.Post("/:id/messages", middleware.AuthMiddleware, messageHandler.Create)
			room.Get("/:id/messages", middleware.AuthMiddleware, messageHandler.FindByRoomID)
			room.Get("/:id/online-users", middleware.AuthMiddleware, wsHandler.GetOnlineUsers)
		}
	}

	app.Get("/ws/:room_id", middleware.AuthMiddleware, wsHandler.Handle)

	// เริ่ม server ใน goroutine แยก เพื่อให้ main goroutine รอ signal ได้
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	// รอ SIGINT (Ctrl+C) หรือ SIGTERM (docker stop / k8s)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")

	// หยุดรับ request ใหม่ และรอให้ request ที่ค้างอยู่เสร็จก่อน (max 5 วิ)
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	// ปิด DB connection pool
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("server exited cleanly")
}
