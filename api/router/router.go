package router

import (
	"fmt"
	"net/http"

	categoriesHandler "github.com/pandusatrianura/kasir_api_service/internal/categories/delivery/http"
	healthHandler "github.com/pandusatrianura/kasir_api_service/internal/health/delivery/http"
	productsHandler "github.com/pandusatrianura/kasir_api_service/internal/products/delivery/http"
	transactionsHandler "github.com/pandusatrianura/kasir_api_service/internal/transactions/delivery/http"
	"github.com/pandusatrianura/kasir_api_service/pkg/scalar"
)

type Router struct {
	categories   *categoriesHandler.CategoryHandler
	products     *productsHandler.ProductHandler
	health       *healthHandler.HealthHandler
	transactions *transactionsHandler.TransactionHandler
}

func NewRouter(categoriesHandler *categoriesHandler.CategoryHandler, productHandler *productsHandler.ProductHandler, healthHandler *healthHandler.HealthHandler, transactionHandler *transactionsHandler.TransactionHandler) *Router {
	return &Router{
		categories:   categoriesHandler,
		products:     productHandler,
		health:       healthHandler,
		transactions: transactionHandler,
	}
}

func (h *Router) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /health/service", h.health.API)
	r.HandleFunc("GET /health/db", h.health.DB)
	r.HandleFunc("GET /products/health", h.products.API)
	r.HandleFunc("POST /products", h.products.CreateProduct)
	r.HandleFunc("GET /products", h.products.GetAllProducts)
	r.HandleFunc("GET /products/{id}", h.products.GetProductByID)
	r.HandleFunc("PUT /products/{id}", h.products.UpdateProduct)
	r.HandleFunc("DELETE /products/{id}", h.products.DeleteProduct)
	r.HandleFunc("GET /categories/health", h.categories.API)
	r.HandleFunc("POST /categories", h.categories.CreateCategory)
	r.HandleFunc("GET /categories", h.categories.GetAllCategories)
	r.HandleFunc("GET /categories/{id}", h.categories.GetCategoryByID)
	r.HandleFunc("PUT /categories/{id}", h.categories.UpdateCategory)
	r.HandleFunc("DELETE /categories/{id}", h.categories.DeleteCategory)
	r.HandleFunc("POST /transactions/checkout", h.transactions.Checkout)
	r.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Test Kasir API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		_, err = fmt.Fprintln(w, htmlContent)
		if err != nil {
			return
		}
	})
	return r
}
