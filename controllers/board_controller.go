package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rakafajars/go-manajemen-project/models"
	"github.com/rakafajars/go-manajemen-project/services"
	"github.com/rakafajars/go-manajemen-project/utils"
)

// BoardController bertugas menangani request HTTP/API yang masuk terkait Board.
// Controller ini menghubungkan input dari User (via HTTP) ke Business Logic di Service.
type BoardController struct {
	// service adalah dependency (ketergantungan) yang dibutuhkan controller ini.
	// Kita menggunakan interface 'BoardService', bukan struct konkretnya (boardService).
	service services.BoardService
}

// NewBoardController adalah fungsi konstruktor untuk membuat instance BoardController.
// Menerima 'BoardService' sebagai parameter (Dependency Injection).
func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

// CreateBoard menangani endpoint POST /boards untuk membuat board baru.
func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error
	board := new(models.Board)

	// 1. Ambil Data User dari JWT Token (Middleware)
	// Saat login, user dapat token. Middleware memvalidasi token itu & menyimpannya di ctx.Locals("user").
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// 2. Parsing Body Request: Mengubah JSON dari body request menjadi struct Board.
	// Jika format JSON salah atau tidak sesuai model, return error 400 (Bad Request).
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal membaca request", err.Error())
	}

	// 3. Ambil Public ID User dari Claims Token
	// Kita ambil "public_id" yang tersimpan di dalam token JWT.
	// Ini memastikan board yang dibuat otomatis menjadi milik user yang sedang login.
	userID, err = uuid.Parse(claims["public_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Gagal membaca ID user", err.Error())
	}

	// Set OwnerPublicID di struct board dengan ID user yang login.
	board.OwnerPublicID = userID

	// 4. Panggil Service: Menjalankan logic pembuatan board (validasi owner, generate ID, simpan ke DB).
	// Jika service mengembalikan error (misal owner tidak ketemu), return error 400.
	if err := c.service.CreateBoard(board); err != nil {
		return utils.BadRequest(ctx, "Gagal menyimpan data board", err.Error())
	}

	// 5. Response Sukses: Jika berhasil, kembalikan response 200/201 dengan data board yang baru dibuat.
	return utils.Success(ctx, "Board berhasil dibuat", board)
}

// UpdateBoard menangani request untuk mengubah data board.
func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	// 1. Ambil ID dari URL (misal: /boards/:id)
	publicID := ctx.Params("id")
	board := new(models.Board)

	// 2. Parse Body: Ubah JSON dari request body menjadi struct Board.
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	// 3. Validasi UUID: Pastikan ID yang dikirim formatnya benar.
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	// 4. Cek Keberadaan Data: Pastikan board yang mau diedit benar-benar ada.
	existingBoard, err := c.service.GetByPublicID(publicID) // Menggunakan method baru di service
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan", err.Error())
	}

	// 5. Set ID Internal: Kita harus memastikan board yang diupdate adalah board yang sama.
	// Kita ambil internal_id dan public_id dari data lama (existingBoard) dan pasang ke object baru (board).
	// Ini penting agar GORM tahu record mana yang harus diupdate (berdasarkan primary key / unique key).
	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.OwnerID = existingBoard.OwnerID
	board.CreatedAt = existingBoard.CreatedAt

	// 6. Panggil Service Update: Lakukan proses penyimpanan perubahan ke database.
	if err = c.service.UpdateBoard(board); err != nil {
		return utils.BadRequest(ctx, "Gagal update board", err.Error())
	}

	// 7. Berikan Response Sukses ke client.
	return utils.Success(ctx, "Board berhasil diupdate", board)
}

// AddBoardMembers menangani endpoint POST /boards/:id/members untuk menambahkan anggota.
func (c *BoardController) AddBoardMembers(ctx *fiber.Ctx) error {
	// 1. Ambil ID Board dari URL Parameter (misal: /boards/:id/members)
	publicID := ctx.Params("id")

	// 2. Siapkan variabel untuk menampung list ID User yang akan ditambahkan.
	// Kita mengharapkan format JSON body berupa array of strings: ["user_id_1", "user_id_2"]
	var userIDs []string

	// 3. Parsing Body Request: Mengubah JSON body menjadi slice string.
	// Jika format JSON tidak sesuai (bukan array string), return error 400.
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	// 4. Panggil Service: Jalankan logic penambahan member di layer service.
	// Service akan memvalidasi apakah board dan user valid, serta memastikan tidak ada duplikasi.
	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Member", err.Error())
	}

	// 5. Response Sukses: Kembalikan status 200 OK jika berhasil.
	return utils.Success(ctx, "Berhasil Menambahkan Member ke Board", nil)
}

// RemoveBoardMembers menangani endpoint DELETE /boards/:id/members untuk menghapus anggota.
func (c *BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	// 1. Ambil ID Board dari URL Parameter (misal: /boards/:id/members)
	publicID := ctx.Params("id")

	// 2. Siapkan variabel untuk menampung list ID User yang akan DIHAPUS.
	// Kita mengharapkan format JSON body berupa array of strings: ["user_id_1", "user_id_2"]
	var userIDs []string

	// 3. Parsing Body Request: Mengubah JSON body menjadi slice string.
	// Jika format JSON tidak sesuai (bukan array string), return error 400.
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	// 4. Panggil Service: Jalankan logic PENGHAPUSAN member di layer service.
	// Service akan memvalidasi data dan menghapus member yang sesuai.
	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Member", err.Error())
	}

	// 5. Response Sukses: Kembalikan status 200 OK jika berhasil.
	return utils.Success(ctx, "Berhasil Menghapus Member dari Board", nil)
}
