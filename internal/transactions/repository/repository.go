package repository

import (
	"errors"

	productEntity "github.com/pandusatrianura/kasir_api_service/internal/products/entity"
	productRepository "github.com/pandusatrianura/kasir_api_service/internal/products/repository"
	"github.com/pandusatrianura/kasir_api_service/internal/transactions/entity"
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
)

type ITransactionsRepository interface {
	Checkout([]entity.CheckoutRequest) (*entity.CheckoutResponse, error)
}

type TransactionsRepository struct {
	db          *database.DB
	productRepo productRepository.ProductRepository
}

func NewTransactionsRepository(db *database.DB, productRepository productRepository.ProductRepository) TransactionsRepository {
	return TransactionsRepository{
		db:          db,
		productRepo: productRepository,
	}
}

func (t *TransactionsRepository) Checkout(requests []entity.CheckoutRequest) (*entity.CheckoutResponse, error) {
	var (
		totalAmount      int
		subTotal         int
		checkoutProducts []entity.CheckoutProduct
		detailProducts   []entity.CheckoutProductDetail
		updateProducts   []productEntity.Product
		err              error
	)

	totalAmount = 0
	subTotal = 0

	detailProducts, err = t.getDetailProductByID(requests)
	if err != nil {
		return nil, err
	}

	if detailProducts == nil {
		return nil, errors.New("product not found")
	}

	for _, product := range detailProducts {
		subTotal = product.Price * product.Quantity
		totalAmount += subTotal

		checkoutProduct := entity.CheckoutProduct{
			ID:           product.ID,
			Name:         product.Name,
			Quantity:     product.Quantity,
			Subtotal:     subTotal,
			CategoryID:   product.CategoryID,
			CategoryName: product.CategoryName,
		}

		updateProduct := productEntity.Product{
			ID:    product.ID,
			Stock: product.Stock - product.Quantity,
		}

		updateProducts = append(updateProducts, updateProduct)
		checkoutProducts = append(checkoutProducts, checkoutProduct)
	}

	transactionID, err := t.createTransaction(totalAmount, checkoutProducts)
	if err != nil {
		return nil, err
	}

	err = t.updateProductStock(updateProducts)
	if err != nil {
		return nil, err
	}

	response := entity.CheckoutResponse{
		Transaction: entity.Transaction{
			ID:          transactionID,
			TotalAmount: totalAmount,
		},
		CheckoutProducts: checkoutProducts,
	}

	return &response, nil
}

func (t *TransactionsRepository) getDetailProductByID(requests []entity.CheckoutRequest) ([]entity.CheckoutProductDetail, error) {
	var (
		product  entity.CheckoutProductDetail
		products []entity.CheckoutProductDetail
		err      error
		query    string
	)

	query = "SELECT products.id, products.name, products.price, products.stock, categories.id as category_id, categories.name as category_name FROM products JOIN categories ON products.category_id = categories.id WHERE products.id = $1"

	for _, request := range requests {
		err = t.db.WithStmt(query, func(stmt *database.Stmt) error {
			err = stmt.Query(func(rows *database.Rows) error {
				if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.CategoryName); err != nil {
					return err
				}
				return nil
			}, request.ProductID)

			return err
		})

		if err != nil {
			return nil, err
		}

		product.Quantity = request.Quantity
		products = append(products, product)
	}

	return products, nil
}

func (t *TransactionsRepository) createTransaction(totalAmount int, checkoutProducts []entity.CheckoutProduct) (int, error) {
	var (
		query        string
		err          error
		lastInsertId int64
	)

	query = "INSERT INTO transactions (total_amount, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id;"

	err = t.db.WithTx(func(tx *database.Tx) error {
		rows := t.db.QueryRow(query, totalAmount, "now()", "now()")
		if rows.Error() != "" {
			return errors.New(rows.Error())
		}

		if err := rows.Scan(&lastInsertId); err != nil {
			return err
		}

		err = t.createTransactionDetail(int(lastInsertId), checkoutProducts)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}

func (t *TransactionsRepository) createTransactionDetail(transactionId int, checkoutProducts []entity.CheckoutProduct) error {
	var (
		query string
		err   error
	)

	query = "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	for _, product := range checkoutProducts {
		err = t.db.WithTx(func(tx *database.Tx) error {
			err = tx.WithStmt(query, func(stmt *database.Stmt) error {
				_, err = stmt.Exec(transactionId, product.ID, product.Quantity, product.Subtotal, "now()", "now()")
				return err
			})

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TransactionsRepository) updateProductStock(updateProducts []productEntity.Product) error {
	var (
		err error
	)

	if len(updateProducts) > 0 {
		for _, product := range updateProducts {
			err = t.productRepo.UpdateStockByProductID(&product)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("product not found")
	}

	return nil
}
