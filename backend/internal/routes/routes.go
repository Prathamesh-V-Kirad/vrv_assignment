package routes

import (
	"backend/internal/controllers"
	"github.com/gofiber/fiber/v2"
)


func Stepup(app *fiber.App){
	app.Post("/api/register", controllers.Register)
	app.Post("api/login",controllers.Login)
	app.Get("api/user",controllers.User)
	app.Post("api/logout",controllers.Logout)

	app.Post("/api/tasks", controllers.CreateTask)    
    app.Get("/api/tasks", controllers.GetTasks)      
    app.Put("/api/tasks/:id", controllers.UpdateTask) 
    app.Delete("/api/tasks/:id", controllers.DeleteTask)
}