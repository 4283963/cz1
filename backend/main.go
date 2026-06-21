package main

import (
	"log"

	"guppy-breeding/database"
	"guppy-breeding/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		tanks := api.Group("/tanks")
		{
			tanks.GET("", handlers.ListTanks)
			tanks.POST("", handlers.CreateTank)
		}

		water := api.Group("/water")
		{
			water.POST("", handlers.CreateWaterRecord)
			water.GET("", handlers.ListWaterRecords)
			water.GET("/tank/:tank_id/latest", handlers.GetLatestWaterRecord)
		}

		breeding := api.Group("/breeding")
		{
			breeding.POST("", handlers.CreateBreedingRecord)
			breeding.GET("", handlers.ListBreedingRecords)
			breeding.GET("/:id", handlers.GetBreedingRecord)
			breeding.PUT("/:id", handlers.UpdateBreedingRecord)
			breeding.DELETE("/:id", handlers.DeleteBreedingRecord)
		}
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
