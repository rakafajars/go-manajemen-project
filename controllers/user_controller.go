package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// GetUser menangani request untuk melihat profil user lain (berdasarkan Public ID)
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	// 1. Ambil parameter 'id' dari URL (contoh: /api/users/123e4567-e89b...)
	// Ingat: ini adalah Public ID, bukan Internal ID database
	id := ctx.Params("id")

	// 2. Panggil service untuk cari user. Service akan query ke repo pakai Public ID.
	user, err := c.service.GetByPublicId(id)

	// 3. Jika user tidak ditemukan, return error 404 (Not Found)
	if err != nil {
		return utils.NotFound(ctx, "User tidak ditemukan", err.Error())
	}

	// 4. Mapping data user ke struct UserResponse (DTO)
	// Tujuannya agar field sensitif seperti password/ID internal tidak terekspos ke publik
	var userResponse models.UserResponse
	err = copier.Copy(&userResponse, &user)

	// 5. Cek error saat mapping (sangat jarang terjadi, tapi good practice untuk di-handle)
	if err != nil {
		return utils.BadRequest(ctx, "Internal Server Error", err.Error())
	}

	// 6. Kembalikan response JSON berisi data user
	return utils.Success(ctx, "User ditemukan", userResponse)
}

func (c *UserController) GetUserPagination(ctx *fiber.Ctx) error {
	// 1. Menangkap Query Parameters (Input dari URL)
	// Contoh: users/page?page=1&limit=10&sort=name&filter=search
	page, _ := strconv.Atoi(ctx.Query("page", "1"))    // Default page 1
	limit, _ := strconv.Atoi(ctx.Query("limit", "10")) // Default limit 10
	filter := ctx.Query("filter", "")                  // Keyword pencarian
	sort := ctx.Query("sort", "")                      // Format sorting

	// 2. Menghitung Offset
	// Menentukan data mulai diambil dari urutan ke berapa
	// Rumus: (page - 1) * limit
	offset := (page - 1) * limit

	// 3. Memanggil Service (Logika Bisnis)
	// Meminta data ke layer Service dengan parameter pagination
	users, total, err := c.service.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil data", err.Error())
	}

	// 4. Mapping Data (Konversi Struct)
	// Mengubah raw data database (models.User) ke format response (models.UserResponse)
	// agar aman dan formatnya sesuai kebutuhan client
	var userResp []models.UserResponse
	copier.Copy(&userResp, &users)

	// 5. Menyusun Metadata Pagination
	// Informasi pelengkap untuk frontend (total halaman, halaman saat ini, dll)
	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))), // Total data dibagi limit, dibulatkan ke atas
		Filter:    filter,
		Sort:      sort,
	}

	// 6. Mengirim Response
	if total == 0 {
		return utils.NotFoundPagination(ctx, "User tidak ditemukan", userResp, meta)
	}

	return utils.SuccessPagination(ctx, "User ditemukan", userResp, meta)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	// 1. Ambil ID dari URL (misal: /users/123-abc)
	id := ctx.Params("id")

	// 2. Parsing ID (Validasi format UUID)
	// Kita pastikan string ID yang dikirim user adalah format UUID yang valid
	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	// 3. Body Parser (Ambil data JSON dari request body)
	// Masukkan data JSON {"name": "Budi"} ke variabel user
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	// 4. Set Public ID ke object user
	// Agar repository tahu user mana yang mau diupdate
	user.PublicID = publicID

	// 5. Panggil Service Update
	// Lakukan proses update di database
	if err := c.service.Update(&user); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Data", err.Error())
	}

	// 6. Ambil Data Terbaru (Optional tapi Recommended)
	// Setelah update, kita ambil lagi data terbarunya dari DB untuk response
	// Supaya user langsung lihat hasil perubahannya
	userUpdated, err := c.service.GetByPublicId(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Ambil Data", err.Error())
	}

	// 7. Mapping ke Response
	// Pindahkan data ke format UserResponse agar bersih (tanpa password dll)
	var userResp models.UserResponse
	err = copier.Copy(&userResp, &userUpdated)

	if err != nil {
		return utils.InternalServerError(ctx, "Error parsing data", err.Error())
	}

	return utils.Success(ctx, "Berhasil Update Data", userResp)
}
