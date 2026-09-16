package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"pulse/handler"
	"pulse/repository"
	"pulse/service"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatal("error loading .env.dev")
	}

	dsn := os.Getenv("DATABASE_URL")
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

	userService := service.NewUserService(userRepo)
	followService := service.NewFollowService(followRepo, userRepo, notificationRepo, db)

	userHandler := handler.NewUserHandler(userService)
	followHandler := handler.NewFollowHandler(followService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.POST("/register", userHandler.RegisterUser)
	r.POST("/login", userHandler.LoginUser)

	auth := r.Group("/")
	auth.Use(handler.AuthMiddleware)
	auth.GET("/me", userHandler.Me)
	auth.POST("/users/:id/follow", followHandler.FollowUser)
	auth.DELETE("/users/:id/follow", followHandler.UnfollowUser)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	}
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
