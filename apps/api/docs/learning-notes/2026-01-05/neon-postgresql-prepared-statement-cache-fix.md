# Catatan Pembelajaran - Blog Go API

## 1. Masalah Login Setelah Migrasi Database

### Konteks Masalah
Setelah melakukan migrasi tabel `users`, terdapat perubahan nama kolom dari `rest_password_token` (typo) menjadi `reset_password_token` (benar). Hal ini menyebabkan error pada saat login.

### Database yang Digunakan
**Neon PostgreSQL** (Serverless Database)

### Error yang Muncul
```
prepare failed after ParseComplete: ERROR: cached plan must not change result type (SQLSTATE 0A000)
```

### Penyebab Masalah
Error ini terjadi karena:
- PostgreSQL menyimpan **prepared statement cache** untuk query yang sering digunakan
- Setelah migrasi, struktur tabel berubah (nama kolom berbeda)
- Cache lama masih menyimpan struktur tabel yang lama
- Ketika GORM mencoba menggunakan cache lama dengan struktur tabel baru, terjadi konflik

### Solusi yang Diterapkan

#### Konfigurasi Database Sebelumnya:
```go
db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
```

#### Konfigurasi Database yang Benar untuk Neon PostgreSQL:
```go
db, err := gorm.Open(postgres.New(postgres.Config{
    DSN:                  config.DatabaseURL,
    PreferSimpleProtocol: true, // WAJIB untuk Neon PostgreSQL - menonaktifkan prepared statements
}), &gorm.Config{
    PrepareStmt: false, // Juga menonaktifkan cache prepared statement dari GORM
})
```

### Penjelasan Solusi

#### 1. `PreferSimpleProtocol: true`
- **Fungsi**: Menggunakan protokol sederhana alih-alih protokol extended PostgreSQL
- **Kegunaan untuk Neon**:
  - Neon adalah serverless database yang bisa scale up/down secara dinamis
  - Koneksi bisa terputus dan dibuat ulang kapan saja
  - Prepared statements tidak reliabel di lingkungan serverless
- **Hasil**: Mengirim query SQL lengkap setiap kali, bukan menggunakan prepared statement

#### 2. `PrepareStmt: false`
- **Fungsi**: Menonaktifkan cache prepared statement dari GORM
- **Kegunaan**: Mencegah GORM menyimpan dan menggunakan kembali prepared statements
- **Hasil**: Setiap query dibuat fresh tanpa menggunakan cache

### Kesimpulan
Untuk menggunakan **Neon PostgreSQL** dengan GORM, **WAJIB** menggunakan konfigurasi di atas untuk menghindari masalah dengan prepared statement cache. Ini adalah best practice khusus untuk serverless PostgreSQL database.

### Referensi
- File: `internal/db/db.go`
- Dokumentasi: Neon PostgreSQL Serverless Database
- Framework: GORM v2 dengan driver PostgreSQL

---

**Dibuat pada**: 5 Januari 2026
**Penulis**: Blog Go API Development
