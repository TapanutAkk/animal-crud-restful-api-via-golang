package main

import (
	"animal-crud-api/controllers"
	"animal-crud-api/database"
	"animal-crud-api/middlewares"
	"animal-crud-api/models"
	"animal-crud-api/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
    var req LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
        return
    }

    var user models.User

    result := database.DB.Where("username = ?", req.Username).First(&user)

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    
    if err != nil {
        if err == bcrypt.ErrMismatchedHashAndPassword {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error during password verification"})
        return
    }

	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token": accessToken, 
        "refresh_token": refreshToken, 
        "message": "Login successful.",
    })
}

func RefreshToken(c *gin.Context) {
    var request struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
        return
    }

    claims, err := utils.ValidateToken(request.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
        return
    }

    const requiredAudience = "refresh"
    found := false
    
    for _, aud := range claims.Audience {
        if aud == requiredAudience {
            found = true
            break
        }
    }
    
    if !found {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not a valid refresh token (Audience mismatch)"})
        return
    }
    
    newAccessToken, newRefreshToken, err := utils.GenerateTokenPair(claims.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token": newAccessToken, 
        "refresh_token": newRefreshToken, 
    })
}

func main() {
	if err := godotenv.Load(); err != nil { 
		log.Println("Error loading .env file. Using system environment variables or defaults.") 
	}

    database.Connect()

	r := gin.Default()

	api := r.Group("/api")
	api.POST("/login", Login)
	api.POST("/refresh", RefreshToken)

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