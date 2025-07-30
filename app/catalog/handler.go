package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/repositories"
	"github.com/mytheresa/go-hiring-challenge/models"
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

// HandleGet handles the request to get a list of products based on query parameters.
// It returns a paginated list of products with optional filters for category and price.
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	params, err := parseRequestParams(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request parameters")
		return
	}

	products, total, err := h.repo.GetAllProducts(params.Offset, params.Limit, params.Category, params.PriceLt)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	response := h.mapResponse(products, total)

	api.OKResponse(w, response)
}

// mapResponse converts the list of products and total count into the Response format.
func (h *Handler) mapResponse(products []models.Product, total int64) Response {
	respProducts := make([]Product, len(products))
	for i, p := range products {
		price, ok := p.Price.Float64()
		if !ok {
			price = p.Price.InexactFloat64()
		}

		respProducts[i] = Product{
			Code:     p.Code,
			Price:    price,
			Category: p.Category.Name,
		}
	}

	response := Response{
		Products: respProducts,
		Total:    total,
	}
	return response
}

// parseRequestParams extracts and validates request parameters from the HTTP request.
// It returns a RequestParams struct or an error if the parameters are invalid.
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
