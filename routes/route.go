package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rakafajars/go-manajemen-project/controllers"
)

// Setup (sebaiknya diganti jadi Capital agar bisa diakses dari main.go)
// Fungsi ini bertugas mendaftarkan semua endpoint URL ke aplikasi Fiber.
// Parameter:
// - app: instance dari Fiber App yang dibuat di main.go
// - uc: instance dari UserController yang sudah diinjeksi service
func Setup(app *fiber.App, uc *controllers.UserController) {
	// Memuat variabel environment dari file .env (misal: DB_PASSWORD, PORT, dll)
	// Catatan: Biasanya ini dipanggil sekali saja di main.go, tapi tidak apa-apa di sini untuk belajar.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	// Grouping Route (Optional): Bisa juga dibuat group misal api := app.Group("/api")

	// Definisi Route:
	// Method: POST
	// URL: /v1/auth/register
	// Handler: uc.Register (fungsi yang ada di user_controller.go)
	app.Post("/v1/auth/register", uc.Register)
	app.Post("/v1/auth/login", uc.Login)
}
