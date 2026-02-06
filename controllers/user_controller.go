package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/rakafajars/go-manajemen-project/models"
	"github.com/rakafajars/go-manajemen-project/services"
	"github.com/rakafajars/go-manajemen-project/utils"
)

// UserController struct bertanggung jawab untuk menangani request yang masuk (Handler).
// Struct ini memiliki field 'service' untuk mengakses logika bisnis (Business Logic).
// Konsep ini disebut Dependency Injection: Controller bergantung pada Service.
type UserController struct {
	service services.UserService
}

// NewUserController adalah constructor function untuk membuat instance baru dari UserController.
// Fungsi ini menerima dependency 's' (UserService) dan mengembalikan pointer ke UserController yang baru.
func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

// Register adalah method (handler) yang menangani endpoint registrasi user.
// (c *UserController) adalah Receiver, artinya fungsi ini milik struct UserController.
// ctx *fiber.Ctx berisi informasi request (body, header) dan fungsi untuk kirim response.
func (c *UserController) Register(ctx *fiber.Ctx) error {
	// 1. Siapkan variabel kosong untuk menampung data user dari request body
	user := new(models.User)

	// 2. BodyParser membaca JSON dari body request dan memasukkan datanya ke variabel 'user'.
	// Jika format JSON salah atau data tidak sesuai, akan mengembalikan error.
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	// 3. Panggil fungsi Register di layer Service untuk memproses logika bisnis (validasi, simpan DB, dll).
	// Controller mendelegasikan tugas "berat" ke Service.
	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	// --- REFACTOR: Menggunakan DTO (Data Transfer Object) untuk Response ---

	// Siapkan variabel struct khusus untuk response (yang tidak ada passwordnya)
	var userResponse models.UserResponse

	// Gunakan library 'copier' untuk menyalin data dari struct 'user' (Model DB) ke 'userResponse' (Response API).
	// Copier akan otomatis menyalin field yang NAMANYA SAMA (misal user.Name -> userResponse.Name).
	// Field 'Password' tidak akan tersalin, karena di struct UserResponse tidak ada field Password. Aman!
	_ = copier.Copy(&userResponse, &user)

	// 4. Jika sukses, kirim response berhasil ke client.
	// Yang dikirim sekarang adalah 'userResponse', bukan object 'user' mentah lagi.
	return utils.Success(ctx, "Registerasi Berhasil", userResponse)
}

// Login menangani proses autentikasi user
func (c *UserController) Login(ctx *fiber.Ctx) error {
	// 1. Definisikan struct sementara untuk menampung input JSON (hanya butuh email & password)
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 2. Parse body request ke struct 'body'. Jika format JSON salah, return error.
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}

	// 3. Panggil fungsi Login di layer Service. Logika pengecekan password dan user ada di sana.
	user, err := c.service.Login(body.Email, body.Password)

	// 4. Jika login gagal (user tidak ketemu atau password salah), kembalikan error Unauthorized (401)
	if err != nil {
		return utils.Unauthorized(ctx, "Login Gagal", err.Error())
	}

	// 5. Jika sukses, buat Token JWT (Access Token) untuk user tersebut
	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)

	// 6. Buat juga Refresh Token agar user tidak perlu login berulang kali saat token utama expired
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	// 7. Siapkan DTO response agar field sensitif (seperti password) tidak ikut terkirim
	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)

	// 8. Kirim response sukses:
	// - token: Access Token (dipakai untuk request API)
	// - refresh_token: Refresh Token (dipakai untuk perpanjang sesi)
	// - user: Data user tanpa password
	return utils.Success(ctx, "Login Succesfuly", fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          userResp,
	})
}
