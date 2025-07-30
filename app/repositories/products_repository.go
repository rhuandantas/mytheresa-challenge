package repositories

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	// GetAllProducts retrieves products with optional filters for category and price.
	// It returns a slice of products, total count, and an error if any.
	// page is the offset for pagination, limit is the maximum number of products to return.
	// category is the category code to filter products, and priceLt is the maximum price to filter products.
	GetAllProducts(page, limit int, category string, priceLt float64) ([]models.Product, int64, error)
	// GetProductDetails retrieves a product by ID, including its variants and category.
	// It returns the product details or an error if not found.
	// id is the product ID to retrieve.
	GetProductDetails(id uint64) (*models.Product, error)
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

	query := r.db.Model(&models.Product{}).Preload("Category")

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

func (r *ProductsRepository) GetProductDetails(id uint64) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Variants").Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
