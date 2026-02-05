package service

import (
	"github.com/pandusatrianura/kasir_api_service/internal/reports/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/reports/repository"
)

type IReportService interface {
	Report(startDate string, endDate string) (*entity.ReportTransaction, error)
	API() entity.HealthCheck
}

type ReportService struct {
	transactionsRepository repository.IReportsRepository
}

func NewReportService(repo repository.IReportsRepository) IReportService {
	return &ReportService{transactionsRepository: repo}
}

func (t *ReportService) API() entity.HealthCheck {
	return entity.HealthCheck{
		Name:      "Reports API",
		IsHealthy: true,
	}
}

func (s *ReportService) Report(startDate string, endDate string) (*entity.ReportTransaction, error) {
	return s.transactionsRepository.Report(startDate, endDate)
}
