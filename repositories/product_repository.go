package repositories

import (
	"KASIR-API/models"
	"database/sql"
	"errors"
)

type ProductRepository struct {
	db *sql.DB
}

type ProductWithCategory struct {
	Product  models.Product
	Category models.Category
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(name string) ([]ProductWithCategory, error) {
	query := `
	SELECT
		p.id,
		p.name,
		p.price,
		p.stock,
		p.category_id,
		c.id,
		c.name,
		c.description
	FROM products p
	JOIN categories c ON c.id = p.category_id
	`

	args := []interface{}{}

	if name != "" {
		query += "where p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ProductWithCategory

	for rows.Next() {
		var p models.Product
		var c models.Category

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&c.ID,
			&c.Name,
			&c.Description,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, ProductWithCategory{
			Product:  p,
			Category: c,
		})
	}

	return results, nil
}

func (repo *ProductRepository) GetByID(id int) (*ProductWithCategory, error) {

	query := `
	SELECT
		p.id,
		p.name,
		p.price,
		p.stock,
		p.category_id,
		c.id,
		c.name,
		c.description
	FROM products p
	JOIN categories c ON c.id = p.category_id
	WHERE p.id = $1
	`

	var p models.Product
	var c models.Category
	err := repo.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&c.ID,
		&c.Name,
		&c.Description,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &ProductWithCategory{
		Product:  p,
		Category: c,
	}, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) FindByIDWithCategory(id int) (*ProductWithCategory, error) {
	query := `
		SELECT 
			p.id, p.name, p.price, p.stock, p.category_id,
			c.id, c.name, c.description
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`

	var p models.Product
	var c models.Category

	err := repo.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&c.ID,
		&c.Name,
		&c.Description,
	)
	if err != nil {
		return nil, err
	}

	return &ProductWithCategory{
		Product:  p,
		Category: c,
	}, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}
