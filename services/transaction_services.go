package services

import (
	"KASIR-API/dto"
	"KASIR-API/models"
	"KASIR-API/repositories"
	"KASIR-API/utils"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

// func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
// 	return s.repo.CreateTransaction(items)
// }

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*dto.TransactionResponse, error) {
	trx, err := s.repo.CreateTransaction(items)
	if err != nil {
		return nil, err
	}

	// Convert waktu UTC → Asia/Jakarta
	created := trx.CreatedAt.In(utils.AppLocation)

	var detailsResponse []dto.TransactionDetailResponse
	for _, d := range trx.Details {
		detailsResponse = append(detailsResponse, dto.TransactionDetailResponse{
			ID:            d.ID,
			TransactionID: d.TransactionID,
			ProductID:     d.ProductID,
			ProductName:   d.ProductName,
			Quantity:      d.Quantity,
			Subtotal:      d.Subtotal,
		})
	}

	response := dto.TransactionResponse{
		ID:          trx.ID,
		TotalAmount: trx.TotalAmount,
		CreatedAt:   created,
		Details:     detailsResponse,
	}
	return &response, nil
}
