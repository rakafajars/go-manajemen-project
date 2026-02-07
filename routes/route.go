package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
	"github.com/rakafajars/go-manajemen-project/config"
	"github.com/rakafajars/go-manajemen-project/controllers"
	"github.com/rakafajars/go-manajemen-project/utils"
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

	// =========================================================================
	// JWT MIDDLEWARE CONFIGURATION
	// =========================================================================
	// Middleware ini bertugas sebagai "Satpam" yang memeriksa setiap request.
	// Jika request tidak membawa token valid, maka akan ditolak (401 Unauthorized).

	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		// 1. SigningKey: Kunci rahasia untuk memvalidasi tanda tangan token.
		// Harus SAMA PERSIS dengan secret key yang dipakai saat generate token.
		SigningKey: []byte(config.AppConfig.JWTSecret),

		// 2. ContextKey: Nama variabel untuk menyimpan claims token di context Fiber.
		// Jadi nanti di controller kita bisa akses data user lewat c.Locals("user").
		ContextKey: "user",

		// 3. ErrorHandler: Apa yang dilakukan jika token tidak valid atau tidak ada.
		// Di sini kita return JSON error standard 401 Unauthorized.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Unauthorized", err.Error())
		},
	}))

	// =========================================================================
	// PROTECTED ROUTES
	// =========================================================================
	// Grouping untuk endpoint user (prefix: /api/v1/users)
	userGroup := api.Group("/users")

	// GET /api/v1/users/:id
	// Endpoint ini otomatis dilindungi oleh JWT Middleware di atas.
	// Hanya user yang login (punya token valid) yang bisa akses.
	userGroup.Get("/page", uc.GetUserPagination)
	userGroup.Get("/:id", uc.GetUser)
	userGroup.Put("/:id", uc.UpdateUser)
	userGroup.Delete("/:id", uc.DeleteUser)
}
