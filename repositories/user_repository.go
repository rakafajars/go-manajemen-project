// Package repositories berisi semua repository untuk akses database
// Repository Pattern memisahkan logic bisnis dari operasi database
// Keuntungan: mudah testing, mudah ganti database, kode lebih terstruktur
package repositories

import (
	"strings"

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

	// FindByID mencari user berdasarkan ID (primary key).
	//
	// Parameter:
	//   - id: uint (unsigned integer) ID user yang dicari.
	//
	// Return:
	//   - *models.User: Pointer ke data user yang ditemukan.
	//   - error: Error object (nil jika sukses, berisi error jika gagal, misal record not found).
	//
	// Cara kerja:
	//   - config.DB.First(&user, id):
	//     Ini adalah shorthand GORM untuk mencari record berdasarkan Primary Key.
	//     Sama dengan query SQL: SELECT * FROM users WHERE id = [id] ORDER BY id LIMIT 1
	FindByID(id uint) (*models.User, error)

	// FindByPublicID mencari user berdasarkan Public ID (bukan primary key).
	//
	// Parameter:
	//   - publicID: string Public ID user yang dicari.
	//
	// Return:
	//   - *models.User: Pointer ke data user yang ditemukan.
	//   - error: Error object (nil jika sukses, berisi error jika gagal, misal record not found).
	//
	// Cara kerja:
	//   - config.DB.Where("public_id = ?", publicID):
	//     Membuat query SQL: SELECT * FROM users WHERE public_id = [publicID] ORDER BY id LIMIT 1
	//   - First(&user): Mengambil record pertama yang cocok dan menyimpannya ke variabel 'user'.
	FindByPublicID(publicID string) (*models.User, error)

	FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)

	Update(user *models.User) error
	Delete(user *models.User) error
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

// FindByID mencari user berdasarkan ID (primary key).
//
// Parameter:
//   - id: uint (unsigned integer) ID user yang dicari.
//
// Return:
//   - *models.User: Pointer ke data user yang ditemukan.
//   - error: Error object (nil jika sukses, berisi error jika gagal, misal record not found).
//
// Cara kerja:
//   - config.DB.First(&user, id):
//     Ini adalah shorthand GORM untuk mencari record berdasarkan Primary Key.
//     Sama dengan query SQL: SELECT * FROM users WHERE id = [id] ORDER BY id LIMIT 1
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	// First(&user, id) otomatis mencari berdasarkan primary key karena argument kedua adalah integer/string ID.
	// Hasil query akan dimasukkan ke variabel 'user' (pass by reference).
	err := config.DB.First(&user, id).Error

	return &user, err
}

// FindByPublicID mencari user berdasarkan Public ID (bukan primary key).
//
// Parameter:
//   - publicID: string Public ID user yang dicari.
//
// Return:
//   - *models.User: Pointer ke data user yang ditemukan.
//   - error: Error object (nil jika sukses, berisi error jika gagal, misal record not found).
//
// Cara kerja:
//   - config.DB.Where("public_id = ?", publicID):
//     Membuat query SQL: SELECT * FROM users WHERE public_id = [publicID] ORDER BY id LIMIT 1
//   - First(&user): Mengambil record pertama yang cocok dan menyimpannya ke variabel 'user'.
func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User

	err := config.DB.Where("public_id = ?", publicID).First(&user).Error

	return &user, err
}

// FindAllPagination mengambil daftar user dengan fitur Filtering, Sorting, dan Pagination.
//
// Parameter:
//   - filter: Keyword pencarian (nama atau email).
//   - sort: Format sorting (contoh: "name" untuk ASC, "-name" untuk DESC).
//   - limit: Jumlah data per halaman.
//   - offset: Jumlah data yang dilewati (untuk paging).
//
// Return:
//   - []models.User: List user yang ditemukan.
//   - int64: Total data (untuk kalkulasi total pages di frontend).
//   - error: Error jika ada.
func (r *userRepository) FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Inisialisasi query builder GORM
	db := config.DB.Model(&models.User{})

	// =========================================================================
	// 1. FILTERING (Pencarian)
	// =========================================================================
	if filter != "" {
		// Gunakan wildcard '%' untuk pencarian parsial
		// Gunakan ILIKE (Case Insensitive LIKE) agar huruf besar/kecil dianggap sama
		// Contoh: filter "eri" akan mencocokkan "Eri", "ERI", "eRi", dll
		// ILIKE adalah fitur spesifik PostgreSQL
		filterPattern := "%" + filter + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ?", filterPattern, filterPattern)
	}

	// =========================================================================
	// 2. COUNTING (Hitung Total Data)
	// =========================================================================
	// Hitung total data SEBELUM limit & offset diterapkan.
	// Penting untuk frontend tahu berapa total halaman yang tersedia.
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// =========================================================================
	// 3. SORTING (Pengurutan)
	// =========================================================================
	if sort != "" {
		// Mapping khusus: jika user kirim "id", kita map ke "internal_id" (nama kolom DB)
		switch sort {
		case "-id":
			sort = "-internal_id"
		case "id":
			sort = "internal_id"
		}

		// Logic detect Ascending/Descending
		// Jika diawali "-", berarti DESC (Descending/Menurun)
		// Contoh: "-name" -> "name DESC"
		if strings.HasPrefix(sort, "-") {
			sort = strings.TrimPrefix(sort, "-") + " DESC"
		} else {
			// Jika tidak ada "-", berarti ASC (Ascending/Naik)
			sort = sort + " ASC"
		}

		// Aplikasikan sorting ke query
		db = db.Order(sort)
	}

	// =========================================================================
	// 4. PAGINATION & EXECUTION
	// =========================================================================
	// Limit: Berapa banyak data yang diambil
	// Offset: Berapa banyak data yang dilewati
	// Find: Eksekusi query final
	err := db.Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err

}

// Update memperbarui data user yang sudah ada.
//
// Parameter:
//   - user: Pointer ke struct user yang berisi data baru.
//
// Cara kerja:
//   - db.Model(&models.User{}): Memberitahu GORM tabel mana yang akan diupdate.
//   - .Where("public_id = ?", user.PublicID): Filter user mana yang akan diupdate.
//   - .Updates(map[string]interface{...}): Melakukan update HANYA pada kolom yang ditentukan di dalam map.
//     Kenapa pakai map? Agar kita bisa memilih secara spesifik kolom apa saja yang mau diubah (SELECTIVE UPDATE).
//     Jika pakai struct, GORM akan mengabaikan field yang kosong, dan kita tidak punya kontrol penuh.
func (r *userRepository) Update(user *models.User) error {
	return config.DB.Model(&models.User{}).Where("public_id = ?", user.PublicID).Updates(map[string]any{
		"name": user.Name, // Hanya kolom 'name' yang diupdate, kolom lain (email, password) aman.
	}).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return config.DB.Where("public_id", user.PublicID).Delete(&models.User{}).Error
}
