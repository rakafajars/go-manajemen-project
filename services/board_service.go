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
	// CreateBoard menangani pembuatan board baru.
	CreateBoard(board *models.Board) error
	// UpdateBoard memperbarui data board yang ada.
	UpdateBoard(board *models.Board) error
	// GetByPublicID mengambil data board berdasarkan ID publiknya.
	GetByPublicID(publicID string) (*models.Board, error)
	AddMembers(boardPublicID string, userPublicIDS []string) error
}

type boardService struct {
	boardRepo       repositories.BoardRepository
	userRepo        repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

func NewBoardService(boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
) BoardService {
	return &boardService{boardRepo, userRepo, boardMemberRepo}
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

// UpdateBoard menangani logika bisnis untuk update board.
// Saat ini hanya memanggil repository, tapi nanti bisa ditambah validasi lain (misal: cek kepemilikan).
func (s *boardService) UpdateBoard(board *models.Board) error {
	return s.boardRepo.Update(board)
}

// GetByPublicID mengambil data detail board.
// Berguna untuk ditampilkan di halaman detail atau form edit.
func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}

// AddMembers menambahkan user ke dalam board.
// Menerima ID Board (public) dan list ID User (public).
// Fungsi ini memastikan user yang ditambahkan valid dan belum menjadi anggota.
func (s *boardService) AddMembers(boardPublicID string, userPublicIDS []string) error {
	// 1. Validasi Board: Pastikan board dengan ID tersebut ada.
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	// 2. Validasi User: Loop semua user ID yang dikirim.
	// Kita perlu mendapatkan Internal ID mereka untuk disimpan di database (karena tabel relasi pakai internal ID).
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDS {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user board not found: " + userPublicID)
		}

		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	// 3. Cek Anggota Lama: Ambil daftar anggota yang SUDAH ada di board ini.
	// Tujuannya agar tidak menambahkan user yang sama dua kali (duplikat).
	existingsMembers, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return err
	}

	// 4. Buat Map Anggota Lama: Ubah list anggota lama menjadi MAP agar pencarian lebih cepat (O(1)).
	// Key: User Internal ID, Value: true
	memberMap := make(map[uint]bool)
	for _, member := range existingsMembers {
		memberMap[uint(member.InternalID)] = true
	}

	// 5. Filter Anggota Baru: Loop user yang mau ditambahkan.
	// Jika user ID TIDAK ada di memberMap, berarti dia anggota baru -> tambahkan ke list 'newMembersIDs'.
	var newMembersIDs []uint
	for _, userID := range userInternalIDs {
		if !memberMap[userID] {
			newMembersIDs = append(newMembersIDs, userID)
		}
	}

	// 6. Cek Apakah Ada Anggota Baru: Jika list kosong (semua user sudah jadi anggota), langsung return sukses.
	// Tidak perlu panggil repository.
	if len(newMembersIDs) == 0 {
		return nil
	}

	// 7. Simpan ke Database: Panggil repository untuk INSERT anggota-anggota baru tersebut.
	return s.boardRepo.AddMember(uint(board.InternalID), newMembersIDs)

}
