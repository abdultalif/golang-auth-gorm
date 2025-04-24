# 🛡️ Golang Auth GORM

Proyek ini adalah RESTful API Authentication Service yang dibangun dengan bahasa Go (Golang), menggunakan GORM sebagai ORM, dan PostgreSQL sebagai database. API ini menyediakan fitur autentikasi lengkap seperti registrasi, login, verifikasi email, dan middleware untuk proteksi rute.

## 📌 Deskripsi

Sistem ini dirancang untuk menyediakan backend otentikasi yang aman dan terstruktur menggunakan JWT dan email verification. Cocok digunakan sebagai starter untuk aplikasi berbasis user-authentication.

---

## ⚙️ Teknologi yang Digunakan

- [Golang](https://golang.org/) - Bahasa pemrograman utama
- [GORM](https://gorm.io/) - ORM untuk Golang
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Docker](https://www.docker.com/) & Docker Compose - Untuk environment dan containerisasi
- [JWT](https://jwt.io/) - JSON Web Token untuk autentikasi
- [Net/HTTP](https://pkg.go.dev/net/http) - HTTP router dan server
- Custom Middleware dan Validations

---

## ✨ Fitur

- ✅ Register pengguna
- ✅ Verifikasi OTP dari email
- ✅ Kirim ulang kode OTP
- ✅ Login dan pembuatan token JWT
- ✅ Refresh token JWT
- ✅ Lupa password
- ✅ Cek validitas token reset password
- 🔐 Ambil profil pengguna yang sedang login
- ✅ Logging request dan error handling global

---

## 📚 Dokumentasi API

Dokumentasi lengkap API beserta contoh request/response tersedia di Postman:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://documenter.getpostman.com/view/27633194/2sB2cbZJCH)

---

## 🚀 Cara Menjalankan Program

### 1. Clone Repositori

```bash
git clone https://github.com/abdultalif/golang-auth-gorm.git
cd golang-auth-gorm
cp .env.example .env.local & .env.production
docker-compose up --build
```
