-- Membuat tabel board_members untuk menyimpan data anggota board
-- Hubungan ini bersifat banyak-ke-banyak (many-to-many) antara board dan user
CREATE TABLE board_members (
    -- board_internal_id: ID board tempat user menjadi anggota
    -- Merujuk ke tabel boards kolom internal_id (Foreign Key)
    -- ON DELETE CASCADE: Jika board dihapus, data keanggotaan ini ikut terhapus
    board_internal_id BIGINT NOT NULL REFERENCES boards(internal_id) ON DELETE CASCADE,

-- user_internal_id: ID user yang menjadi anggota board
-- Merujuk ke tabel users kolom internal_id (Foreign Key)
-- NOTE: Sebelumnya merujuk ke boards(internal_id), tetapi seharusnya merujuk ke users(internal_id) agar logis
-- ON DELETE CASCADE: Jika user dihapus, data keanggotaan ini ikut terhapus
user_internal_id BIGINT NOT NULL REFERENCES users (internal_id) ON DELETE CASCADE,

-- joined_at: Waktu kapan user bergabung
-- DEFAULT NOW(): Jika tidak diisi, otomatis menggunakan waktu saat ini
joined_at TIMESTAMP NOT NULL DEFAULT NOW(),

-- PRIMARY KEY (Composite Key): Kombinasi board_id dan user_id harus unik
-- Mencegah user yang sama ditambahkan dua kali ke board yang sama
PRIMARY KEY (board_internal_id, user_internal_id) );