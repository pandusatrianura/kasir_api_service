package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/pandusatrianura/kasir_api_service/api/middleware"
	categoriesHandler "github.com/pandusatrianura/kasir_api_service/internal/categories/delivery/http"
	healthHandler "github.com/pandusatrianura/kasir_api_service/internal/health/delivery/http"
	indexHandler "github.com/pandusatrianura/kasir_api_service/internal/index/delivery/http"
	productsHandler "github.com/pandusatrianura/kasir_api_service/internal/products/delivery/http"
	reportHandler "github.com/pandusatrianura/kasir_api_service/internal/reports/delivery/http"
	transactionsHandler "github.com/pandusatrianura/kasir_api_service/internal/transactions/delivery/http"
	userHandler "github.com/pandusatrianura/kasir_api_service/internal/users/delivery/http"
)

type Router struct {
	categories   *categoriesHandler.CategoryHandler
	products     *productsHandler.ProductHandler
	health       *healthHandler.HealthHandler
	transactions *transactionsHandler.TransactionHandler
	index        *indexHandler.IndexHandler
	report       *reportHandler.ReportHandler
	user         *userHandler.UserHandler
}

func NewRouter(categoriesHandler *categoriesHandler.CategoryHandler, productHandler *productsHandler.ProductHandler,
	healthHandler *healthHandler.HealthHandler, transactionHandler *transactionsHandler.TransactionHandler,
	indexHandler *indexHandler.IndexHandler, reportHandler *reportHandler.ReportHandler, userHandler *userHandler.UserHandler) *Router {
	return &Router{
		categories:   categoriesHandler,
		products:     productHandler,
		health:       healthHandler,
		transactions: transactionHandler,
		index:        indexHandler,
		report:       reportHandler,
		user:         userHandler,
	}
}

func (h *Router) RegisterProductRoutes() chi.Router {
	r := chi.NewRouter()
	products := h.products
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth, middleware.JWTAuthMiddleware)
		r.Post("/", products.CreateProduct)
		r.Put("/{id}", products.UpdateProduct)
		r.Delete("/{id}", products.DeleteProduct)
		r.Get("/{id}", products.GetProductByID)
	})
	r.Get("/", products.GetAllProducts)
	r.Get("/health", products.API)
	return r
}

func (h *Router) RegisterCategoriesRoutes() chi.Router {
	r := chi.NewRouter()
	categories := h.categories
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth, middleware.JWTAuthMiddleware)
		r.Post("/", categories.CreateCategory)
		r.Get("/{id}", categories.GetCategoryByID)
		r.Put("/{id}", categories.UpdateCategory)
		r.Delete("/{id}", categories.DeleteCategory)
	})

	r.Get("/", categories.GetAllCategories)
	r.Get("/health", categories.API)
	return r
}

func (h *Router) RegisterTransactionRoutes() chi.Router {
	r := chi.NewRouter()
	transactions := h.transactions
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth, middleware.JWTAuthMiddleware)
		r.Post("/checkout", transactions.Checkout)
	})
	r.Get("/health", transactions.API)
	return r
}

func (h *Router) RegisterReportRoutes() chi.Router {
	r := chi.NewRouter()
	report := h.report
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth, middleware.JWTAuthMiddleware)
		r.Get("/hari-ini", report.Today)
		r.Get("/", report.Report)
	})
	r.Get("/health", report.API)
	return r
}

func (h *Router) RegisterHealthRoutes() chi.Router {
	r := chi.NewRouter()
	health := h.health
	r.Get("/service", health.API)
	r.Get("/db", health.DB)
	return r
}

func (h *Router) RegisterDocsRoutes() chi.Router {
	r := chi.NewRouter()
	index := h.index
	r.Get("/", index.Docs)
	return r
}

func (h *Router) RegisterIndexRoutes() chi.Router {
	r := chi.NewRouter()
	index := h.index
	r.Get("/", index.Index)
	return r
}

func (h *Router) RegisterUserRoutes() chi.Router {
	r := chi.NewRouter()
	user := h.user
	r.Post("/login", user.Login)
	return r
}
