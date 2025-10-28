package database

import (
	"animal-crud-api/models"
	"animal-crud-api/utils"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    if err := godotenv.Load(); err != nil {
        log.Println("Error loading .env file, using default environment variables")
    }

    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
        dbHost, dbUser, dbPassword, dbName, dbPort,
    )
    
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatal("Failed to connect to database! \n", err)
    }

    fmt.Println("Connected to database successfully!")

    DB.AutoMigrate(&models.Animal{}, &models.User{})

    seedAnimals()
    seedUsers();
}

func seedAnimals() {
    var count int64
    DB.Model(&models.Animal{}).Count(&count)

    if count == 0 {
        numToSeed := 10
        var animals []models.Animal
        
        usedNames := make(map[string]bool) 

        for i := 0; i < numToSeed; i++ {
            var name string
            
            for {
                name = utils.RandomName()
                if !usedNames[name] {
                    usedNames[name] = true
                    break
                }
                if len(usedNames) == len(utils.AnimalNames) {
                    break 
                }
            }
            
            species := utils.RandomSpecies()
            age := utils.RandomAge()

            animal := models.Animal{
                Name:    name,
                Species: species,
                Age:     age,
            }
            animals = append(animals, animal)
        }
        
        for _, animal := range animals {
            if err := DB.Create(&animal).Error; err != nil {
                log.Printf("Could not seed animal %s: %v", animal.Name, err)
            }
        }

        fmt.Printf("Database seeded with %d initial unique Animal data entries.\n", len(animals))
    } else {
        fmt.Println("Animal table already has data, skipping seeding.")
    }
}

func seedUsers() {
    var count int64
    DB.Model(&models.User{}).Count(&count)

    if count == 0 {
        plainPassword := "password123" 

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
        if err != nil {
            log.Fatalf("Failed to hash password for seeding: %v", err)
            return
        }

        user := models.User{
            Username: "admin",
            Password: string(hashedPassword),
        }

        if err := DB.Create(&user).Error; err != nil {
            log.Printf("Could not seed user %s: %v", user.Username, err)
        }
        fmt.Println("Database seeded with initial User data: admin/password123.")
    } else {
        fmt.Println("User table already has data, skipping user seeding.")
    }
}