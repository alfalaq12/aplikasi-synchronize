package server

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register user baru
func Register(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password"})
	}

	user := User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	// Simpan ke database
	if err := DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan user"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil didaftarkan!"})
}

// Login user
func Login(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user User
	if err := DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Password salah"})
	}

	return c.JSON(fiber.Map{"message": "Login berhasil!"})
}
