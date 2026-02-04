package service

import (
	"github.com/pandusatrianura/kasir_api_service/internal/transactions/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/transactions/repository"
)

type ITransactionsService interface {
	Checkout([]entity.CheckoutRequest) (*entity.CheckoutResponse, error)
	API() entity.HealthCheck
}

type TransactionsService struct {
	transactionsRepository repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return TransactionsService{transactionsRepository: repo}
}

func (t *TransactionsService) API() entity.HealthCheck {
	return entity.HealthCheck{
		Name:      "Transactions/Checkout API",
		IsHealthy: true,
	}
}

func (t *TransactionsService) Checkout(requests []entity.CheckoutRequest) (*entity.CheckoutResponse, error) {
	response, err := t.transactionsRepository.Checkout(requests)
	if err != nil {
		return nil, err
	}

	return response, nil
}
