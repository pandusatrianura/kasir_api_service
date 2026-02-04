package repository

import (
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
)

type HealthRepository struct {
	db *database.DB
}

type IHealthRepository interface {
	DB() error
}

func NewHealthRepository(db *database.DB) HealthRepository {
	return HealthRepository{db: db}
}

func (h *HealthRepository) DB() error {
	err := h.db.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}
