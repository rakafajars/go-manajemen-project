// Package utils berisi fungsi-fungsi helper yang digunakan di seluruh aplikasi
// File ini khusus untuk menstandarisasi format response API
package utils

import "github.com/gofiber/fiber/v2"

// =============================================================================
// STRUCT DEFINITIONS (Definisi Struktur Data)
// =============================================================================

// Response adalah struktur standar untuk response API
// Semua response akan menggunakan format yang sama agar konsisten
//
// Contoh output JSON:
//
//	{
//	  "status": "Success",
//	  "response_code": 200,
//	  "message": "Login successful",
//	  "data": { ... }
//	}
//
// Penjelasan field:
//   - Status: status operasi (Success, Error, Created, dll)
//   - ResponseCode: HTTP status code (200, 400, 404, 500, dll)
//   - Message: pesan yang menjelaskan hasil operasi
//   - Data: data yang dikembalikan (optional, menggunakan interface{} agar bisa menerima tipe data apapun)
//   - Error: pesan error jika terjadi kesalahan (optional)
//
// Tag `json:"..."` digunakan untuk menentukan nama field saat di-serialize ke JSON
// Tag `omitempty` artinya field tidak akan ditampilkan jika nilainya kosong/nil
type Response struct {
	Status       string      `json:"status"`
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        string      `json:"error,omitempty"`
}

// ResponsePaginated adalah struktur response untuk data yang menggunakan pagination
// Mirip dengan Response, tapi ditambah field Meta untuk informasi pagination
//
// Contoh output JSON:
//
//	{
//	  "status": "Success",
//	  "response_code": 200,
//	  "message": "Data retrieved successfully",
//	  "data": [ ... ],
//	  "meta": {
//	    "page": 1,
//	    "limit": 10,
//	    "total": 100,
//	    "total_pages": 10
//	  }
//	}
type ResponsePaginated struct {
	Status       string         `json:"status"`
	ResponseCode int            `json:"response_code"`
	Message      string         `json:"message,omitempty"`
	Data         interface{}    `json:"data,omitempty"`
	Error        string         `json:"error,omitempty"`
	Meta         PaginationMeta `json:"meta"`
}

// PaginationMeta berisi informasi tentang pagination
// Digunakan untuk membantu client mengetahui halaman saat ini, total data, dll
//
// Penjelasan field:
//   - Page: halaman saat ini (contoh: 1, 2, 3)
//   - Limit: jumlah data per halaman (contoh: 10, 20, 50)
//   - Total: total keseluruhan data di database
//   - TotalPage: total halaman yang tersedia
//   - Filter: filter yang digunakan (contoh: "nama=triady")
//   - Sort: urutan data (contoh: "-id" artinya descending by id)
//
// Tag `example:"..."` digunakan untuk dokumentasi Swagger/OpenAPI
type PaginationMeta struct {
	Page      int    `json:"page" example:"1"`
	Limit     int    `json:"limit" example:"10"`
	Total     int    `json:"total" example:"100"`
	TotalPage int    `json:"total_pages" example:"10"`
	Filter    string `json:"filter" example:"nama=triady"`
	Sort      string `json:"sort" example:"-id"`
}

// =============================================================================
// SUCCESS RESPONSES (Response untuk operasi yang berhasil)
// =============================================================================

// Success mengirim response sukses dengan HTTP status 200 (OK)
// Digunakan untuk operasi yang berhasil seperti: GET data, Login, dll
//
// Parameter:
//   - c: Fiber context (untuk mengirim response ke client)
//   - message: pesan sukses yang akan ditampilkan
//   - data: data yang akan dikembalikan (bisa struct, slice, map, dll)
//
// Contoh penggunaan di handler:
//
//	func GetUser(c *fiber.Ctx) error {
//	    user := User{Name: "John", Email: "john@example.com"}
//	    return utils.Success(c, "User retrieved successfully", user)
//	}
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:       "Success",
		ResponseCode: fiber.StatusOK, // 200
		Message:      message,
		Data:         data,
	})
}

// SuccessPagination mengirim response sukses dengan pagination info
// Digunakan untuk operasi yang mengembalikan list data dengan pagination
//
// Contoh penggunaan:
//
//	func GetUsers(c *fiber.Ctx) error {
//	    users := []User{...}
//	    meta := utils.PaginationMeta{Page: 1, Limit: 10, Total: 100, TotalPage: 10}
//	    return utils.SuccessPagination(c, "Users retrieved", users, meta)
//	}
func SuccessPagination(c *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(ResponsePaginated{
		Status:       "Success",
		ResponseCode: fiber.StatusOK,
		Message:      message,
		Data:         data,
		Meta:         meta,
	})
}

// Created mengirim response sukses dengan HTTP status 201 (Created)
// Digunakan setelah berhasil membuat data baru (POST request)
//
// Contoh penggunaan:
//
//	func CreateUser(c *fiber.Ctx) error {
//	    newUser := User{...}
//	    // simpan ke database...
//	    return utils.Created(c, "User created successfully", newUser)
//	}
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:       "Created",
		ResponseCode: fiber.StatusCreated, // 201
		Message:      message,
		Data:         data,
	})
}

// =============================================================================
// ERROR RESPONSES (Response untuk operasi yang gagal)
// =============================================================================

// BadRequest mengirim response error dengan HTTP status 400 (Bad Request)
// Digunakan ketika request dari client tidak valid
// Contoh: validation error, format data salah, field required tidak diisi
//
// Contoh penggunaan:
//
//	func CreateUser(c *fiber.Ctx) error {
//	    if err := c.BodyParser(&user); err != nil {
//	        return utils.BadRequest(c, "Invalid request body", err.Error())
//	    }
//	}
func BadRequest(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Status:       "Error Bad Request",
		ResponseCode: fiber.StatusBadRequest, // 400
		Message:      message,
		Error:        err,
	})
}

// NotFound mengirim response error dengan HTTP status 404 (Not Found)
// Digunakan ketika data yang dicari tidak ditemukan di database
//
// Contoh penggunaan:
//
//	func GetUser(c *fiber.Ctx) error {
//	    user, err := findUserByID(id)
//	    if err != nil {
//	        return utils.NotFound(c, "User not found", err.Error())
//	    }
//	}
func NotFound(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Status:       "Error Not Found",
		ResponseCode: fiber.StatusNotFound, // 404
		Message:      message,
		Error:        err,
	})
}

// NotFoundPagination mengirim response not found dengan pagination info
// Digunakan untuk list data dengan pagination yang hasilnya kosong
func NotFoundPagination(c *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusNotFound).JSON(ResponsePaginated{
		Status:       "Not Found",
		ResponseCode: fiber.StatusNotFound,
		Message:      message,
		Data:         data,
		Meta:         meta,
	})
}

// Unauthorized mengirim response error dengan HTTP status 401 (Unauthorized)
// Digunakan ketika user tidak terautentikasi atau token tidak valid
// Contoh: token expired, token tidak ada, login gagal
//
// Contoh penggunaan:
//
//	func ProtectedRoute(c *fiber.Ctx) error {
//	    token := c.Get("Authorization")
//	    if token == "" {
//	        return utils.Unauthorized(c, "Token required", "No token provided")
//	    }
//	}
func Unauthorized(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Status:       "Error Unauthorized",
		ResponseCode: fiber.StatusUnauthorized, // 401
		Message:      message,
		Error:        err,
	})
}

// InternalServerError mengirim response error dengan HTTP status 500 (Internal Server Error)
// Digunakan ketika terjadi error di sisi server yang tidak terduga
// Contoh: database connection error, panic, file system error
//
// PENTING: Jangan expose detail error internal ke client di production!
// Gunakan pesan generic dan log error detail di server
//
// Contoh penggunaan:
//
//	func GetUsers(c *fiber.Ctx) error {
//	    users, err := db.Find(&users).Error
//	    if err != nil {
//	        return utils.InternalServerError(c, "Failed to fetch users", "Database error")
//	    }
//	}
func InternalServerError(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Status:       "Internal Server Error",
		ResponseCode: fiber.StatusInternalServerError, // 500
		Message:      message,
		Error:        err,
	})
}
