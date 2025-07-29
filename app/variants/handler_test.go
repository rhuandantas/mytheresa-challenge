package variants

import (
	"encoding/json"
	"errors"
	models2 "github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	repo_mock "github.com/mytheresa/go-hiring-challenge/app/repositories/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

var _ = Describe("Variant Handler", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *repo_mock.MockProductRepository
		handler  *Handler
		rr       *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = repo_mock.NewMockProductRepository(ctrl)
		handler = NewVariantHandler(mockRepo)
		rr = httptest.NewRecorder()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("returns product details with variants on success", func() {
		productID := uint64(1)
		product := &models2.Product{
			Code:     "P1",
			Price:    decimal.NewFromFloat(100.0),
			Category: models2.Category{Name: "Shoes"},
			Variants: []models2.Variant{
				{ID: 1, Name: "Red", Price: decimal.NewFromFloat(0)},
				{ID: 2, Name: "Blue", Price: decimal.NewFromFloat(120)},
			},
		}
		mockRepo.EXPECT().GetProductDetails(productID).Return(product, nil)

		req := httptest.NewRequest(http.MethodGet, "/catalog/"+strconv.FormatUint(productID, 10), nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))
		var resp Response
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp.Code).To(Equal("P1"))
		Expect(resp.Price).To(Equal(100.0))
		Expect(resp.Category).To(Equal("Shoes"))
		Expect(resp.Variants).To(HaveLen(2))
		variantPrice1, exact := resp.Variants[0].Price.Float64()
		Expect(exact).To(BeTrue())
		Expect(variantPrice1).To(Equal(100.0)) // inherited price
		variantPrice2, exact := resp.Variants[1].Price.Float64()
		Expect(exact).To(BeTrue())
		Expect(variantPrice2).To(Equal(120.0))
	})

	It("returns 404 if product not found", func() {
		productID := uint64(2)
		mockRepo.EXPECT().GetProductDetails(productID).Return(nil, gorm.ErrRecordNotFound)

		req := httptest.NewRequest(http.MethodGet, "/catalog/"+strconv.FormatUint(productID, 10), nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusNotFound))
	})

	It("returns 400 for invalid product ID", func() {
		req := httptest.NewRequest(http.MethodGet, "/catalog/abc", nil)
		handler.HandleGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusBadRequest))
	})

	It("returns 400 if product ID is missing", func() {
		req := httptest.NewRequest(http.MethodGet, "/catalog/", nil)
		handler.HandleGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusBadRequest))
	})

	It("returns 500 on repository error", func() {
		productID := uint64(3)
		mockRepo.EXPECT().GetProductDetails(productID).Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/catalog/"+strconv.FormatUint(productID, 10), nil)
		handler.HandleGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusInternalServerError))
	})
})

func TestTaskAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Variant API Suite")
}
