package operasi

import (
	"encoding/json"
	"net/http"

	"Apalah/models" // Ganti dengan modul Anda

	"gorm.io/gorm"
)

func AddGudang(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var gudang models.Gudang
		err := json.NewDecoder(r.Body).Decode(&gudang)
		if err != nil {
			http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
			return
		}
		if err := db.Create(&gudang).Error; err != nil {
			http.Error(w, "Failed to insert gudang data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gudang data added successfully"})
	}
}
