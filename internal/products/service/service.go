package service

import (
	"errors"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	categoryRepository "github.com/pandusatrianura/kasir_api_service/internal/categories/repository"
	"github.com/pandusatrianura/kasir_api_service/internal/products/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/products/repository"
)

type ProductService struct {
	productRepository  repository.IProductRepository
	categoryRepository categoryRepository.ICategoryRepository
}

type IProductService interface {
	CreateProduct(product *entity.RequestProduct) error
	UpdateProduct(id int64, product *entity.RequestProduct) error
	DeleteProduct(id int64) error
	GetProductByID(id int64) (*entity.ResponseProductWithCategories, error)
	GetAllProducts(name string) ([]entity.ResponseProductWithCategories, error)
	API() entity.HealthCheck
}

func NewProductService(productRepository repository.IProductRepository, categoryRepository categoryRepository.ICategoryRepository) IProductService {
	return &ProductService{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
	}
}

func (s *ProductService) API() entity.HealthCheck {
	return entity.HealthCheck{
		Name:      "Products API",
		IsHealthy: true,
	}
}

func (s *ProductService) CreateProduct(requestProduct *entity.RequestProduct) error {
	_, err := s.categoryRepository.GetCategoryByID(int64(requestProduct.CategoryID))
	if err != nil {
		return errors.New("category not found")
	}

	product := &entity.Product{
		Name:       requestProduct.Name,
		Price:      requestProduct.Price,
		Stock:      requestProduct.Stock,
		CategoryID: requestProduct.CategoryID,
	}

	return s.productRepository.CreateProduct(product)
}

func (s *ProductService) UpdateProduct(id int64, requestProduct *entity.RequestProduct) error {
	_, err := s.productRepository.GetProductByID(id)
	if err != nil {
		return errors.New(constants.ErrProductNotFound)
	}

	_, err = s.categoryRepository.GetCategoryByID(int64(requestProduct.CategoryID))
	if err != nil {
		return errors.New(constants.ErrCategoryNotFound)
	}

	product := &entity.Product{
		Name:       requestProduct.Name,
		Price:      requestProduct.Price,
		Stock:      requestProduct.Stock,
		CategoryID: requestProduct.CategoryID,
	}

	return s.productRepository.UpdateProduct(id, product)
}

func (s *ProductService) DeleteProduct(id int64) error {
	_, err := s.productRepository.GetProductByID(id)
	if err != nil {
		return errors.New(constants.ErrProductNotFound)
	}

	return s.productRepository.DeleteProduct(id)
}

func (s *ProductService) GetProductByID(id int64) (*entity.ResponseProductWithCategories, error) {
	result, err := s.productRepository.GetProductByID(id)
	return result, err
}

func (s *ProductService) GetAllProducts(name string) ([]entity.ResponseProductWithCategories, error) {
	return s.productRepository.GetAllProducts(name)
}
