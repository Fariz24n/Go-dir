package operasi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func MoveStockHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		itemID := req["item_id"].(string)
		quantity := int(req["quantity"].(float64))

		itemID = strings.TrimSpace(itemID)

		var item struct {
			No       int64   `json:"no"`
			Supplier string  `json:"supplier"`
			Item     string  `json:"item"`
			Quantity int     `json:"quantity"`
			Ket      string  `json:"ket"`
			Price    float64 `json:"price"`
		}

		if err := db.Table("gudangs").
			Where("item = ?", itemID).
			First(&item).Error; err != nil {
			http.Error(w, "Item not found in warehouse", http.StatusNotFound)
			return
		}

		if item.Quantity < quantity {
			http.Error(w, "Not enough stock in warehouse", http.StatusBadRequest)
			return
		}

		// Konversi
		var itemQuantity int
		if strings.Contains(strings.ToLower(item.Ket), "box") {
			itemQuantity = quantity * 100
		} else {
			itemQuantity = quantity
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(`
				UPDATE gudangs
				SET quantity = quantity - ?
				WHERE item = ?`, quantity, itemID).Error; err != nil {
				return fmt.Errorf("failed to update warehouse stock: %w", err)
			}

			if err := tx.Exec(`
				INSERT INTO aisles (no, item, quantity, supplier, ket, price)
				VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT (item)
				DO UPDATE SET quantity = aisles.quantity + EXCLUDED.quantity`,
				item.No, item.Item, item.Quantity, item.Supplier, item.Ket, item.Price).Error; err != nil {
				return fmt.Errorf("failed to insert into aisles: %w", err)
			}

			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Successfully moved %d of %s to aisles", itemQuantity, itemID)))
	}
}
