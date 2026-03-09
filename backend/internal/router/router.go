package router

import (
	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/config"
	"github.com/godongmin/open_talk/backend/internal/handler"
	"github.com/godongmin/open_talk/backend/internal/middleware"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"github.com/godongmin/open_talk/backend/internal/service"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// Repositories
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	friendRepo := repository.NewFriendRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, friendRepo, chatRepo)
	chatService := service.NewChatService(chatRepo, messageRepo)
	friendService := service.NewFriendService(friendRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, authService)
	chatHandler := handler.NewChatHandler(chatService, messageRepo)
	friendHandler := handler.NewFriendHandler(friendService)

	// WebSocket
	hub := handler.NewHub(chatRepo, messageRepo)
	go hub.Run()
	wsHandler := handler.NewWSHandler(hub)

	// Public routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("/auth/logout", authHandler.Logout)

		users := protected.Group("/users")
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
			users.GET("/search", userHandler.SearchUsers)
			users.PUT("/password", userHandler.ChangePassword)
			users.DELETE("/me", userHandler.DeleteAccount)
		}

		chats := protected.Group("/chats")
		{
			chats.POST("", chatHandler.CreateChat)
			chats.GET("", chatHandler.GetChats)
			chats.GET("/:id/messages", chatHandler.GetMessages)
			chats.POST("/:id/messages", chatHandler.SendMessage)
			chats.PUT("/:id/read", chatHandler.MarkMessagesRead)
		}

		friends := protected.Group("/friends")
		{
			friends.POST("/request", friendHandler.SendRequest)
			friends.PUT("/:id/accept", friendHandler.AcceptRequest)
			friends.PUT("/:id/reject", friendHandler.RejectRequest)
			friends.GET("", friendHandler.GetFriends)
			friends.DELETE("/:id", friendHandler.BlockFriend)
		}
	}

	// WebSocket (protected)
	r.GET("/ws", middleware.AuthMiddleware(cfg), wsHandler.HandleWebSocket)

	return r
}
