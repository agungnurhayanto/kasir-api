package main

import (
	"KASIR-API/database"
	"KASIR-API/handlers"
	"KASIR-API/repositories"
	"KASIR-API/services"
	"KASIR-API/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string
	DBConn string
}

func main() {

	// ===== CONFIG =====
	utils.InitTimezone()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.DBConn == "" {
		log.Fatal("❌ DB_CONN belum diset")
	}

	// ===== ROUTER =====
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// ===== DB INIT =====
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("❌ DB gagal connect:", err)
	}

	repo := repositories.NewProductRepository(db)
	service := services.NewProductService(repo)
	handler := handlers.NewProductHandler(service)

	mux.HandleFunc("/api/produk", handler.HandleProducts)
	mux.HandleFunc("/api/produk/", handler.HandleProductByID)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// untuk bikin transaksi
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	mux.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	// untuk report transaksi
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	mux.HandleFunc("/api/report", reportHandler.GetReport)               // ini untuk end point option chalenge get all trx by star dan end date
	mux.HandleFunc("/api/report/hari-ini", reportHandler.GetTodayReport) // ini untuk end point hari-ini

	// ===== SERVER =====
	srv := &http.Server{
		Addr:         "0.0.0.0:" + config.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("🚀 Server running on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
