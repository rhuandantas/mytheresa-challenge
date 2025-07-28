package catalog

import (
	"encoding/json"
	"github.com/mytheresa/go-hiring-challenge/app/repositories"
	"net/http"
	"strconv"
)

type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category,omitempty"`
}

type Handler struct {
	repo repositories.ProductRepository
}

func NewCatalogHandler(r repositories.ProductRepository) *Handler {
	return &Handler{
		repo: r,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse offset and limit
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	products, total, err := h.repo.GetAllProducts(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	respProducts := make([]Product, len(products))
	for i, p := range products {
		respProducts[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: respProducts,
		Total:    total,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
