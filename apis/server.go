package apis

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rLukoyanov/w/apis/handlers"
	"github.com/rLukoyanov/w/apis/middlewares"
	"github.com/rLukoyanov/w/apis/ws"
	"github.com/rLukoyanov/w/store"
	"github.com/rLukoyanov/w/web"
)

type Server struct {
	app       *fiber.App
	store     store.Store
	port      int
	hub       *ws.Hub
	jwtSecret string
}

func NewServer(st store.Store, hub *ws.Hub, port int, jwtSecret string) *Server {
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
		app:       app,
		store:     st,
		port:      port,
		hub:       hub,
		jwtSecret: jwtSecret,
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

	// WebSocket endpoint
	s.app.Get("/ws", websocket.New(s.handleWebSocket, websocket.Config{
		// Allow connections from any origin for development
		// TODO: Restrict in production
	}))

	// Serve embedded SvelteKit frontend
	distFS, err := fs.Sub(web.DistFS, "build")
	if err != nil {
		panic(fmt.Sprintf("failed to get dist subdirectory: %v", err))
	}

	s.app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(distFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html", // SPA fallback
	}))
}

func (s *Server) handleWebSocket(c *websocket.Conn) {
	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		c.WriteJSON(fiber.Map{"error": "missing token"})
		c.Close()
		return
	}

	// Validate JWT token
	claims := &jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		c.WriteJSON(fiber.Map{"error": "invalid token"})
		c.Close()
		return
	}

	// Extract user ID from claims
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		c.WriteJSON(fiber.Map{"error": "invalid user_id in token"})
		c.Close()
		return
	}

	// Create new client
	client := ws.NewClient(userID, c, s.hub)

	// Register client with hub
	s.hub.Register(client)

	// Start goroutines
	go client.WritePump()
	client.ReadPump() // Blocks until connection closes
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
