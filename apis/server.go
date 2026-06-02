package apis

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rLukoyanov/w/apis/handlers"
	"github.com/rLukoyanov/w/apis/middlewares"
	"github.com/rLukoyanov/w/store"
)

type Server struct {
	app   *fiber.App
	store store.Store
	port  int
}

func NewServer(st store.Store, port int, jwtSecret string) *Server {
	app := fiber.New(fiber.Config{
		AppName: "W",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middlewares.Zerolog())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	s := &Server{
		app:   app,
		store: st,
		port:  port,
	}

	s.setupRoutes(jwtSecret)
	return s
}

func (s *Server) setupRoutes(jwtSecret string) {
	api := s.app.Group("/api")

	// Auth routes
	authHandler := handlers.NewAuthHandler(s.store, jwtSecret)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("")
	protected.Use(middlewares.JWTAuth(jwtSecret))

	protected.Get("/auth/me", authHandler.Me)

	// Servers
	serversHandler := handlers.NewServersHandler(s.store)
	protected.Post("/servers", serversHandler.Create)
	protected.Get("/servers/:id", serversHandler.GetByID)

	// Channels
	channelsHandler := handlers.NewChannelsHandler(s.store)
	protected.Post("/servers/:id/channels", channelsHandler.Create)
	protected.Get("/channels/:id", channelsHandler.GetByID)

	// Messages
	messagesHandler := handlers.NewMessagesHandler(s.store)
	protected.Get("/channels/:id/messages", messagesHandler.GetByChannelID)
	protected.Post("/channels/:id/messages", messagesHandler.Create)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
