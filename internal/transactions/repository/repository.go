package repository

import (
	"database/sql"
	"errors"
	"fmt"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/internal/transactions/entity"
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
)

type ITransactionsRepository interface {
	Checkout([]entity.CheckoutRequest) (*entity.CheckoutResponse, error)
}

type TransactionsRepository struct {
	db *database.DB
}

func NewTransactionsRepository(db *database.DB) ITransactionsRepository {
	return &TransactionsRepository{
		db: db,
	}
}

func (t *TransactionsRepository) Checkout(requests []entity.CheckoutRequest) (*entity.CheckoutResponse, error) {
	var (
		totalAmount      int
		subTotal         int
		checkoutProducts []entity.CheckoutProduct
		detailProducts   []entity.CheckoutProductDetail
		updateProducts   []entity.UpdatedProduct
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

		updateProduct := entity.UpdatedProduct{
			ID:    product.ID,
			Stock: product.Stock - product.Quantity,
		}

		updateProducts = append(updateProducts, updateProduct)
		checkoutProducts = append(checkoutProducts, checkoutProduct)
	}

	transactionID, checkoutProducts, err := t.createTransaction(totalAmount, checkoutProducts, updateProducts)
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

	products = make([]entity.CheckoutProductDetail, 0)

	query = "SELECT products.id, products.name, products.price, products.stock, categories.id as category_id, categories.name as category_name FROM products JOIN categories ON products.category_id = categories.id WHERE products.id = $1"

	for _, request := range requests {
		err = t.db.WithStmt(query, func(stmt *database.Stmt) error {
			scanFn := func(rows *database.Rows) error {
				return rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.CategoryName)
			}

			err = stmt.Query(scanFn, request.ProductID)

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

func (t *TransactionsRepository) createTransaction(totalAmount int, checkoutProducts []entity.CheckoutProduct, updateProducts []entity.UpdatedProduct) (int, []entity.CheckoutProduct, error) {
	var (
		query          string
		err            error
		lastInsertId   int64
		checkoutWithID []entity.CheckoutProduct
	)

	checkoutWithID = make([]entity.CheckoutProduct, 0)

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

		err = t.updateProductsStock(updateProducts)
		if err != nil {
			return err
		}

		checkoutWithID = t.getTransactionsDetailByTransactionID(int(lastInsertId))
		if checkoutWithID == nil {
			return errors.New(constants.ErrProductNotFound)
		}

		return nil
	})

	if err != nil {
		return 0, nil, err
	}

	return int(lastInsertId), checkoutWithID, nil
}

func (t *TransactionsRepository) createTransactionDetail(transactionId int, checkoutProducts []entity.CheckoutProduct) error {
	var (
		query string
		err   error
		args  []interface{}
	)

	query = "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal, created_at, updated_at) VALUES "

	numFields := 6
	for i, product := range checkoutProducts {
		p := i * numFields
		query = fmt.Sprintf("%s ($%d, $%d, $%d, $%d, $%d, $%d)", query, p+1, p+2, p+3, p+4, p+5, p+6)
		if i < len(checkoutProducts)-1 {
			query += ","
		}
		args = append(args, transactionId, product.ProductID, product.Quantity, product.Subtotal, "now()", "now()")
	}

	err = t.db.WithTx(func(tx *database.Tx) error {
		_, err = t.db.Exec(query, args...)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionsRepository) updateProductsStock(updateProducts []entity.UpdatedProduct) error {
	var (
		err error
	)

	if len(updateProducts) > 0 {
		for _, product := range updateProducts {
			err = t.updateStockByProductID(&product)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New(constants.ErrProductNotFound)
	}

	return nil
}

func (t *TransactionsRepository) updateStockByProductID(product *entity.UpdatedProduct) error {
	var (
		query string
		err   error
	)

	query = "UPDATE products SET stock = $1, updated_at = $2 WHERE id = $3"

	err = t.db.WithTx(func(tx *database.Tx) error {
		return tx.WithStmt(query, func(stmt *database.Stmt) error {
			_, err := stmt.Exec(product.Stock, "now()", product.ID)
			return err
		})
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionsRepository) getTransactionsDetailByTransactionID(transactionId int) []entity.CheckoutProduct {
	var (
		checkoutProducts []entity.CheckoutProduct
		checkoutProduct  entity.CheckoutProduct
		query            string
		err              error
	)

	checkoutProducts = make([]entity.CheckoutProduct, 0)

	query = "SELECT products.id, products.name, categories.id as category_id, categories.name as category_name, transaction_details.id, transaction_details.transaction_id, transaction_details.quantity, transaction_details.subtotal FROM transaction_details JOIN products ON transaction_details.product_id = products.id JOIN categories ON products.category_id = categories.id WHERE transaction_details.transaction_id = $1"

	err = t.db.WithStmt(query, func(stmt *database.Stmt) error {
		scanFn := func(rows *database.Rows) error {
			if err := rows.Scan(&checkoutProduct.ProductID, &checkoutProduct.Name, &checkoutProduct.CategoryID, &checkoutProduct.CategoryName, &checkoutProduct.TransactionDetailID, &checkoutProduct.TransactionID, &checkoutProduct.Quantity, &checkoutProduct.Subtotal); err != nil {
				return err
			}
			checkoutProducts = append(checkoutProducts, checkoutProduct)
			return nil
		}

		return stmt.Query(scanFn, transactionId)
	})

	if err != nil {
		return nil
	}

	return checkoutProducts
}
