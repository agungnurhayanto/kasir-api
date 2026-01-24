package storage

import "KASIR-API/models"

var Categories = []models.Category{
	{ID: 1, Name: "Komputer", Description: "Peralatan Komputer dan Aksesoris"},
	{ID: 2, Name: "Furniture", Description: "Perabotan rumah dan kantor"},
	{ID: 3, Name: "Fashion", Description: "Baju dan aksesoris fashion"},
}

var Produk = []models.Produk{
	{ID: 1, Nama: "Power Supplay", Harga: 275000, Stok: 10},
	{ID: 2, Nama: "Motherboard", Harga: 1250000, Stok: 15},
	{ID: 3, Nama: "Processor", Harga: 750000, Stok: 20},
}
