package service

import (
	"github.com/pandusatrianura/kasir_api_service/internal/health/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/health/repository"
)

type HealthService struct {
	healthRepository repository.HealthRepository
}

type IHealthService interface {
	API() entity.HealthCheck
	DB() (entity.HealthCheck, error)
}

func NewHealthService(healthRepo repository.HealthRepository) HealthService {
	return HealthService{healthRepository: healthRepo}
}

func (h *HealthService) API() entity.HealthCheck {

	return entity.HealthCheck{
		Name:      "Connection to Kasir API",
		IsHealthy: true,
	}
}

func (h *HealthService) DB() (entity.HealthCheck, error) {
	err := h.healthRepository.DB()
	if err != nil {
		return entity.HealthCheck{}, err
	}

	return entity.HealthCheck{
		Name:      "Connection to Kasir Database",
		IsHealthy: true,
	}, nil
}
