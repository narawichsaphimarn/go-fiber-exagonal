package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	adapters "example.com/practice/fiber/internal/adapters/auth/jwt"
	http "example.com/practice/fiber/internal/adapters/http/handlers"
	"example.com/practice/fiber/internal/adapters/http/middleware"
	"example.com/practice/fiber/internal/adapters/repo"
	usecaseBook "example.com/practice/fiber/internal/usecase/book"
	usecaseUser "example.com/practice/fiber/internal/usecase/user"
	"example.com/practice/fiber/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := pkg.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	cfgPool, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName, cfg.DB.Options))
	if err != nil {
		log.Fatal(err)
	}
	cfgPool.MaxConns = 10
	cfgPool.MinConns = 1
	cfgPool.MaxConnLifetime = time.Hour
	pool, err := pgxpool.NewWithConfig(context.Background(), cfgPool)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	appV1 := app.Group("/v1")
	authProvider := adapters.NewProvider([]byte(cfg.Auth.JWT.Secret), 15*time.Minute)
	appProtectV1 := app.Group("/v1/auth", middleware.NewAuthMiddleware(authProvider).Protect)
	{
		bookRepo := repo.NewBookRepo(pool)
		bookService := usecaseBook.NewBookService(bookRepo)
		http.NewBookHandler(bookService).RegisterRoutes(appProtectV1)
	}

	{
		userRepo := repo.NewUserRepo(pool)
		userService := usecaseUser.NewUserService(userRepo, authProvider)
		http.NewUserHandler(userService).RegisterNotProtected(appV1)
		http.NewUserHandler(userService).RegisterProtected(appProtectV1)
	}
	log.Fatal(app.Listen(":" + strconv.Itoa(cfg.App.Port)))
}
