package catalog

import (
	"encoding/json"
	"errors"
	repo_mock "github.com/mytheresa/go-hiring-challenge/app/repositories/mocks"
	models2 "github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"testing"

	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Catalog Handler", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *repo_mock.MockProductRepository
		handler  *Handler
		rr       *httptest.ResponseRecorder
		req      *http.Request
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = repo_mock.NewMockProductRepository(ctrl)
		handler = NewCatalogHandler(mockRepo)
		rr = httptest.NewRecorder()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("returns products and total on success", func() {
		products := []models2.Product{
			{
				Code:  "P1",
				Price: decimal.NewFromFloat(1),
				Category: models2.Category{
					Name: "Category1",
				},
			},
		}
		mockRepo.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(products, int64(len(products)), nil)

		req = httptest.NewRequest(http.MethodGet, "/catalog", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))

		var resp Response
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp.Total).To(Equal(int64(1)))
		Expect(resp.Products).To(HaveLen(1))
		Expect(resp.Products[0].Code).To(Equal("P1"))
	})

	It("returns error on repository failure", func() {
		mockRepo.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, int64(0), errors.New("repository failure"))

		req = httptest.NewRequest(http.MethodGet, "/catalog", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		var resp map[string]string
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp["error"]).To(Equal("failed to fetch products"))
	})
	It("returns empty response when no products found", func() {
		mockRepo.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, int64(0), nil)

		req = httptest.NewRequest(http.MethodGet, "/catalog", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))

		var resp Response
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp.Total).To(Equal(int64(0)))
		Expect(resp.Products).To(BeEmpty())
	})
	It("should send limit equals 10 when not specified ", func() {
		mockRepo.EXPECT().GetAllProducts(0, 10, "", 0.0).
			Return(nil, int64(0), nil)

		req = httptest.NewRequest(http.MethodGet, "/catalog?offset=0", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))

		var resp Response
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp.Total).To(Equal(int64(0)))
		Expect(resp.Products).To(BeEmpty())
	})
	It("should send minimum limit when passing less than 1", func() {
		mockRepo.EXPECT().GetAllProducts(0, 1, "", 0.0).
			Return(nil, int64(0), nil)

		req = httptest.NewRequest(http.MethodGet, "/catalog?offset=0&limit=0", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))

		var resp Response
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp.Total).To(Equal(int64(0)))
		Expect(resp.Products).To(BeEmpty())
	})
	It("should send limit 100 when passing more than 100", func() {
		mockRepo.EXPECT().GetAllProducts(0, 100, "", 0.0).
			Return(nil, int64(0), nil)

		req = httptest.NewRequest(http.MethodGet, "/catalog?offset=0&limit=101", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))

		var resp Response
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp.Total).To(Equal(int64(0)))
		Expect(resp.Products).To(BeEmpty())
	})
	It("should return error when passing invalid limit", func() {
		req = httptest.NewRequest(http.MethodGet, "/catalog?offset=0&limit=invalid", nil)
		handler.HandleGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusBadRequest))
		var resp map[string]string
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp["error"]).To(Equal("invalid request parameters"))
	})
	It("should return error when passing invalid offset", func() {
		req = httptest.NewRequest(http.MethodGet, "/catalog?offset=invalid&limit=10", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusBadRequest))
		var resp map[string]string
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp["error"]).To(Equal("invalid request parameters"))
	})
	It("should send offset equals 0 when less than 0", func() {
		mockRepo.EXPECT().GetAllProducts(0, 10, "", 0.0).
			Return(nil, int64(0), nil)

		req = httptest.NewRequest(http.MethodGet, "/catalog?offset=-1&limit=10", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))

		var resp Response
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp.Total).To(Equal(int64(0)))
		Expect(resp.Products).To(BeEmpty())
	})
	It("should return error when passing an invalid price_lt", func() {
		req = httptest.NewRequest(http.MethodGet, "/catalog?price_lt=invalid", nil)
		handler.HandleGet(rr, req)

		Expect(rr.Code).To(Equal(http.StatusBadRequest))
		var resp map[string]string
		err := json.NewDecoder(rr.Body).Decode(&resp)
		Expect(err).To(BeNil())
		Expect(resp["error"]).To(Equal("invalid request parameters"))
	})
})

func TestTaskAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Catalog API Suite")
}
