package repositories

import (
	"KASIR-API/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	var createdAt time.Time

	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	if len(details) > 0 {

		valueStrings := []string{}
		valueArgs := []interface{}{}
		argID := 1

		for _, d := range details {
			valueStrings = append(valueStrings,
				fmt.Sprintf("($%d,$%d,$%d,$%d)", argID, argID+1, argID+2, argID+3))

			valueArgs = append(valueArgs,
				transactionID,
				d.ProductID,
				d.Quantity,
				d.Subtotal,
			)

			argID += 4
		}

		query := fmt.Sprintf(`
			INSERT INTO transaction_details
			(transaction_id, product_id, quantity, subtotal)
			VALUES %s
			RETURNING id
		`, strings.Join(valueStrings, ","))

		rows, err := tx.Query(query, valueArgs...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		i := 0
		for rows.Next() {
			if err := rows.Scan(&details[i].ID); err != nil {
				return nil, err
			}
			i++
		}
	}

	for i := range details {
		details[i].TransactionID = transactionID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}, nil
}
