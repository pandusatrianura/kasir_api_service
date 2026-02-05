package repository

import (
	"database/sql"
	"errors"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
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
	productRepo productRepository.IProductRepository
}

func NewTransactionsRepository(db *database.DB, productRepository productRepository.IProductRepository) ITransactionsRepository {
	return &TransactionsRepository{
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

	if detailProducts == nil || len(detailProducts) == 0 {
		return nil, errors.New(constants.ErrProductNotFound)
	}

	for _, product := range detailProducts {
		if product.ID == 0 {
			return nil, errors.New(constants.ErrProductNotFound)
		}

		if product.Stock < product.Quantity {
			return nil, errors.New(constants.ErrStockNotEnough)
		}

		if product.Stock == 0 {
			return nil, errors.New(constants.ErrStockEmpty)
		}

		subTotal = product.Price * product.Quantity
		totalAmount += subTotal

		checkoutProduct := entity.CheckoutProduct{
			ProductID:    product.ID,
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

	transactionID, checkoutProducts, err := t.createTransaction(totalAmount, checkoutProducts)
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

			if err != nil {
				return err
			}

			if errors.Is(sql.ErrNoRows, err) {
				return errors.New(constants.ErrProductNotFound)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		product.Quantity = request.Quantity
		products = append(products, product)
	}

	return products, nil
}

func (t *TransactionsRepository) createTransaction(totalAmount int, checkoutProducts []entity.CheckoutProduct) (int, []entity.CheckoutProduct, error) {
	var (
		query          string
		err            error
		lastInsertId   int64
		checkoutWithID []entity.CheckoutProduct
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

		checkoutWithID, err = t.createTransactionDetail(int(lastInsertId), checkoutProducts)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, nil, err
	}

	return int(lastInsertId), checkoutWithID, nil
}

func (t *TransactionsRepository) createTransactionDetail(transactionId int, checkoutProducts []entity.CheckoutProduct) ([]entity.CheckoutProduct, error) {
	var (
		query                  string
		err                    error
		lastInsertId           int64
		checkoutProductsWithID []entity.CheckoutProduct
	)

	query = "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"

	for _, product := range checkoutProducts {
		err = t.db.WithTx(func(tx *database.Tx) error {
			rows := t.db.QueryRow(query, transactionId, product.ProductID, product.Quantity, product.Subtotal, "now()", "now()")
			if rows.Error() != "" {
				return errors.New(rows.Error())
			}

			if err := rows.Scan(&lastInsertId); err != nil {
				return err
			}

			checkoutProductsWithID = append(checkoutProductsWithID, entity.CheckoutProduct{
				ProductID:           product.ProductID,
				TransactionDetailID: int(lastInsertId),
				TransactionID:       transactionId,
				Name:                product.Name,
				Quantity:            product.Quantity,
				Subtotal:            product.Subtotal,
				CategoryID:          product.CategoryID,
				CategoryName:        product.CategoryName,
			})

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return checkoutProductsWithID, nil
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
		return errors.New(constants.ErrProductNotFound)
	}

	return nil
}
