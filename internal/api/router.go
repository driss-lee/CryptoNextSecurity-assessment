package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router sets up the HTTP router with all routes and middleware
type Router struct {
	handler *Handler
}

// NewRouter creates a new router instance
func NewRouter(handler *Handler, logger interface{}) *Router {
	return &Router{
		handler: handler,
	}
}

// Setup configures the router with all routes and middleware
func (r *Router) Setup() *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// API routes
	api := router.Group("/api/v1")
	{
		// Packet routes
		packets := api.Group("/packets")
		{
			packets.GET("", r.handler.GetPackets)
			packets.GET(":id", r.handler.GetPacketByID)
			packets.DELETE(":id", r.handler.DeletePacketByID)
			packets.DELETE("", r.handler.ClearPackets)
		}

		// Sniffing control routes
		sniffing := api.Group("/sniffing")
		{
			sniffing.POST("/start", r.handler.StartSniffing)
			sniffing.POST("/stop", r.handler.StopSniffing)
			sniffing.GET("/status", r.handler.SniffingStatus)
		}

		// Health and stats
		api.GET("/health", r.handler.Health)
		api.GET("/stats", r.handler.Stats)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
