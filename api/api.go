// Package api Server represents an HTTP server with an address for listening to incoming requests.
package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pandusatrianura/kasir_api_service/api/middleware"
	route "github.com/pandusatrianura/kasir_api_service/api/router"
	categoryHandler "github.com/pandusatrianura/kasir_api_service/internal/categories/delivery/http"
	categoryRepository "github.com/pandusatrianura/kasir_api_service/internal/categories/repository"
	categoryService "github.com/pandusatrianura/kasir_api_service/internal/categories/service"
	healthHandler "github.com/pandusatrianura/kasir_api_service/internal/health/delivery/http"
	healthRepository "github.com/pandusatrianura/kasir_api_service/internal/health/repository"
	healthService "github.com/pandusatrianura/kasir_api_service/internal/health/service"
	indexHandler "github.com/pandusatrianura/kasir_api_service/internal/index/delivery/http"
	productHandler "github.com/pandusatrianura/kasir_api_service/internal/products/delivery/http"
	productRepository "github.com/pandusatrianura/kasir_api_service/internal/products/repository"
	productService "github.com/pandusatrianura/kasir_api_service/internal/products/service"
	reportHandler "github.com/pandusatrianura/kasir_api_service/internal/reports/delivery/http"
	reportRepository "github.com/pandusatrianura/kasir_api_service/internal/reports/repository"
	reportService "github.com/pandusatrianura/kasir_api_service/internal/reports/service"
	transactionsHandler "github.com/pandusatrianura/kasir_api_service/internal/transactions/delivery/http"
	transactionsRepository "github.com/pandusatrianura/kasir_api_service/internal/transactions/repository"
	transactionsService "github.com/pandusatrianura/kasir_api_service/internal/transactions/service"

	"github.com/go-chi/chi/v5"
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
)

type Server struct {
	addr string
	db   *database.DB
}

// NewAPIServer initializes and returns a new Server instance configured to listen to the specified address.
func NewAPIServer(addr string, db *database.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

// Run starts the server, initializes dependencies, registers routes, and listens for incoming HTTP requests.
func (s *Server) Run() error {

	categoriesRepo := categoryRepository.NewCategoryRepository(s.db)
	categoriesSvc := categoryService.NewCategoryService(categoriesRepo)
	categoriesHandle := categoryHandler.NewCategoryHandler(categoriesSvc)

	productsRepo := productRepository.NewProductRepository(s.db)
	productsSvc := productService.NewProductService(productsRepo, categoriesRepo)
	productsHandle := productHandler.NewProductHandler(productsSvc)

	healthRepo := healthRepository.NewHealthRepository(s.db)
	healthSvc := healthService.NewHealthService(healthRepo)
	healthHandle := healthHandler.NewHealthHandler(healthSvc)

	transactionsRepo := transactionsRepository.NewTransactionsRepository(s.db)
	transactionsSvc := transactionsService.NewTransactionsService(transactionsRepo)
	transactionsHandle := transactionsHandler.NewTransactionsHandler(transactionsSvc)

	reportsRepo := reportRepository.NewReportsRepository(s.db)
	reportsService := reportService.NewReportService(reportsRepo)
	reportsHandle := reportHandler.NewReportHandler(reportsService)

	indexHandle := indexHandler.NewIndexHandler()

	routers := route.NewRouter(categoriesHandle, productsHandle, healthHandle, transactionsHandle, indexHandle, reportsHandle)
	productRoute := routers.RegisterProductRoutes()
	indexRoutes := routers.RegisterIndexRoutes()
	docsRoutes := routers.RegisterDocsRoutes()
	categoryRoutes := routers.RegisterCategoriesRoutes()
	transactionRoutes := routers.RegisterTransactionRoutes()
	reportRoutes := routers.RegisterReportRoutes()

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware, middleware.ErrorHandlingMiddleware, middleware.CORS)
	r.Route("/api", func(r chi.Router) {
		r.Mount("/products", productRoute)
		r.Mount("/categories", categoryRoutes)
		r.Mount("/transactions", transactionRoutes)
		r.Mount("/reports", reportRoutes)
		r.Mount("/docs", docsRoutes)
	})
	r.Route("/", func(r chi.Router) {
		r.Mount("/", indexRoutes)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	addr := fmt.Sprintf("%s%s", "0.0.0.0", s.addr)
	log.Println("Starting server on", addr)
	return http.ListenAndServe(s.addr, r)
}
