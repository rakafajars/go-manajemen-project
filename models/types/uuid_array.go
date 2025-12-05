package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UUIDArray adalah tipe data kustom (slice of uuid.UUID).
// Kita butuh ini karena GORM/Go secara default kadang bingung cara baca array dari database (misal PostgreSQL array).
type UUIDArray []uuid.UUID

// Scan adalah method khusus dari interface `sql.Scanner`.
//
// KENAPA PAKAI POINTER (*UUIDArray)?
// Karena method ini bertujuan untuk MENGUBAH (mengisi) nilai variable `a` dengan data dari database.
// Jika tidak pakai pointer, yang terisi hanya "copy"-an variable-nya, sedangkan variable aslinya tetap kosong.
func (a *UUIDArray) Scan(value interface{}) error {
	var str string

	// 1. Cek tipe data yang dikirim database.
	switch v := value.(type) {
	case []byte:
		str = string(v) // Ubah byte jadi string
	case string:
		str = v
	default:
		return errors.New("failed to parse UUID array: unsupported data type")
	}

	// 2. Bersihkan format array database.
	// Dari: "{uuid1,uuid2,uuid3}"
	// Jadi: "uuid1,uuid2,uuid3"
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	// 3. Pecah string berdasarkan koma "," menjadi slice string.
	parts := strings.Split(str, ",")

	// 4. Siapkan tempat penampungan (slice) baru.
	*a = make(UUIDArray, 0, len(parts))

	// 5. Loop setiap potongan string UUID.
	for _, s := range parts {
		// Bersihkan spasi atau tanda kutip jika ada.
		s = strings.TrimSpace(strings.Trim(s, `"`))
		if s == "" {
			continue
		}

		// Parse string menjadi object uuid.UUID.
		u, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("failed to parse UUID array: %v", err)
		}

		// Masukkan ke slice.
		*a = append(*a, u)
	}
	return nil
}

// Value adalah method khusus dari interface `driver.Valuer`.
//
// KENAPA TIDAK PAKAI POINTER (UUIDArray)?
// Karena method ini hanya butuh MEMBACA data untuk dikirim ke database.
// Kita tidak mengubah isi array-nya, jadi cukup kirim "value"-nya saja (copy) sudah aman dan efisien.
func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	// Ubah setiap UUID jadi string
	var parts []string
	for _, u := range a {
		parts = append(parts, u.String())
	}

	// Gabungkan dengan koma dan bungkus kurung kurawal
	// strings.Join menggabungkan ["a", "b"] jadi "a,b"
	return fmt.Sprintf("{%s}", strings.Join(parts, ",")), nil
}

// GormDataType memberi tahu GORM tipe data apa yang harus dibuat di database saat AutoMigrate.
// Tanpa ini, GORM mungkin bingung dan menganggapnya sebagai "text" atau "bytes".
// Dengan me-return "uuid[]", kita memaksa PostgreSQL membuat kolom bertipe ARRAY of UUID.
func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
