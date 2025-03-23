package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Simulasi database lisensi yang valid
var validLicenses = map[string]bool{
	"ABC123": true,
	"XYZ789": true,
}

func StartLicenseServer() {
	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Key string `json:"key"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Cek apakah lisensi valid
		_, valid := validLicenses[req.Key]
		resp := map[string]bool{"valid": valid}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server lisensi berjalan di port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
