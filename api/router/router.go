package router

import (
	"net/http"

	categoriesHandler "github.com/pandusatrianura/kasir_api_service/internal/categories/delivery/http"
	healthHandler "github.com/pandusatrianura/kasir_api_service/internal/health/delivery/http"
	indexHandler "github.com/pandusatrianura/kasir_api_service/internal/index/delivery/http"
	productsHandler "github.com/pandusatrianura/kasir_api_service/internal/products/delivery/http"
	reportHandler "github.com/pandusatrianura/kasir_api_service/internal/reports/delivery/http"
	transactionsHandler "github.com/pandusatrianura/kasir_api_service/internal/transactions/delivery/http"
)

type Router struct {
	categories   *categoriesHandler.CategoryHandler
	products     *productsHandler.ProductHandler
	health       *healthHandler.HealthHandler
	transactions *transactionsHandler.TransactionHandler
	index        *indexHandler.IndexHandler
	report       *reportHandler.ReportHandler
}

func NewRouter(categoriesHandler *categoriesHandler.CategoryHandler, productHandler *productsHandler.ProductHandler,
	healthHandler *healthHandler.HealthHandler, transactionHandler *transactionsHandler.TransactionHandler,
	indexHandler *indexHandler.IndexHandler, reportHandler *reportHandler.ReportHandler) *Router {
	return &Router{
		categories:   categoriesHandler,
		products:     productHandler,
		health:       healthHandler,
		transactions: transactionHandler,
		index:        indexHandler,
		report:       reportHandler,
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
	r.HandleFunc("GET /transactions/health", h.transactions.API)
	r.HandleFunc("POST /transactions/checkout", h.transactions.Checkout)
	r.HandleFunc("GET /reports/health", h.report.API)
	r.HandleFunc("GET /reports/hari-ini", h.report.Today)
	r.HandleFunc("GET /reports", h.report.Report)
	r.HandleFunc("GET /docs", h.index.Docs)
	return r
}

func (h *Router) RegisterIndex() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("GET /", h.index.Index)
	return r
}
