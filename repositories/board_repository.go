package repositories

import (
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
