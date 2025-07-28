package repositories

import (
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllProducts(offset, limit int) ([]models.Product, int64, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) ProductRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	// Count total products
	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated products
	if err := r.db.Preload("Variants").Preload("Category").
		Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
