package main

import (
	"KASIR-API/handlers"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	//GET  localhost:8080/api/categories -> menampilkan semua kategori dan add data
	http.HandleFunc("/api/categories", handlers.CategoriesHandler)

	//GET  localhost:8080/api/produk -> menampilkan semua produk dan add data
	http.HandleFunc("/api/produk", handlers.ProductsHandler)
	// END POINT CRUD CATEGORY BY ID

	// END POINT CRUD CETEGORY
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handlers.GetCategoryByID(w, r)
		} else if r.Method == "PUT" {
			handlers.UpdateCategoryByID(w, r)
		} else if r.Method == "DELETE" {
			handlers.DeleteCategoryByID(w, r)
		}
	})

	// END POINT CRUD PRODUK
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handlers.GetProdukByID(w, r)
		} else if r.Method == "PUT" {
			handlers.UpdateProdukByID(w, r)
		} else if r.Method == "DELETE" {
			handlers.DeleteProdukByID(w, r)
		}
	})

	// cek aplikasi running / tidak / health check / status aplikasi
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Kasir Api Siap",
		})
	})

	fmt.Println("Server is running on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Kasir Api gagal berjalan")
	}

}
