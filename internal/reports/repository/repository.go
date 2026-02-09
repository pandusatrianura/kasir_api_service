package repository

import (
	"database/sql"
	"errors"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/internal/reports/entity"
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
)

type IReportsRepository interface {
	Report(startDate string, endDate string) (*entity.ReportTransaction, error)
}

type ReportsRepository struct {
	db *database.DB
}

func NewReportsRepository(db *database.DB) IReportsRepository {
	return &ReportsRepository{
		db: db,
	}
}

func (r *ReportsRepository) Report(startDate string, endDate string) (*entity.ReportTransaction, error) {
	var (
		totalRevenue     int
		totalTransaction int
		soldsProduct     []entity.MostSoldProduct
		err              error
	)

	soldsProduct = make([]entity.MostSoldProduct, 0)

	if startDate == "" || endDate == "" {
		return nil, errors.New(constants.ErrRequiredDate)
	}

	totalRevenue, totalTransaction, err = r.getRevenueAndTransaction(startDate, endDate)
	if err != nil {
		return nil, err
	}

	soldsProduct, err = r.getMostSoldProduct(startDate, endDate)
	if err != nil {
		return nil, err
	}

	soldsProduct = r.findMostSoldProducts(soldsProduct)
	if soldsProduct == nil {
		soldsProduct = []entity.MostSoldProduct{}
	}

	report := entity.ReportTransaction{
		TotalRevenue:      int64(totalRevenue),
		TotalTransactions: totalTransaction,
		MostSoldProduct:   soldsProduct,
	}

	return &report, nil
}

func (r *ReportsRepository) getRevenueAndTransaction(startDate string, endDate string) (int, int, error) {
	var (
		totalRevenue     int
		totalTransaction int
		query            string
		err              error
	)

	if startDate == "" || endDate == "" {
		return 0, 0, errors.New(constants.ErrRequiredDate)
	}

	query = "SELECT SUM(total_amount) AS total_revenue, COUNT(id) AS total_transaction FROM transactions WHERE created_at BETWEEN $1 AND $2"

	err = r.db.WithStmt(query, func(stmt *database.Stmt) error {
		scanFn := func(rows *database.Rows) error {
			return rows.Scan(&totalRevenue, &totalTransaction)
		}

		err = stmt.Query(scanFn, startDate, endDate)

		if err != nil {
			return errors.New(constants.ErrTransactionNotFound)
		}

		if errors.Is(sql.ErrNoRows, err) {
			return errors.New(constants.ErrTransactionNotFound)
		}

		return nil
	})

	if err != nil {
		return 0, 0, err
	}

	return totalRevenue, totalTransaction, nil
}

func (r *ReportsRepository) getMostSoldProduct(startDate string, endDate string) ([]entity.MostSoldProduct, error) {
	var (
		soldsProduct []entity.MostSoldProduct
		soldProduct  entity.MostSoldProduct
		query        string
		err          error
	)

	soldsProduct = make([]entity.MostSoldProduct, 0)

	if startDate == "" || endDate == "" {
		return nil, errors.New(constants.ErrRequiredDate)
	}

	query = "SELECT c.name, SUM(a.quantity) AS sum_quantity FROM transaction_details a JOIN transactions b ON a.transaction_id = b.id JOIN products c ON a.product_id = c.id WHERE b.created_at >= $1 AND b.created_at < $2 GROUP BY a.product_id, c.name ORDER BY sum_quantity DESC;"

	err = r.db.WithStmt(query, func(stmt *database.Stmt) error {
		scanFn := func(rows *database.Rows) error {
			if err := rows.Scan(&soldProduct.Name, &soldProduct.QtySold); err != nil {
				return err
			}

			soldsProduct = append(soldsProduct, soldProduct)
			return nil
		}

		err = stmt.Query(scanFn, startDate, endDate)
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

	return soldsProduct, nil
}

func (r *ReportsRepository) findMostSoldProducts(products []entity.MostSoldProduct) []entity.MostSoldProduct {
	var (
		prods []entity.MostSoldProduct
	)

	prods = make([]entity.MostSoldProduct, 0)

	if len(products) == 0 {
		return nil
	}

	limitMaxProdSold := r.limitMaxProdSold(products)

	for i := 0; i < len(products); i++ {
		if products[i].QtySold > limitMaxProdSold {
			prod := products[i]
			prods = append(prods, prod)
		} else if products[i].QtySold == limitMaxProdSold {
			prods = append(prods, products[i])
		}
	}

	return prods
}

func (r *ReportsRepository) limitMaxProdSold(products []entity.MostSoldProduct) int {
	maxObj := products[0]

	for _, something := range products {
		if something.QtySold > maxObj.QtySold {
			maxObj = something
		}
	}

	return maxObj.QtySold
}
