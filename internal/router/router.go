package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/aryadisastra/main-app-be/internal/handlers"
	"github.com/aryadisastra/main-app-be/internal/httpx"
	"github.com/aryadisastra/main-app-be/internal/middleware"
)

func New(db *gorm.DB, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(httpx.RecoverJSON(), httpx.NotFoundAsJSON())

	r.GET("/health", func(c *gin.Context) { httpx.OK(c, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(jwtSecret))
	{
		api.POST("/shipments", handlers.CreateShipment(db))
		api.GET("/shipments", handlers.GetShipments(db))
		api.GET("/shipments/track/:trackingNumber", handlers.TrackShipment(db))

		admin := api.Group("/shipments")
		admin.Use(middleware.RequireRoles("admin"))
		{
			admin.PATCH("/:trackingNumber/status", handlers.UpdateShipmentStatus(db))
		}
	}
	return r
}
