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
