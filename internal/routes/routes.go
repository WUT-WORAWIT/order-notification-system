package routes

import (
	"order-notification-system/internal/api"
	"order-notification-system/internal/handlers"
	"order-notification-system/internal/middleware"
	"order-notification-system/internal/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configures the application routes.
// It takes the Gin engine and necessary handler instances as arguments.
func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// --- Initialize API handlers with DB dependency ---
	// These handlers are now created in main.go and passed here,
	// or they can be created here if they only need `db`.
	// For consistency with current main.go, let's assume they are created in main
	// and we'd pass them as arguments.
	// However, to keep this function self-contained for route setup,
	// we can also initialize them here if they only depend on `db`.

	orderAPIHandler := api.NewOrderAPI(db)
	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	profileHandler := handlers.NewProfileHandler(db)

	// Public routes
	// Grouping public routes under /api prefix
	publicAPIRoutes := r.Group("/api")
	{
		publicAPIRoutes.GET("/", userHandler.GetRemark)
		publicAPIRoutes.POST("/users", userHandler.CreateUser)
		publicAPIRoutes.POST("/login", authHandler.Login)
	}
	// Order related public route
	r.POST("/order", orderAPIHandler.CreateOrder)
	r.GET("/products", orderAPIHandler.GetProducts)
	// Protected routes
	// Grouping protected routes under /api prefix and applying JWT middleware
	protectedAPIRoutes := r.Group("/api")
	protectedAPIRoutes.Use(middleware.JWTMiddleware())
	{
		// protectedAPIRoutes.GET("/users", userHandler.GetUsersAll) // Assuming GetUsersAll exists
		protectedAPIRoutes.GET("/users/:username", userHandler.GetUserByID)
		protectedAPIRoutes.PUT("/users/:username", userHandler.UpdateUser)
		protectedAPIRoutes.DELETE("/users/:username", userHandler.DeleteUser)
		protectedAPIRoutes.GET("/profile", profileHandler.GetProfile)

		// Product routes (protected)
		// protectedAPIRoutes.GET("/products", orderAPIHandler.GetProducts)       // New route for getting all products
		protectedAPIRoutes.POST("/getproduct", orderAPIHandler.GetProduct)     // Existing route, kept for consistency if needed, but GET /products/:id is more RESTful
		protectedAPIRoutes.POST("/editproduct", orderAPIHandler.CreateProduct) // Existing route, consider changing to POST /products for creation
	}

	// WebSocket and Order Status routes (protected)
	r.GET("/ws", middleware.JWTMiddleware(), gin.WrapF(websocket.HandleWebSocket))
	r.PATCH("/orders/:id/status", middleware.JWTMiddleware(), orderAPIHandler.UpdateOrderStatus)
}
