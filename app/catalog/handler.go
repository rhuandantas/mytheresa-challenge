package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/repositories"
	"net/http"
	"strconv"
)

type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

type RequestParams struct {
	Offset   int
	Limit    int
	Category string
	PriceLt  float64
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
	params, err := parseRequestParams(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request parameters")
		return
	}

	// Validate request parameters
	products, total, err := h.repo.GetAllProducts(params.Offset, params.Limit, params.Category, params.PriceLt)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch products")
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

	api.OKResponse(w, response)
}
func parseRequestParams(r *http.Request) (*RequestParams, error) {
	var (
		offset = 0
		limit  = 10
		err    error
	)

	if param := r.URL.Query().Get("offset"); param != "" {
		if offset, err = strconv.Atoi(param); err != nil {
			return nil, err
		}
		if offset < 0 {
			offset = 0
		}
	}

	if param := r.URL.Query().Get("limit"); param != "" {
		if limit, err = strconv.Atoi(param); err != nil {
			return nil, err
		}
		if limit < 1 {
			limit = 1 // Minimum limit enforced
		}
		if limit > 100 {
			limit = 100 // Maximum limit enforced
		}
	}

	category := r.URL.Query().Get("category")
	priceLt := 0.0
	if v := r.URL.Query().Get("price_lt"); v != "" {
		if priceLt, err = strconv.ParseFloat(v, 64); err != nil {
			return nil, err
		}
	}

	return &RequestParams{
		Offset:   offset,
		Limit:    limit,
		Category: category,
		PriceLt:  priceLt,
	}, nil
}
