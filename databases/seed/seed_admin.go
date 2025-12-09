// Package seed berisi fungsi-fungsi untuk mengisi data awal (seed) ke database.
// Seeding berguna untuk membuat data default seperti admin pertama atau data testing.
package seed

import (
	"log"

	"github.com/rakafajars/go-manajemen-project/config"
	"github.com/rakafajars/go-manajemen-project/models"
	"github.com/rakafajars/go-manajemen-project/utils"
)

// SeedAdmin membuat akun admin pertama di database.
// Fungsi ini biasanya dipanggil sekali saat aplikasi pertama kali dijalankan.
func SeedAdmin() {
	// 1. Hash password "admin" menggunakan bcrypt.
	// Kita TIDAK BOLEH menyimpan password plain text di database!
	// Underscore (_) mengabaikan error, tapi di production sebaiknya di-handle.
	password, _ := utils.HashPassword("admin")

	// 2. Buat objek User dengan data admin.
	admin := models.User{
		Name:     "Admin",
		Email:    "admin@admin.com",
		Password: password, // Password yang sudah di-hash
		Role:     "admin",
	}

	// 3. Simpan ke database menggunakan FirstOrCreate.
	//
	// FirstOrCreate adalah fitur GORM yang sangat berguna:
	// - Pertama, GORM akan MENCARI user dengan Email = "admin@admin.com".
	// - Jika DITEMUKAN: Tidak melakukan apa-apa (menghindari duplikasi).
	// - Jika TIDAK DITEMUKAN: Buat user baru dengan data dari variabel `admin`.
	//
	// Ini memastikan fungsi ini bisa dipanggil berkali-kali tanpa membuat duplikat.
	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil {
		log.Println("Failed to seed admin:", err)
	} else {
		log.Println("Admin seeded successfully")
	}
}
