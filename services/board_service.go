package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rakafajars/go-manajemen-project/models"
	"github.com/rakafajars/go-manajemen-project/repositories"
)

// BoardService mendefinisikan "Business Logic" seputar Board.
// Service tidak berurusan langsung dengan database (itu tugas Repository).
// Service menggabungkan beberapa repository jika perlu (misal: cek user dulu, baru simpan board).
type BoardService interface {
	CreateBoard(board *models.Board) error
}

type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo  repositories.UserRepository
}

func NewBoardService(boardRepo repositories.BoardRepository, userRepo repositories.UserRepository) BoardService {
	return &boardService{boardRepo, userRepo}
}

// CreateBoard menangani logika pembuatan board baru.
// Menerima data board, memvalidasi owner, membuat ID baru, lalu menyimpannya.
func (s *boardService) CreateBoard(board *models.Board) error {
	// 1. Validasi Owner: Cek apakah user pemilik board (berdasarkan ID yang dikirim) ada di database.
	// Kita gunakan OwnerPublicID karena itulah yang diisi dari Token JWT.
	// board.PublicID saat ini masih kosong (0000...) karena belum di-generate.
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner not found") // Return error jika user tidak ketemu
	}

	// 2. Generate ID Baru: Kita buat UUID baru yang unik untuk board ini.
	// Ini akan menimpa ID lama (jika ada) karena ID board harus selalu baru saat create.
	board.PublicID = uuid.New()

	// 3. Set Owner Internal ID: Hubungkan board ini dengan user pemiliknya (Foreign Key).
	// Kita ambil internal_id dari user yang tadi kita temukan.
	board.OwnerID = user.InternalID

	// 4. Simpan ke Database: Panggil repository untuk query INSERT.
	return s.boardRepo.Create(board)
}
