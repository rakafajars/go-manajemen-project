package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rakafajars/go-manajemen-project/config"
	"github.com/rakafajars/go-manajemen-project/controllers"
	"github.com/rakafajars/go-manajemen-project/databases/seed"
	"github.com/rakafajars/go-manajemen-project/repositories"
	"github.com/rakafajars/go-manajemen-project/routes"
	"github.com/rakafajars/go-manajemen-project/services"
)

// Main Function adalah entry point (titik masuk) dari seluruh aplikasi Go.
// Ketika aplikasi dijalankan, kode di dalam fungsi inilah yang dieksekusi pertama kali.
func main() {
	// 1. Load Konfigurasi & Database
	// Membaca file .env untuk mengambil password DB, Port, dll.
	config.LoadEnv()
	// Membuka koneksi ke Database (PostgreSQL) menggunakan GORM.
	config.ConnectDB()

	// 2. Database Seeding (Optional)
	// Membuat data awal (misal: Admin superuser) jika belum ada di database.
	seed.SeedAdmin()

	// 3. Inisialisasi Framework Fiber
	// Membuat instance aplikasi web server baru.
	app := fiber.New()

	// 4. Manual Dependency Injection
	// Kita merakit layer aplikasi dari paling dalam ke paling luar:
	// Repository (akses DB) -> Service (Business Logic) -> Controller (Handler HTTP)

	// a. Buat Repository (Layer akses data)
	userRepo := repositories.NewUserRepository()

	// b. Buat Service (Layer logika), butuh Repository
	userService := services.NewUserService(userRepo)

	// c. Buat Controller (Layer interface), butuh Service
	userController := controllers.NewUserController(userService)

	// Board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	// 5. Setup Routes
	// Mendaftarkan URL (endpoint) ke Controller yang sesuai.
	// Kita passing 'app' dan 'userController' agar route bisa menghubungkan URL ke fungsi di controller.
	routes.Setup(app, userController, boardController)

	// 6. Jalankan Server
	// Mengambil PORT dari konfigurasi (misal "8080")
	port := config.AppConfig.AppPort
	log.Println("Server running on port :", port)

	// app.Listen akan memblokir main thread dan standby menerima request.
	log.Fatal(app.Listen(":" + port))

}
