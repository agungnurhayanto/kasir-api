package services

import (
	"KASIR-API/models"
	"KASIR-API/repositories"
	"errors"

	"github.com/lib/pq"
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

func (s *ProductService) Create(product *models.Product) error {
	err := s.repo.Create(product)
	if err != nil {

		// 🔥 khusus foreign key violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				return errors.New("category_id tidak ditemukan")
			}
		}

		return err
	}

	return s.repo.FindByIDWithCategory(product.ID, product)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) error {
	err := s.repo.Update(product)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				return errors.New("category_id tidak ditemukan")
			}
		}

		return err
	}

	return nil
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
