package api

import (
	"fmt"
	"log"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/hamiddarani/web-api-fiber/docs"
	"github.com/hamiddarani/web-api-fiber/internal/api/middlewares"
	"github.com/hamiddarani/web-api-fiber/internal/api/routes"
	"github.com/hamiddarani/web-api-fiber/internal/api/validations"
	"github.com/hamiddarani/web-api-fiber/internal/config"
	"github.com/hamiddarani/web-api-fiber/pkg/logging"
)

type Server struct {
	logger logging.Logger
	cfg    *config.Config
}

func NewServer(cfg *config.Config, lg logging.Logger) *Server {
	return &Server{
		logger: lg,
		cfg:    cfg,
	}
}

func (s *Server) InitServer() {
	r := fiber.New()

	RegisterValidators()

	r.Use(recover.New())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
		AllowMethods:     "POST, GET, OPTIONS, PUT, DELETE, UPDATE",
		MaxAge:           21600,
	}))

	r.Use(limiter.New(limiter.Config{
		Max:               10,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	RegisterRoutes(r, s.cfg)
	RegisterSwagger(r)

	r.Listen(fmt.Sprintf(":%s", "5000"))

}

func RegisterRoutes(r *fiber.App, cfg *config.Config) {
	v1 := r.Group("/api/v1")
	{
		health := v1.Group("/health")
		routes.Health(health)

		auth := v1.Group("/auth")
		routes.AuthRoutes(auth, cfg)

		files := v1.Group("/files", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}))
		routes.FilesRoutes(files, cfg)

		posts := v1.Group("/posts", middlewares.Authentication(cfg), middlewares.Authorization([]string{"admin"}))
		routes.PostRoutes(posts, cfg)
	}

}

func RegisterValidators() {
	validate := validations.NewValidator()
	err := validate.RegisterValidation("mobile", validations.IranianMobileNumberValidator)
	if err != nil {
		log.Print(err.Error())
	}
}

func RegisterSwagger(r *fiber.App) {
	docs.SwaggerInfo.Title = "golang web api"
	docs.SwaggerInfo.Description = "golang web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", "5000")
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.Get("/swagger/*", swagger.HandlerDefault) // default
}
