package operasi

import (
	"Apalah/handlers"
	"Apalah/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil || !claims["is_admin"].(bool) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func UserOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userInput models.Customer
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userInput); err != nil {
			response := map[string]string{"message": err.Error()}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer r.Body.Close()

		var user models.Customer
		if err := db.Where("name = ?", userInput.Name).First(&user).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				response := map[string]string{"message": "Username atau password salah"}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			default:
				response := map[string]string{"message": err.Error()}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
			response := map[string]string{"message": "Username atau password salah"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		expTime := time.Now().Add(time.Minute * 1)
		claims := &handlers.JWTClaim{
			Name: user.Name,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-jwt-mux",
				ExpiresAt: jwt.NewNumericDate(expTime),
			},
		}

		tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token, err := tokenAlgo.SignedString(handlers.JWT_KEY)
		if err != nil {
			response := map[string]string{"message": err.Error()}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Path:     "/",
			Value:    token,
			HttpOnly: true,
		})

		response := map[string]string{"message": "login berhasil"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userInput models.Customer
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userInput); err != nil {
			response := map[string]string{"message": err.Error()}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer r.Body.Close()

		// Validasi apakah address ada
		if userInput.Address == "" {
			response := map[string]string{"message": "Address is required"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Hash password
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		userInput.Password = string(hashPassword)

		// Simpan ke database
		if err := db.Create(&userInput).Error; err != nil {
			response := map[string]string{"message": err.Error()}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Respon sukses
		response := map[string]string{"message": "success"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func Logout(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Path:     "/",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})

		response := map[string]string{"message": "logout berhasil"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}
