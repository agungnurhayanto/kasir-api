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
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string
	DB_CONN string
}

func main() {

	// ===== CONFIG =====
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.ReadInConfig()
	}

	config := Config{
		Port:    viper.GetString("PORT"),
		DB_CONN: viper.GetString("DB_CONN"),
	}
	if config.Port == "" {
		config.Port = "8080"
	}

	// ===== ROUTES =====
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	srv := &http.Server{
		Addr:         "0.0.0.0:" + config.Port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("🚀 Server running on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	db, err := database.InitDB(config.DB_CONN)
	if err != nil {
		log.Println("⚠️ DB belum siap:", err)
	} else {
		repo := repositories.NewProductRepository(db)
		service := services.NewProductService(repo)
		handler := handlers.NewProductHandler(service)

		http.HandleFunc("/api/produk", handler.HandleProducts)
		http.HandleFunc("/api/produk/", handler.HandleProductByID)
	}

	select {}
}
