package dto

import "time"

type BestProduct struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type ReportResponse struct {
	TotalRevenue   int         `json:"total_revenue"`
	TotalTransaksi int         `json:"total_transaksi"`
	ProdukTerlaris BestProduct `json:"produk_terlaris"`
}

type TransactionRangeItem struct {
	ID          int       `json:"id"`
	TotalAmount int       `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReportRangeResponse struct {
	Transactions []TransactionRangeItem `json:"transactions"`
}
