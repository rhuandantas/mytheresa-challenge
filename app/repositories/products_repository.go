package repositories

import (
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllProducts(page, limit int, category string, priceLt float64) ([]models.Product, int64, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int, category string, priceLt float64) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Preload("Variants").Preload("Category")

	if category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").Where("categories.code = ?", category)
	}
	if priceLt > 0 {
		query = query.Where("products.price < ?", priceLt)
	}

	// Count total with filters
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
