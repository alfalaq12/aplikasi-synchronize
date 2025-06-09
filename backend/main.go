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
	"strings"

	"licensi-app/backend/server"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var licenseKey string
var licenseServer = "http://localhost:8080/validate" // bisa pointing ke server pusat

func main() {
	// Koneksi ke DB
	server.ConnectDB()

	// Jalankan server lokal lisensi (jika ada)
	go server.StartLicenseServer()

	app := fiber.New()
	app.Use(cors.New())

	// Serve file static
	app.Static("/", "../frontend/public")

	// Home route
	app.Get("/", func(c *fiber.Ctx) error {
		if c.Cookies("logged_in") != "true" {
			return c.Redirect("/login.html")
		}
		if checkLicense() {
			return c.Redirect("/dashboard.html")
		}
		return c.SendFile("../frontend/index.html")
	})

	// Route validasi lisensi
	app.Post("/validate", func(c *fiber.Ctx) error {
		type Request struct {
			Key string `json:"key"`
		}
		var req Request

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Baca dari file license.key
		data, err := os.ReadFile("license.key")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca file lisensi"})
		}

		fileLicense := string(data)
		if req.Key == fileLicense {
			return c.JSON(fiber.Map{"valid": true})
		} else {
			return c.JSON(fiber.Map{"valid": false})
		}
		// Validasi ke server pusat
		valid, err := validateLicenseWithServer(req.Key)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal validasi ke server pusat"})
		}
		if !valid {
			return c.Status(403).JSON(fiber.Map{"error": "Lisensi tidak valid"})
		}

		// Simpan lisensi
		licenseKey = req.Key
		err = ioutil.WriteFile("license.key", []byte(licenseKey), 0644)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan lisensi"})
		}

		go restartApp()

		return c.JSON(fiber.Map{"message": "Lisensi valid! Restart aplikasi..."})
	})

	app.Post("/register", server.Register)
	app.Post("/login", server.Login)

	fmt.Println("Server berjalan di port 441 🚀")
	log.Fatal(app.Listen(":441"))
}

// Validasi ke server pusat
func validateLicenseWithServer(key string) (bool, error) {
	data, _ := json.Marshal(map[string]string{"key": key})

	// Hindari request ke diri sendiri
	if strings.Contains(licenseServer, ":441") {
		fmt.Println("📛 Skip validasi ke diri sendiri")
		return key == "SYNXCHRO-1234-VALID", nil // dummy lisensi
	}

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

// Cek lisensi
func checkLicense() bool {
	data, err := os.ReadFile("license.key")
	return err == nil && len(data) > 0
}

// Restart
func restartApp() {
	fmt.Println("Restarting service...")
	cmd := exec.Command(os.Args[0])
	cmd.Start()
	os.Exit(0)
}
