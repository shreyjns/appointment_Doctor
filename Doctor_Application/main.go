package main

import (
	"Doctor_Application/Appointment"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", Appointment.Register)
	app.Post("/api/login", Appointment.Login)
	app.Post("/api/logout", Appointment.Logout)
	app.Get("/api/user", Appointment.UserL)
	app.Post("/api/createdoctor", Appointment.CreateDoctor)
	app.Get("/api/getdoctor", Appointment.GetAvailability)
	app.Post("/api/book", Appointment.BookAppointment)

}
func main() {
	Appointment.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	Setup(app)

	app.Listen(":8000")
}
