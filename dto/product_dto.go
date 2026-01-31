package dto

// Untuk POST /products
type ProductCreateRequest struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}

// Untuk PUT /products/{id}
type ProductUpdateRequest struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}

// ==========================
// RESPONSE DTO
// ==========================

type CategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductResponse struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Price    int              `json:"price"`
	Stock    int              `json:"stock"`
	Category CategoryResponse `json:"category"`
}
