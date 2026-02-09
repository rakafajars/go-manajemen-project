package repositories

import (
	"github.com/rakafajars/go-manajemen-project/config"
	"github.com/rakafajars/go-manajemen-project/models"
)

type BoardMemberRepository interface {
	GetMembers(boardPublicID string) ([]models.User, error)
}

type boardMemberRepository struct{}

func NewBoardMemberRepository() BoardMemberRepository {
	return &boardMemberRepository{}
}

func (r *boardMemberRepository) GetMembers(boardPublicID string) ([]models.User, error) {
	var users []models.User

	// Query ini bertujuan untuk mendapatkan daftar user yang menjadi anggota dari sebuah board tertentu.
	// Kita menggunakan JOIN antar tabel karena relasinya many-to-many (User <-> BoardMember <-> Board).
	err := config.DB.
		// 1. JOIN ke tabel intermediate 'board_members'.
		//    Menghubungkan 'users.internal_id' dengan 'board_members.user_internal_id'.
		Joins("JOIN board_members ON board_members.user_internal_id = users.internal_id").

		// 2. JOIN ke tabel 'boards'.
		//    Menghubungkan 'board_members.board_internal_id' dengan 'boards.internal_id'.
		//    Ini diperlukan agar kita bisa memfilter berdasarkan 'public_id' milik board.
		Joins("JOIN boards ON boards.internal_id = board_members.board_internal_id").

		// 3. Filter berdasarkan 'public_id' board yang dicari.
		Where("boards.public_id = ?", boardPublicID).

		// 4. Eksekusi query dan simpan hasilnya ke variabel 'users'.
		Find(&users).Error

	return users, err
}
