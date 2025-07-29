package variants

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/repositories"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Response struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category string           `json:"category,omitempty"`
	Variants []models.Variant `json:"variants,omitempty"`
}

type Handler struct {
	repo repositories.ProductRepository
}

func NewVariantHandler(r repositories.ProductRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

// HandleGet handles the request to get product details by ID
// It returns the product code, price, category, and variants.
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from URL parameters
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		api.ErrorResponse(w, http.StatusBadRequest, "product ID is required")
		return
	}

	// parse to uint
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := h.repo.GetProductDetails(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ErrorResponse(w, http.StatusNotFound, "product not found")
			return
		}

		api.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch product details")
		return
	}

	// Inherit product price for variants without a specific price
	for i := range product.Variants {
		if product.Variants[i].Price.IsZero() {
			product.Variants[i].Price = product.Price
		}
	}

	// Convert price to float64, handling both exact and inexact cases
	price, ok := product.Price.Float64()
	if !ok {
		price = product.Price.InexactFloat64()
	}

	res := Response{
		Code:     product.Code,
		Price:    price,
		Category: product.Category.Name,
		Variants: product.Variants,
	}

	api.OKResponse(w, res)
}
