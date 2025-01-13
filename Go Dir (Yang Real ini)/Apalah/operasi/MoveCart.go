package operasi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func MoveItemToCartHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		itemID := req["item_id"].(string)
		quantity := int(req["quantity"].(float64))
		customerID := uint(req["customer_id"].(float64))

		itemID = strings.TrimSpace(itemID)

		var item struct {
			No       int64   `json:"no"`
			Item     string  `json:"item"`
			Quantity int     `json:"quantity"`
			Price    float64 `json:"price"`
		}

		if err := db.Table("aisles").
			Where("item = ?", itemID).
			First(&item).Error; err != nil {
			http.Error(w, "Item not found in aisles", http.StatusNotFound)
			return
		}

		if item.Quantity < quantity {
			http.Error(w, "Not enough stock in aisles", http.StatusBadRequest)
			return
		}
		if strings.Contains(strings.ToLower(item.Item), "box") {
			item.Price = item.Price / 100
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(`
				UPDATE aisles
				SET quantity = quantity - ?
				WHERE item = ?`, quantity, itemID).Error; err != nil {
				return fmt.Errorf("failed to update aisle stock: %w", err)
			}

			if err := tx.Exec(`
				INSERT INTO carts (customer_id, item, quantity, price, total)
				VALUES (?, ?, ?, ?, ?)
				ON CONFLICT (item)
				DO UPDATE SET quantity = carts.quantity + EXCLUDED.quantity,
							 total = carts.quantity * EXCLUDED.price
			`, customerID, item.Item, quantity, item.Price, item.Price*float64(quantity)).Error; err != nil {
				return fmt.Errorf("failed to insert/update into cart: %w", err)
			}

			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Successfully moved %d of %s to cart", quantity, itemID)))
	}
}
