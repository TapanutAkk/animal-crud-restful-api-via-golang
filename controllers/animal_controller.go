package controllers

import (
	"animal-crud-api/database"
	"animal-crud-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateAnimal(c *gin.Context) {
	var animal models.Animal
	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&animal)
	c.JSON(http.StatusOK, animal)
}

func FindAnimals(c *gin.Context) {
	var animals []models.Animal
	database.DB.Find(&animals)

	c.JSON(http.StatusOK, animals)
}

func FindAnimal(c *gin.Context) {
    id := c.Param("id") 
    var animal models.Animal

    if err := database.DB.First(&animal, id).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Animal not found!"})
        return
    }

    c.JSON(http.StatusOK, animal)
}

func UpdateAnimal(c *gin.Context) {
    id := c.Param("id")
    var animal models.Animal
    
    if err := database.DB.First(&animal, id).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Animal not found!"})
        return
    }

    var input models.Animal
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    database.DB.Model(&animal).Updates(input)
    
    c.JSON(http.StatusOK, animal)
}

func DeleteAnimal(c *gin.Context) {
    id := c.Param("id")
    var animal models.Animal
    
    if err := database.DB.First(&animal, id).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Animal not found!"})
        return
    }

    database.DB.Delete(&animal)

    c.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}