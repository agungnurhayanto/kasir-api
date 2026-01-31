package services

import (
	"KASIR-API/dto"
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

func (s *ProductService) GetAll() ([]dto.ProductResponse, error) {
	data, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.ProductResponse
	for _, item := range data {
		result = append(result, dto.ProductResponse{
			ID:    item.Product.ID,
			Name:  item.Product.Name,
			Price: item.Product.Price,
			Stock: item.Product.Stock,
			Category: dto.CategoryResponse{
				ID:          item.Category.ID,
				Name:        item.Category.Name,
				Description: item.Category.Description,
			},
		})
	}

	return result, nil
}

func (s *ProductService) Create(req dto.ProductCreateRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
	}

	err := s.repo.Create(&product)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				return nil, errors.New("category_id tidak ditemukan")
			}
		}
		return nil, err
	}

	// ambil data lengkap (JOIN)
	data, err := s.repo.GetByID(product.ID)
	if err != nil {
		return nil, err
	}

	res := dto.ProductResponse{
		ID:    data.Product.ID,
		Name:  data.Product.Name,
		Price: data.Product.Price,
		Stock: data.Product.Stock,
		Category: dto.CategoryResponse{
			ID:          data.Category.ID,
			Name:        data.Category.Name,
			Description: data.Category.Description,
		},
	}

	return &res, nil
}

func (s *ProductService) GetByID(id int) (*dto.ProductResponse, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	res := dto.ProductResponse{
		ID:    data.Product.ID,
		Name:  data.Product.Name,
		Price: data.Product.Price,
		Stock: data.Product.Stock,
		Category: dto.CategoryResponse{
			ID:          data.Category.ID,
			Name:        data.Category.Name,
			Description: data.Category.Description,
		},
	}

	return &res, nil
}

func (s *ProductService) Update(id int, req dto.ProductUpdateRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		ID:         id,
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
	}

	if err := s.repo.Update(&product); err != nil {
		return nil, err
	}

	data, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	res := dto.ProductResponse{
		ID:    data.Product.ID,
		Name:  data.Product.Name,
		Price: data.Product.Price,
		Stock: data.Product.Stock,
		Category: dto.CategoryResponse{
			ID:          data.Category.ID,
			Name:        data.Category.Name,
			Description: data.Category.Description,
		},
	}

	return &res, nil
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
