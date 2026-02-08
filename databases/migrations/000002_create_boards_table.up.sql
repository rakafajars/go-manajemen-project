-- Membuat tabel 'boards' untuk menyimpan data papan kerja (board)
CREATE TABLE boards (
    -- ID internal (integer) sebagai Primary Key, otomatis bertambah (BIGSERIAL).
    -- Digunakan untuk relasi database internal agar lebih cepat.
    internal_id BIGSERIAL PRIMARY KEY,

-- ID publik (UUID) yang aman untuk diekspos ke API/Frontend.
-- DEFAULT gen_random_uuid() akan otomatis membuat UUID acak saat data dibuat.
public_id UUID NOT NULL DEFAULT gen_random_uuid (),

-- Judul board, tipe string maksimal 255 karakter, wajib diisi (NOT NULL).
title VARCHAR(255) NOT NULL,

-- Deskripsi board, tipe TEXT untuk teks panjang, boleh kosong (NULL).
description TEXT,

-- ID internal pemilik board, merujuk ke tabel 'users' kolom 'internal_id'.
-- Tipe datanya BIGINT harus sama dengan users(internal_id).
owner_internal_id BIGINT NOT NULL REFERENCES users (internal_id),

-- ID publik pemilik board, disalin untuk kemudahan query API.
owner_public_id UUID NOT NULL,

-- Waktu pembuatan, otomatis diisi waktu sekarang (NOW()) jika tidak diisi.
created_at TIMESTAMP NOT NULL DEFAULT NOW(),

-- CONSTRAINT (Batasan) untuk memastikan public_id selalu unik.
CONSTRAINT boards_public_id_unique UNIQUE (public_id),

-- Menambahkan Foreign Key constraint secara eksplisit pada owner_public_id
-- yang merujuk ke users(public_id).
-- ON DELETE CASCADE: Jika user dihapus, maka semua board milik user tersebut juga otomatis dihapus.
CONSTRAINT fk_boards_owner_public_id FOREIGN KEY (owner_public_id) REFERENCES users (public_id) ON DELETE CASCADE
);