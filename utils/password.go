// Package utils berisi fungsi-fungsi pembantu (helper) yang bisa dipakai di mana saja.
package utils

// Import library bcrypt untuk hashing password.
// bcrypt adalah algoritma hashing yang SANGAT AMAN untuk password.
import "golang.org/x/crypto/bcrypt"

// HashPassword mengubah password plain text menjadi hash yang aman untuk disimpan di database.
//
// KENAPA HARUS DI-HASH?
// - Jika database diretas, hacker tidak bisa langsung tahu password aslinya.
// - Hash bersifat "satu arah", tidak bisa di-reverse menjadi password asli.
//
// KENAPA PAKAI BCRYPT?
// - Bcrypt otomatis menambahkan "salt" (karakter acak) ke setiap password.
// - Artinya, password sama bisa menghasilkan hash yang berbeda.
// - Bcrypt juga dirancang untuk LAMBAT, sehingga hacker sulit melakukan brute-force.
//
// Parameter:
//   - password: Password asli dari user (plain text).
//
// Return:
//   - string: Password yang sudah di-hash (aman untuk disimpan di database).
//   - error: Error jika proses hashing gagal.
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword() menerima:
	// 1. []byte(password): Password diubah jadi byte array.
	// 2. bcrypt.DefaultCost: Tingkat "kerumitan" hashing (default: 10).
	//    Semakin tinggi cost, semakin lama proses hashing, tapi semakin aman.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Ubah hasil hash (byte array) kembali ke string, lalu kembalikan.
	return string(bytes), err
}
