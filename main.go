package main

import (
	"animal-crud-api/controllers"
	"animal-crud-api/database"
	"animal-crud-api/middlewares"
	"animal-crud-api/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Login(c *gin.Context) {
    userID := uint(1) 
    
    token, err := utils.GenerateToken(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token, "message": "Login successful. Use this token for protected routes."})
}

func main() {
	if err := godotenv.Load(); err != nil { 
		log.Println("Error loading .env file. Using system environment variables or defaults.") 
	}

    database.Connect()

	r := gin.Default()

	api := r.Group("/api")
	api.POST("/login", Login)

	protected := api.Group("/animals") 
    protected.Use(middlewares.AuthRequired())
    {
		protected.POST("", controllers.CreateAnimal)
        protected.GET("", controllers.FindAnimals)
        protected.GET("/:id", controllers.FindAnimal)
        protected.PUT("/:id", controllers.UpdateAnimal)
        protected.DELETE("/:id", controllers.DeleteAnimal)
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