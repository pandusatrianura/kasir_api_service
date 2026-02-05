package service

import (
	"errors"

	"github.com/pandusatrianura/kasir_api_service/internal/categories/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/categories/repository"
)

type CategoryService struct {
	categoryRepository repository.ICategoryRepository
}

type ICategoryService interface {
	CreateCategory(requestCategory *entity.RequestCategory) error
	UpdateCategory(id int64, requestCategory *entity.RequestCategory) error
	DeleteCategory(id int64) error
	GetCategoryByID(id int64) (*entity.ResponseCategory, error)
	GetAllCategories() ([]entity.ResponseCategory, error)
	API() entity.HealthCheck
}

func NewCategoryService(categoryRepository repository.ICategoryRepository) ICategoryService {
	return &CategoryService{categoryRepository: categoryRepository}
}

func (s *CategoryService) API() entity.HealthCheck {
	return entity.HealthCheck{
		Name:      "Categories API",
		IsHealthy: true,
	}
}

func (s *CategoryService) CreateCategory(requestCategory *entity.RequestCategory) error {
	category := &entity.Category{
		Name:        requestCategory.Name,
		Description: requestCategory.Description,
	}
	return s.categoryRepository.CreateCategory(category)
}

func (s *CategoryService) UpdateCategory(id int64, requestCategory *entity.RequestCategory) error {
	_, err := s.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	category := &entity.Category{
		Name:        requestCategory.Name,
		Description: requestCategory.Description,
	}
	return s.categoryRepository.UpdateCategory(id, category)
}

func (s *CategoryService) DeleteCategory(id int64) error {
	_, err := s.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	return s.categoryRepository.DeleteCategory(id)
}

func (s *CategoryService) GetCategoryByID(id int64) (*entity.ResponseCategory, error) {
	return s.categoryRepository.GetCategoryByID(id)
}

func (s *CategoryService) GetAllCategories() ([]entity.ResponseCategory, error) {
	return s.categoryRepository.GetAllCategories()
}
