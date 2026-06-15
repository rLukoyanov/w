package apis

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rLukoyanov/w/apis/handlers"
	"github.com/rLukoyanov/w/apis/middlewares"
	"github.com/rLukoyanov/w/apis/ws"
	"github.com/rLukoyanov/w/storage"
	"github.com/rLukoyanov/w/store"
	"github.com/rLukoyanov/w/web"
)

type Server struct {
	app         *fiber.App
	store       store.Store
	port        int
	hub         *ws.Hub
	jwtSecret   string
	fileStorage storage.FileStorage
}

func NewServer(st store.Store, hub *ws.Hub, port int, jwtSecret string, fs storage.FileStorage) *Server {
	app := fiber.New(fiber.Config{
		AppName: "W",
	})

	app.Use(recover.New())
	app.Use(middlewares.Zerolog())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	s := &Server{
		app:         app,
		store:       st,
		port:        port,
		hub:         hub,
		jwtSecret:   jwtSecret,
		fileStorage: fs,
	}

	s.setupRoutes(jwtSecret)
	return s
}

func (s *Server) setupRoutes(jwtSecret string) {
	api := s.app.Group("/api")

	authHandler := handlers.NewAuthHandler(s.store, jwtSecret)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middlewares.JWTAuth(jwtSecret))

	protected.Get("/auth/me", authHandler.Me)

	serversHandler := handlers.NewServersHandler(s.store)
	protected.Post("/servers", serversHandler.Create)
	protected.Get("/servers", serversHandler.GetAll)
	protected.Get("/servers/:id", serversHandler.GetByID)
	protected.Get("/servers/:id/members", serversHandler.GetMembers)
	protected.Patch("/servers/:id", serversHandler.Update)
	protected.Delete("/servers/:id", serversHandler.Delete)

	channelsHandler := handlers.NewChannelsHandler(s.store)
	protected.Post("/servers/:id/channels", channelsHandler.Create)
	protected.Get("/servers/:id/channels", channelsHandler.GetByServerID)
	protected.Get("/channels/:id", channelsHandler.GetByID)
	protected.Delete("/channels/:id", channelsHandler.Delete)

	messagesHandler := handlers.NewMessagesHandler(s.store)
	protected.Get("/channels/:id/messages", messagesHandler.GetByChannelID)
	protected.Post("/channels/:id/messages", messagesHandler.Create)
	protected.Patch("/messages/:id", messagesHandler.Update)
	protected.Delete("/messages/:id", messagesHandler.Delete)

	usersHandler := handlers.NewUsersHandler(s.store, s.hub)
	protected.Get("/users/connected", usersHandler.ListConnected)

	invitesHandler := handlers.NewInvitesHandler(s.store)
	protected.Post("/servers/:id/invites", invitesHandler.Create)
	protected.Get("/servers/:id/invites", invitesHandler.ListByServer)
	protected.Delete("/invites/:id", invitesHandler.Delete)
	protected.Post("/invites/:code/join", invitesHandler.Join)

	filesHandler := handlers.NewFilesHandler(s.store, s.fileStorage)
	protected.Post("/channels/:id/attachments", filesHandler.Upload)
	protected.Get("/attachments/:id", filesHandler.Download)

	adminHandler := handlers.NewAdminHandler(s.store, s.hub)
	admin := protected.Group("/admin")
	admin.Use(middlewares.RequireAdmin(s.store))
	admin.Get("/stats", adminHandler.Stats)
	admin.Get("/users", adminHandler.ListUsers)
	admin.Patch("/users/:id/role", adminHandler.UpdateUserRole)
	admin.Delete("/users/:id", adminHandler.DeleteUser)
	admin.Get("/servers", adminHandler.ListServers)
	admin.Delete("/servers/:id", adminHandler.DeleteServer)

	protected.Get("/openapi.yaml", OpenAPISpec)
	protected.Get("/openapi.json", OpenAPISpecJSON)

	s.app.Get("/ws", websocket.New(s.handleWebSocket, websocket.Config{}))

	s.app.Get("/docs", SwaggerUI)

	distFS, err := fs.Sub(web.DistFS, "build")
	if err != nil {
		panic(fmt.Sprintf("failed to get dist subdirectory: %v", err))
	}

	s.app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(distFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			return strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/ws") || strings.HasPrefix(path, "/docs")
		},
	}))
}

func (s *Server) handleWebSocket(c *websocket.Conn) {
	token := c.Query("token")
	if token == "" {
		c.WriteJSON(fiber.Map{"error": "missing token"})
		c.Close()
		return
	}

	claims := &jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		c.WriteJSON(fiber.Map{"error": "invalid token"})
		c.Close()
		return
	}

	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		c.WriteJSON(fiber.Map{"error": "invalid user_id in token"})
		c.Close()
		return
	}

	client := ws.NewClient(userID, c, s.hub)

	s.hub.Register(client)

	go client.WritePump()
	client.ReadPump()
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
