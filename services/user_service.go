// Package services berisi business logic aplikasi
// Service layer berada di antara Handler (controller) dan Repository (database)
// Service bertugas mengolah data, validasi bisnis, dan koordinasi antar repository
package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rakafajars/go-manajemen-project/models"
	"github.com/rakafajars/go-manajemen-project/repositories"
	"github.com/rakafajars/go-manajemen-project/utils"
)

// =============================================================================
// INTERFACE DEFINITION
// =============================================================================

// UserService adalah interface yang mendefinisikan semua operasi bisnis
// yang berhubungan dengan User
//
// Kenapa pakai interface?
// 1. Memudahkan testing (bisa di-mock)
// 2. Loose coupling - handler tidak tergantung langsung ke implementasi
// 3. Bisa punya multiple implementasi (contoh: UserServiceV1, UserServiceV2)
type UserService interface {
	// Register mendaftarkan user baru ke sistem
	// Melakukan: validasi email, hash password, set role, simpan ke database
	Register(user *models.User) error
}

// =============================================================================
// STRUCT IMPLEMENTATION
// =============================================================================

// userService adalah implementasi dari UserService interface
// Menggunakan huruf kecil (private) agar hanya bisa diakses via constructor
//
// Field repo adalah dependency yang di-inject melalui constructor
// Ini disebut "Dependency Injection" - service tidak membuat repository sendiri,
// tapi menerima dari luar (biasanya dari main.go atau wire.go)
type userService struct {
	repo repositories.UserRepository // dependency ke repository layer
}

// =============================================================================
// CONSTRUCTOR
// =============================================================================

// NewUserService adalah constructor untuk membuat instance UserService
//
// Parameter:
//   - repo: instance dari UserRepository yang sudah dibuat
//
// Return:
//   - UserService interface
//
// Ini adalah contoh Dependency Injection:
// - Service tidak peduli bagaimana repository dibuat
// - Service hanya butuh "sesuatu" yang memenuhi kontrak UserRepository
// - Memudahkan testing: bisa inject mock repository
//
// Contoh penggunaan di main.go:
//
//	userRepo := repositories.NewUserRepository()
//	userService := services.NewUserService(userRepo)
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// =============================================================================
// METHOD IMPLEMENTATIONS
// =============================================================================

// Register mendaftarkan user baru ke sistem
//
// Langkah-langkah yang dilakukan:
// 1. Cek apakah email sudah terdaftar
// 2. Hash password untuk keamanan
// 3. Set default role sebagai "user"
// 4. Generate PublicID (UUID) untuk identifikasi publik
// 5. Simpan user ke database
//
// Parameter:
//   - user: pointer ke models.User yang berisi data registrasi
//
// Return:
//   - error: nil jika sukses, error message jika gagal
//
// Security notes:
//   - Password WAJIB di-hash sebelum disimpan ke database
//   - Jangan pernah simpan plain text password!
//   - PublicID digunakan untuk expose ke API (bukan InternalID)
//
// Contoh penggunaan di handler:
//
//	func RegisterHandler(c *fiber.Ctx) error {
//	    var user models.User
//	    c.BodyParser(&user)
//	    err := userService.Register(&user)
//	    if err != nil {
//	        return utils.BadRequest(c, "Registration failed", err.Error())
//	    }
//	    return utils.Created(c, "User registered", user)
//	}
func (s *userService) Register(user *models.User) error {
	// =========================================================================
	// STEP 1: Cek email sudah terdaftar atau belum
	// =========================================================================
	// Kita cari user dengan email yang sama di database
	// Jika ditemukan (InternalID != 0), berarti email sudah dipakai
	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email sudah terdaftar")
	}

	// =========================================================================
	// STEP 2: Hash password untuk keamanan
	// =========================================================================
	// Password harus di-hash menggunakan bcrypt atau algoritma aman lainnya
	// Ini memastikan password tidak bisa dibaca meski database bocor
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err // jika hashing gagal, return error
	}

	// =========================================================================
	// STEP 3: Set data user sebelum disimpan
	// =========================================================================
	// Ganti plain password dengan hashed password
	user.Password = hashedPassword

	// Set default role sebagai "user" (bukan "admin")
	// Role bisa diubah manual di database atau oleh admin
	user.Role = "user"

	// Generate UUID untuk PublicID
	// PublicID digunakan untuk identifikasi user di API responses
	// Ini lebih aman daripada expose InternalID (auto increment)
	// Karena InternalID bisa ditebak (1, 2, 3, ...) sedangkan UUID random
	user.PublicID = uuid.New()

	// =========================================================================
	// STEP 4: Simpan user ke database
	// =========================================================================
	// Memanggil repository untuk menyimpan data
	// Jika berhasil return nil, jika gagal return error
	return s.repo.Create(user)
}
