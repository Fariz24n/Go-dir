package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func MoveCarts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Item     string `json:"item"`
			Quantity int    `json:"quantity"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.Item) == "" || req.Quantity <= 0 {
			http.Error(w, "Item and quantity must be valid", http.StatusBadRequest)
			return
		}

		procedureCall := `CALL move_carts(?, ?, ?);`
		if err := db.Exec(procedureCall, req.Item, req.Quantity).Error; err != nil {
			log.Printf("Error moving stock: %v", err)
			http.Error(w, "Failed to move stock", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Successfully moved %d of item '%s' to carts.\n", req.Quantity, req.Item)))
	}
}
