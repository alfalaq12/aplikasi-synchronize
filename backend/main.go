package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"licensi-app/backend/server"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Variabel global untuk menyimpan lisensi
var licenseKey string
var licenseServer = "http://localhost:8080/validate" // URL server pusat

func main() {
	// Koneksi ke database PostgreSQL
	server.ConnectDB()

	// Jalankan server lisensi di background
	go server.StartLicenseServer()

	// Inisialisasi Fiber
	app := fiber.New()

	// Middleware CORS untuk menghindari masalah request dari frontend
	app.Use(cors.New())

	// Sajikan file statis dari frontend
	app.Static("/", "../frontend")

	// Route utama (Home)
	app.Get("/", func(c *fiber.Ctx) error {
		if checkLicense() {
			return c.Redirect("/dashboard.html") // Redirect ke dashboard jika lisensi valid
		}
		return c.SendFile("../frontend/index.html") // Jika belum valid, tampilkan form lisensi
	})

	// API untuk validasi lisensi
	app.Post("/validate-license", func(c *fiber.Ctx) error {
		type Request struct {
			Key string `json:"key"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Kirim lisensi ke server pusat
		valid, err := validateLicenseWithServer(req.Key)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal validasi ke server pusat"})
		}

		if !valid {
			return c.Status(403).JSON(fiber.Map{"error": "Lisensi tidak valid"})
		}

		// Simpan lisensi lokal
		licenseKey = req.Key
		err = ioutil.WriteFile("license.key", []byte(licenseKey), 0644)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan lisensi"})
		}

		// Restart aplikasi agar memuat lisensi baru
		go restartApp()

		return c.JSON(fiber.Map{"message": "Lisensi valid! Restart aplikasi..."})
	})

	// Endpoint register
	app.Post("/register", server.Register)

	// Endpoint login
	app.Post("/login", server.Login)

	// Jalankan server di port 441
	fmt.Println("Server berjalan di port 441 🚀")
	log.Fatal(app.Listen(":441"))
}

// Validasi lisensi ke server pusat
func validateLicenseWithServer(key string) (bool, error) {
	data, _ := json.Marshal(map[string]string{"key": key})
	resp, err := http.Post(licenseServer, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result map[string]bool
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result["valid"], nil
}

// Cek lisensi dari file
func checkLicense() bool {
	data, err := os.ReadFile("license.key")
	if err != nil || len(data) == 0 {
		return false
	}
	return true
}

// Restart aplikasi
func restartApp() {
	fmt.Println("Restarting service...")
	cmd := exec.Command(os.Args[0])
	cmd.Start()
	os.Exit(0)
}
