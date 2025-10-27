package main

import (
	"animal-crud-api/controllers"
	"animal-crud-api/database"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("Error loading .env file, using default environment variables for port.")
    }

    database.Connect()

	r := gin.Default()

	api := r.Group("/api")
    {
		api.POST("/animals", controllers.CreateAnimal)
        api.GET("/animals", controllers.FindAnimals)
        api.GET("/animals/:id", controllers.FindAnimal)
        api.PUT("/animals/:id", controllers.UpdateAnimal)
        api.DELETE("/animals/:id", controllers.DeleteAnimal)
	}

    port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

    runAddress := fmt.Sprintf(":%s", port)
    
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(runAddress); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}