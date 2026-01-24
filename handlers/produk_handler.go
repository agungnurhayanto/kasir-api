package handlers

import (
	"KASIR-API/models"
	"KASIR-API/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		encoder.Encode(storage.Produk)

	case http.MethodPost:
		var newProduct models.Produk
		if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newProduct.ID = len(storage.Produk) + 1
		storage.Produk = append(storage.Produk, newProduct)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newProduct)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// GET localhost:8080/api/produk {id}
func GetProdukByID(w http.ResponseWriter, r *http.Request) {
	// mendapatkan id dari url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	// jika error
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// mencari produk berdasarkan id
	for _, p := range storage.Produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// PUT localhost:8080/api/produk {id}
func UpdateProdukByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	// ganti int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduk models.Produk
	err = json.NewDecoder(r.Body).Decode(&updatedProduk)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// loop produk cari id ganti sesuai dari put request
	for i := range storage.Produk {
		if storage.Produk[i].ID == id {
			updatedProduk.ID = id
			storage.Produk[i] = updatedProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduk)
			return
		}
	}

	http.Error(w, "Id produk belum ada", http.StatusNotFound)
}

// DELETE localhost:8080/api/produk {id}
func DeleteProdukByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	// ganti id int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	// loop produk cari ID, dapet dari indext yang mau di hapus
	for i, p := range storage.Produk {
		if p.ID == id {
			// bikin slice baru index data sebelum dan sesudah di hapus
			storage.Produk = append(storage.Produk[:i], storage.Produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Produk berhasil di hapus",
			})

			return
		}

	}
	http.Error(w, "Id produk belum ada", http.StatusNotFound)
}
