package server

import (
	"context"
	"log"
	"os"

	"github.com/YonkLongSchlong/Todo-BE/packages/auth"
	"github.com/YonkLongSchlong/Todo-BE/packages/todo"
	"github.com/YonkLongSchlong/Todo-BE/packages/user"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type ApiServer struct {
	addr string
	db   *sqlx.DB
}

func NewApiServer(addr string, db *sqlx.DB) *ApiServer {
	return &ApiServer{addr: addr, db: db}
}

func (s *ApiServer) Run() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Fatal(err)
	}
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	e := echo.New()

	authGroup := e.Group("/api/v1/auth")
	todoGroup := e.Group("/api/v1")
	userGroup := e.Group("/api/v1")

	todoGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
	}))

	userGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
	}))

	userStore := user.NewStore(s.db)
	authRoute := auth.NewRoute(userStore)
	authRoute.Routes(authGroup)

	todoStore := todo.NewStore(s.db)
	todoRoute := todo.NewRoute(todoStore)
	todoRoute.Routes(todoGroup)

	userRoute := user.NewRoute(userStore, uploader)
	userRoute.Routes(userGroup)

	e.Start(s.addr)
}
