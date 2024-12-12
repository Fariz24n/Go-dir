package routes

import (
	"Apalah/handlers"
	"Apalah/operasi"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) http.Handler {
	// Buat router baru
	r := mux.NewRouter()
	r.HandleFunc("/aisle", handlers.AisleHandler(db)).Methods("GET")

	r.HandleFunc("/gudang", handlers.GudangHandler(db)).Methods("GET")
	r.HandleFunc("/gudang", operasi.AddGudang(db)).Methods("POST")

	r.HandleFunc("/cart", operasi.GetCart(db)).Methods("GET")
	r.HandleFunc("/cart", operasi.AddToCart(db)).Methods("POST")
	r.HandleFunc("/cart", operasi.DeleteFromCart(db)).Methods("DELETE")

	r.HandleFunc("/checkout", operasi.AddToCart(db)).Methods("POST")

	r.HandleFunc("/move-stock", handlers.MoveStock(db)).Methods("POST")
	r.HandleFunc("/move-carts", handlers.MoveCarts(db)).Methods("POST")

	r.HandleFunc("/purchase", operasi.PurchaseHandler(db)).Methods("POST")

	return r
}
