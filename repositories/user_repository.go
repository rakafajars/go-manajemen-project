// Package repositories berisi semua repository untuk akses database
// Repository Pattern memisahkan logic bisnis dari operasi database
// Keuntungan: mudah testing, mudah ganti database, kode lebih terstruktur
package repositories

import (
	"github.com/rakafajars/go-manajemen-project/config"
	"github.com/rakafajars/go-manajemen-project/models"
)

// =============================================================================
// INTERFACE DEFINITION
// =============================================================================

// UserRepository adalah interface (kontrak) yang mendefinisikan
// semua operasi database yang berhubungan dengan User
//
// Kenapa pakai interface?
// 1. Abstraksi - handler/service tidak perlu tahu implementasi detail
// 2. Testing - mudah di-mock untuk unit testing
// 3. Flexibility - mudah ganti implementasi tanpa ubah kode yang menggunakan
//
// Contoh: Jika nanti ingin ganti dari GORM ke raw SQL,
// cukup buat struct baru yang implement interface ini
type UserRepository interface {
	// Create menyimpan user baru ke database
	// Parameter: pointer ke models.User yang akan disimpan
	// Return: error jika gagal, nil jika sukses
	Create(user *models.User) error

	// FindByEmail mencari user berdasarkan email
	// Parameter: email string yang akan dicari
	// Return: pointer ke models.User dan error
	// Jika user tidak ditemukan, akan return error "record not found"
	FindByEmail(email string) (*models.User, error)
}

// =============================================================================
// STRUCT IMPLEMENTATION
// =============================================================================

// userRepository adalah struct yang mengimplementasikan UserRepository interface
// Menggunakan huruf kecil (private) agar hanya bisa diakses melalui constructor
//
// Kenapa private?
// - Memaksa pengguna menggunakan constructor (NewUserRepository)
// - Menjaga enkapsulasi dan kontrol pembuatan instance
type userRepository struct {
	// Kosong karena menggunakan config.DB langsung
	// Alternatif: bisa tambahkan field db *gorm.DB untuk dependency injection
}

// =============================================================================
// CONSTRUCTOR
// =============================================================================

// NewUserRepository adalah constructor untuk membuat instance UserRepository
// Mengembalikan tipe interface, bukan struct (best practice di Go)
//
// Contoh penggunaan:
//
//	userRepo := repositories.NewUserRepository()
//	err := userRepo.Create(&user)
//
// Kenapa return interface?
// - Memungkinkan dependency injection
// - Memudahkan mocking saat testing
// - Menyembunyikan detail implementasi
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// =============================================================================
// METHOD IMPLEMENTATIONS
// =============================================================================

// Create menyimpan user baru ke database menggunakan GORM
//
// Parameter:
//   - user: pointer ke models.User yang akan disimpan
//
// Return:
//   - error: nil jika sukses, error message jika gagal
//
// GORM akan otomatis:
//   - Generate ID (jika pakai auto increment)
//   - Set CreatedAt dan UpdatedAt (jika ada field tersebut)
//   - Validasi constraint database (unique, not null, dll)
//
// Contoh penggunaan:
//
//	user := &models.User{
//	    Name:     "John Doe",
//	    Email:    "john@example.com",
//	    Password: hashedPassword,
//	}
//	err := userRepo.Create(user)
//	if err != nil {
//	    // handle error (email duplicate, validation failed, dll)
//	}
func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// FindByEmail mencari satu user berdasarkan email di database
//
// Parameter:
//   - email: alamat email yang akan dicari
//
// Return:
//   - *models.User: pointer ke user yang ditemukan
//   - error: nil jika ditemukan, error jika tidak ditemukan atau ada masalah
//
// Cara kerja GORM:
//   - Where("email = ?", email) → membuat query: SELECT * FROM users WHERE email = ?
//   - First(&user) → ambil 1 record pertama dan simpan ke variabel user
//   - .Error → mengambil error jika ada (contoh: record not found)
//
// Penggunaan "?" (placeholder):
//   - Mencegah SQL Injection attack
//   - GORM akan otomatis escape karakter berbahaya
//   - JANGAN gunakan: Where("email = '" + email + "'") ← BAHAYA!
//
// Contoh penggunaan:
//
//	user, err := userRepo.FindByEmail("john@example.com")
//	if err != nil {
//	    if errors.Is(err, gorm.ErrRecordNotFound) {
//	        // email tidak terdaftar
//	    }
//	    // handle error lain
//	}
//	// user ditemukan, lanjut proses
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
