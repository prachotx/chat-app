package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fasthttp/websocket"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/prachotx/real-time-chat/api/config"
	"github.com/prachotx/real-time-chat/api/internal/dto"
	"github.com/prachotx/real-time-chat/api/internal/handler"
	"github.com/prachotx/real-time-chat/api/internal/middleware"
	"github.com/prachotx/real-time-chat/api/internal/model"
	"github.com/prachotx/real-time-chat/api/internal/repository"
	"github.com/prachotx/real-time-chat/api/internal/service"
	hub "github.com/prachotx/real-time-chat/api/internal/ws"
	jwtpkg "github.com/prachotx/real-time-chat/api/pkg/jwt"
	"github.com/valyala/fasthttp"
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
		}
	}

	upgrader := websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true },
	}

	h := hub.New()

	app.Get("/ws/:roomId", func(c fiber.Ctx) error {
		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		token, err := jwtpkg.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		claims := token.Claims.(gojwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		roomId := c.Params("roomId")
		roomIdUint, err := strconv.ParseUint(roomId, 10, 64)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
			cl := &hub.Client{Conn: conn, UserID: userID}
			h.Register(roomId, cl)
			fmt.Printf("Connected room=%s userID=%d total=%d\n", roomId, userID, h.Count(roomId))

			defer func() {
				h.Unregister(roomId, cl)
				fmt.Printf("Disconnected room=%s userID=%d total=%d\n", roomId, userID, h.Count(roomId))
			}()

			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					break
				}
				var incoming hub.IncomingMessage
				if err := json.Unmarshal(msg, &incoming); err != nil || incoming.Content == "" {
					continue
				}
				saved, err := messageService.Save(dto.CreateMessageDto{Content: incoming.Content}, uint(roomIdUint), userID)
				if err != nil {
					fmt.Printf("save error: %v\n", err)
					continue
				}
				event, err := hub.NewMessageEvent(saved)
				if err != nil {
					continue
				}
				h.Broadcast(roomId, websocket.TextMessage, event)
			}
		})
	})

	app.Listen(":" + cfg.Port)
}
