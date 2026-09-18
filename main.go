package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"pulse/email"
	"pulse/handler"
	"pulse/repository"
	"pulse/service"
	"pulse/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type NotifyMsg struct {
	UserID int64  `json:"userId"`
	Body   string `json:"body"`
}

func main() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatal("error loading .env.dev")
	}

	dsn := os.Getenv("DATABASE_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5434/pulse?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	notificationRepo := repository.NewNotificationRepository(db)
	userRepo := repository.NewUserRepository(db)
	followRepo := repository.NewFollowRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)

	hub := ws.NewWSConnection()
	wsHandler := ws.NewHandler(hub)

	mailer, err := email.NewSenderFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	userService := service.NewUserService(userRepo)
	followService := service.NewFollowService(followRepo, userRepo, notificationRepo, outboxRepo, db)
	notificationService := service.NewNotificationService(notificationRepo)

	userHandler := handler.NewUserHandler(userService)
	followHandler := handler.NewFollowHandler(followService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.POST("/register", userHandler.RegisterUser)
	r.POST("/login", userHandler.LoginUser)

	auth := r.Group("/")
	auth.Use(handler.AuthMiddleware)
	auth.GET("/me", userHandler.Me)
	auth.GET("/notifications", notificationHandler.GetNotifications)
	auth.POST("/users/:id/follow", followHandler.FollowUser)
	auth.DELETE("/users/:id/follow", followHandler.UnfollowUser)
	auth.GET("/ws", wsHandler.WSHandler)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	}

	//Redis
	const channel = "pulse:notify"

	go func() {
		ctx := context.Background()
		sub := rdb.Subscribe(ctx, channel)
		ch := sub.Channel()
		for msg := range ch {
			var actualMsg NotifyMsg
			err := json.Unmarshal([]byte(msg.Payload), &actualMsg)
			if err != nil {
				log.Print(err)
				continue
			}
			hub.Send(actualMsg.UserID, websocket.TextMessage, []byte(actualMsg.Body))

		}
	}()

	go Run(outboxRepo, userRepo, mailer, rdb)
	go RunSweeper(outboxRepo)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}

}
