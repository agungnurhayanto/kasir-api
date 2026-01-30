package main

import (
	"KASIR-API/database"
	"KASIR-API/handlers"
	"KASIR-API/repositories"
	"KASIR-API/services"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DB_CONN string `mapstructure:"DB_CONN"`
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Gagal membaca .env:", err)
		}
	}

	config := Config{
		Port:    viper.GetString("PORT"),
		DB_CONN: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.DB_CONN == "" {
		log.Fatal("DB_CONN kosong, pastikan .env terbaca")
	}

	db, err := database.InitDB(config.DB_CONN)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// cek aplikasi running / tidak / health check / status aplikasi
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Kasir Api Siap",
		})
	})

	log.Println("🚀 Server running on port", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))

}
