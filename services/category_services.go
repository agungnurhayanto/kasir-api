package services

import (
	"KASIR-API/models"
	"KASIR-API/repositories"
	"errors"

	"github.com/lib/pq"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

// func (s *CategoryService) Delete(id int) error {
// 	return s.repo.Delete(id)
// }

func (s *CategoryService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				return errors.New("category masih digunakan oleh product")
			}
		}

		return err
	}

	return nil
}
