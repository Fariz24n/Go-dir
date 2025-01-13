package handlers

import (
	"encoding/json"
	"net/http"

	"Apalah/models"

	"gorm.io/gorm"
)

func GudangHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var gudangs []models.Gudang
		if err := db.Find(&gudangs).Error; err != nil {
			http.Error(w, "Failed to fetch gudang", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gudangs)
	}
}
func CustomerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customers []models.Customer
		if err := db.Find(&customers).Error; err != nil {
			http.Error(w, "Failed to fetch customer", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}
}

func AisleHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var aisles []models.Aisle
		if err := db.Find(&aisles).Error; err != nil {
			http.Error(w, "Failed to fetch aisle", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(aisles)
	}
}

func ExecuteMoveProcedure(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Contoh pemanggilan prosedur
		err := db.Exec("CALL move_stock(?, ?)", "BOX123", 10).Error
		if err != nil {
			http.Error(w, "Failed to execute procedure: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Procedure executed successfully"))
	}
}
