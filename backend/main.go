package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"backend/internal/initialize"
	"backend/internal/database"
	"backend/internal/routes"
	// "github.com/gofiber/fiber/v2/middleware/csrf"
	// "github.com/gofiber/fiber/v2/middleware/helmet"
)

func main() {
	fmt.Println("Hello")
	err := godotenv.Load()
	if err != nil {
        log.Println("Error loading .env file")
    }
	database.Connect()
	initialize.InitializePermissionsAndRoles()
	app := fiber.New()
    app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173" ,  
		AllowMethods: "GET,POST,PUT,DELETE",  
		AllowHeaders: "Content-Type,Authorization",  
		AllowCredentials: true,  
	}))
	
	// app.Use(csrf.New())
	// app.Use(helmet.New())
	
	routes.Stepup(app)

	port := os.Getenv("PORT")

	fmt.Println(port)

	log.Fatal(app.Listen(":" + port))
}
