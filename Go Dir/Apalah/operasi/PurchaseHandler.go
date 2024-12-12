package operasi

import (
	"Apalah/models"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func PurchaseHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CustomerID uint `json:"customer_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input format", http.StatusBadRequest)
			return
		}

		var cartItems []models.Cart
		if err := db.Where("customer_id = ?", req.CustomerID).Find(&cartItems).Error; err != nil {
			http.Error(w, "Failed to fetch cart items", http.StatusInternalServerError)
			return
		}

		if len(cartItems) == 0 {
			http.Error(w, "No items in cart to purchase", http.StatusBadRequest)
			return
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			for _, item := range cartItems {
				transaction := models.Transaction{
					CustomerID: item.CustomerID,
					Item:       item.Item,
					Quantity:   item.Quantity,
					Price:      item.Price,
					TotalPrice: item.Total,
				}
				if err := tx.Create(&transaction).Error; err != nil {
					return fmt.Errorf("failed to create transaction record: %w", err)
				}
			}

			if err := tx.Where("customer_id = ?", req.CustomerID).Delete(&models.Cart{}).Error; err != nil {
				return fmt.Errorf("failed to clear cart: %w", err)
			}

			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Purchase completed successfully for customer ID %d", req.CustomerID)))
	}
}
