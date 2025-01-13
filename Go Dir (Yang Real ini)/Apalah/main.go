package main

import (
	models "Apalah/models"
	"Apalah/routes"
	"log"
	"net/http"

	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "host=localhost user=postgres password=secret dbname=Apalah port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal membuka koneksi ke database:", err)
	}
	log.Println("Berhasil terhubung ke database")
}

// func corsValidationMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")
// 		allowedOrigins := []string{"http://localhost:8080"}
// 		allowed := false
// 		for _, o := range allowedOrigins {
// 			if origin == o {
// 				allowed = true
// 				break
// 			}
// 		}
// 		if !allowed {
// 			http.Error(w, "CORS validation failed: origin not allowed", http.StatusForbidden)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {
	InitDB()

	router := routes.SetupRoutes(DB)
	DB.AutoMigrate(&models.Gudang{})
	DB.AutoMigrate(&models.Aisle{})
	DB.AutoMigrate(&models.Cart{})
	DB.AutoMigrate(&models.Customer{})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	//routerWithCORSValidation := corsValidationMiddleware(router)

	log.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", corsHandler.Handler(router)); err != nil {
		log.Fatal(err)
	}
}
