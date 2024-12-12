package operasi

import (
	"Apalah/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

func AddToCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cart models.Cart
		err := json.NewDecoder(r.Body).Decode(&cart)
		if err != nil {
			http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
			return
		}
		if err := db.Create(&cart).Error; err != nil {
			http.Error(w, "Failed to add product to cart", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product added to cart successfully"})
	}
}

func GetCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerID := r.URL.Query().Get("customer_id")
		var cartItems []models.Cart
		if err := db.Where("customer_id = ?", customerID).Find(&cartItems).Error; err != nil {
			http.Error(w, "Failed to retrieve cart items", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cartItems)
	}
}

func DeleteFromCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cart models.Cart
		err := json.NewDecoder(r.Body).Decode(&cart)
		if err != nil {
			http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
			return
		}
		if err := db.Delete(&cart).Error; err != nil {
			http.Error(w, "Failed to remove product from cart", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product removed from cart successfully"})
	}
}

func Checkout(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerID := r.URL.Query().Get("customer_id")
		var cartItems []models.Cart
		if err := db.Where("customer_id = ?", customerID).Find(&cartItems).Error; err != nil {
			http.Error(w, "Failed to retrieve cart items", http.StatusInternalServerError)
			return
		}

		if len(cartItems) == 0 {
			http.Error(w, "Cart is empty", http.StatusBadRequest)
			return
		}

		if err := db.Where("customer_id = ?", customerID).Delete(&models.Cart{}).Error; err != nil {
			http.Error(w, "Failed to complete checkout", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Checkout successful"})
	}
}
