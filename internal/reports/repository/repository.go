package repository

import (
	"database/sql"
	"errors"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	productRepository "github.com/pandusatrianura/kasir_api_service/internal/products/repository"
	"github.com/pandusatrianura/kasir_api_service/internal/reports/entity"
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
)

type IReportsRepository interface {
	Report(startDate string, endDate string) ([]entity.ReportTransaction, error)
}

type ReportsRepository struct {
	db          *database.DB
	productRepo productRepository.ProductRepository
}

func NewReportsRepository(db *database.DB, productRepository productRepository.ProductRepository) ReportsRepository {
	return ReportsRepository{
		db:          db,
		productRepo: productRepository,
	}
}

func (r *ReportsRepository) Report(startDate string, endDate string) (*entity.ReportTransaction, error) {
	var (
		totalRevenue     int
		totalTransaction int
		soldsProduct     []entity.MostSoldProduct
		soldProduct      entity.MostSoldProduct
		query            string
		err              error
	)

	if startDate == "" || endDate == "" {
		return nil, errors.New(constants.ErrRequiredDate)
	}

	query = "SELECT SUM(total_amount) AS total_revenue, COUNT(id) AS total_transaction FROM transactions WHERE created_at BETWEEN $1 AND $2"

	err = r.db.WithStmt(query, func(stmt *database.Stmt) error {
		err = stmt.Query(func(rows *database.Rows) error {
			if err := rows.Scan(&totalRevenue, &totalTransaction); err != nil {
				return errors.New(constants.ErrTransactionNotFound)
			}

			return nil
		}, startDate, endDate)

		if err != nil {
			return errors.New(constants.ErrTransactionNotFound)
		}

		if errors.Is(sql.ErrNoRows, err) {
			return errors.New(constants.ErrTransactionNotFound)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	query = "SELECT c.name, SUM(a.quantity) AS sum_quantity FROM transaction_details a JOIN transactions b ON a.transaction_id = b.id JOIN products c ON a.product_id = c.id WHERE b.created_at >= $1 AND b.created_at < $2 GROUP BY a.product_id, c.name ORDER BY sum_quantity DESC;"

	err = r.db.WithStmt(query, func(stmt *database.Stmt) error {
		err = stmt.Query(func(rows *database.Rows) error {
			if err := rows.Scan(&soldProduct.Name, &soldProduct.QtySold); err != nil {
				return err
			}

			soldsProduct = append(soldsProduct, soldProduct)
			return nil
		}, startDate, endDate)

		if err != nil {
			return err
		}

		if errors.Is(sql.ErrNoRows, err) {
			return errors.New(constants.ErrTransactionNotFound)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	report := entity.ReportTransaction{
		TotalRevenue:      int64(totalRevenue),
		TotalTransactions: totalTransaction,
		MostSoldProduct:   soldsProduct,
	}

	return &report, nil
}
