package repositories

import (
	"time"

	"github.com/rakafajars/go-manajemen-project/config"
	"github.com/rakafajars/go-manajemen-project/models"
)

// BoardRepository adalah antarmuka (interface) yang mendefinisikan kontrak
// atau daftar fungsi apa saja yang BISA dilakukan terkait data Board.
// Interface ini berguna untuk abstraksi, memudahkan testing (mocking), dan menjaga kode tetap rapi.
type BoardRepository interface {
	// Create menambahkan data board baru ke database.
	// Menerima pointer ke model Board, dan mengembalikan error jika gagal.
	Create(board *models.Board) error
	Update(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
	AddMember(boardID uint, userIDs []uint) error
	RemoveMembers(boardID uint, userIDs []uint) error
}

// boardRepository adalah struktur (struct) konkret yang mengimplementasikan interface BoardRepository.
// Struct ini biasanya akan menyimpan koneksi database (misalnya *gorm.DB).
// Karena sifatnya 'private' (huruf awal kecil), struct ini tidak bisa diakses langsung dari luar package.
type boardRepository struct {
	// db *gorm.DB // Nanti kita akan tambahkan field ini
}

// NewBoardRepository adalah fungsi konstruktor (factory) untuk membuat instance baru dari boardRepository.
// Fungsi ini mengembalikan tipe 'BoardRepository' (interface), bukan struct konkretnya.
// Ini adalah pola umum di Go agar code lain (seperti Service/Controller) hanya bergantung pada interface.
func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

// Create adalah implementasi dari method Create di interface BoardRepository.
// Fungsi ini menerima data board (sebagai pointer agar hemat memori & bisa diubah jika perlu)
// lalu menyimpannya ke database menggunakan GORM.
func (r *boardRepository) Create(board *models.Board) error {
	// config.DB adalah koneksi database global (bisa juga diinject ke struct boardRepository).
	// .Create(board) akan membuat query SQL INSERT INTO boards ...
	// .Error akan mengembalikan nil jika sukses, atau error object jika gagal.
	return config.DB.Create(board).Error
}

// Update mengubah data board yang sudah ada di database.
// Menggunakan map[string]any agar kita bisa memilih kolom spesifik yang ingin di-update.
func (r *boardRepository) Update(board *models.Board) error {
	// 1. config.DB.Model(&models.Board{}): Memberitahu GORM tabel mana yang akan di-update.
	// 2. .Where("public_id = ?", ...): Filter baris yang mau diubah berdasarkan PublicID.
	// 3. .Updates(map...): Melakukan update ke beberapa kolom sekaligus.
	return config.DB.Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.DueDate,
	}).Error
}

// FindByPublicID mencari satu data board berdasarkan Public ID (UUID).
// Mengembalikan pointer ke model Board dan error jika data tidak ditemukan.
func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board // Siapkan variabel kosong untuk menampung hasil query

	// 1. .Where("public_id = ?", ...): Cari data dengan syarat public_id tertentu.
	// 2. .First(&board): Ambil data PERTAMA yang ditemukan, lalu masukkan ke variabel board.
	err := config.DB.Where("public_id = ?", publicID).First(&board).Error

	return &board, err
}

// AddMember menambahkan banyak member sekaligus ke dalam board.
// Ini menggunakan konsep "Bulk Insert" agar efisien (cukup 1 query ke DB).
func (r *boardRepository) AddMember(boardID uint, userIDs []uint) error {
	// 1. Cek User Kosong: Jika tidak ada user yang mau ditambahkan, langsung return nil.
	// Ini mencegah error atau query kosong ke database.
	if len(userIDs) == 0 {
		return nil
	}

	// 2. Siapkan Data: Kita buat slice/array of struct 'BoardMember'.
	// BoardMember adalah tabel perantara (pivot) untuk relasi Many-to-Many antara Board dan User.
	now := time.Now()
	var members []models.BoardMember

	// 3. Looping: Kita ubah daftar UserID menjadi daftar object BoardMember.
	for _, userID := range userIDs {
		members = append(members, models.BoardMember{
			BoardID:  int64(boardID), // ID Board yang mau diisi member
			UserID:   int64(userID),  // ID User yang jadi member
			JoinedAt: now,            // Waktu join (semua disamakan sekarang)
		})
	}

	// 4. Bulk Insert: GORM pintar, jika kita kasih slice/array ke .Create(),
	// dia akan otomatis membuat query "INSERT INTO ... VALUES (...), (...), (...)"
	// Jadi hanya butuh 1 kali trip ke database untuk menyimpan banyak data sekaligus.
	return config.DB.Create(&members).Error
}

// RemoveMembers menghapus anggota dari board.
// Menerima boardID dan list userIDs yang akan dihapus.
func (r *boardRepository) RemoveMembers(boardID uint, userIDs []uint) error {
	// 1. Cek Kosong: Jika list userIDs kosong, tidak perlu melakukan apa-apa.
	if len(userIDs) == 0 {
		return nil
	}

	// 2. Hapus Data: Melakukan query DELETE untuk menghapus baris di tabel 'board_members'.
	//    Where: board_internal_id = ? AND user_internal_id IN (?)
	//    Ini akan menghapus semua record yang cocok dalam satu kali query.
	return config.DB.Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).Delete(&models.BoardMember{}).Error
}
