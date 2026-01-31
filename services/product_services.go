package services

import (
	"KASIR-API/models"
	"KASIR-API/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

// func (s *ProductService) Create(data *models.Product) error {
// 	return s.repo.Create(data)
// }

func (s *ProductService) Create(product *models.Product) error {
	// simpan product
	if err := s.repo.Create(product); err != nil {
		return err
	}
	// 🔑 ambil ulang + preload category
	return s.repo.FindByIDWithCategory(product.ID, product)

}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
