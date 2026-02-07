package repositories

import (
	"KASIR-API/dto"
	"database/sql"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetReport(startUTC, endUTC time.Time) (int, int, string, int, error) {

	var totalRevenue sql.NullInt64
	var totalTransaksi int

	// Total revenue + total transaksi
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(total_amount), 0),
			COUNT(*)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, startUTC, endUTC).Scan(&totalRevenue, &totalTransaksi)

	if err != nil {
		return 0, 0, "", 0, err
	}

	// Produk terlaris untuk end point hari-ini
	var nama string
	var qty sql.NullInt64

	err = r.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as total_qty
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startUTC, endUTC).Scan(&nama, &qty)

	if err == sql.ErrNoRows {
		return int(totalRevenue.Int64), totalTransaksi, "", 0, nil
	}

	if err != nil {
		return 0, 0, "", 0, err
	}

	return int(totalRevenue.Int64), totalTransaksi, nama, int(qty.Int64), nil
}

// optional change range date report
func (r *ReportRepository) GetTransactionsByRange(startUTC, endUTC time.Time) ([]dto.TransactionRangeItem, error) {

	rows, err := r.db.Query(`
		SELECT id, total_amount, created_at
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at DESC
	`, startUTC, endUTC)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []dto.TransactionRangeItem

	for rows.Next() {
		var t dto.TransactionRangeItem

		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}
