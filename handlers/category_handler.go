package handlers

import (
	"KASIR-API/models"
	"KASIR-API/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// GET localhost:8080/api/categories dan tambah data nya
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		encoder.Encode(storage.Categories)

	case http.MethodPost:
		var newCategory models.Category
		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newCategory.ID = len(storage.Categories) + 1
		storage.Categories = append(storage.Categories, newCategory)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCategory)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// GET localhost:8080/api/categories {id}
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	// mendapatkan id dari url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	// jika error
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// mencari kategori berdasarkan id
	for _, c := range storage.Categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// update categories
func UpdateCategoryByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	// ganti int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var updatedCategory models.Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// loop kategori cari id ganti sesuai dari put request
	for i := range storage.Categories {
		if storage.Categories[i].ID == id {
			updatedCategory.ID = id
			storage.Categories[i] = updatedCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}

	http.Error(w, "Id category belum ada", http.StatusNotFound)
}

// delete category
func DeleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	// ganti id int
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	// loop kategori cari ID, dapet dari indext yang mau di hapus
	for i, p := range storage.Categories {
		if p.ID == id {
			// bikin slice baru index data sebelum dan sesudah di hapus
			storage.Categories = append(storage.Categories[:i], storage.Categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category berhasil di hapus",
			})

			return
		}

	}
	http.Error(w, "Id category belum ada", http.StatusNotFound)
}
