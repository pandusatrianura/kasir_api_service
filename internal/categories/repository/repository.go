package repository

import (
	"errors"

	constants "github.com/pandusatrianura/kasir_api_service/constant"
	"github.com/pandusatrianura/kasir_api_service/internal/categories/entity"
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
	"github.com/pandusatrianura/kasir_api_service/pkg/datetime"
)

type ICategoryRepository interface {
	CreateCategory(category *entity.Category) error
	UpdateCategory(id int64, category *entity.Category) error
	DeleteCategory(id int64) error
	GetCategoryByID(id int64) (*entity.ResponseCategory, error)
	GetAllCategories() ([]entity.ResponseCategory, error)
}

type CategoryRepository struct {
	db *database.DB
}

func NewCategoryRepository(db *database.DB) CategoryRepository {
	return CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(category *entity.Category) error {
	var (
		query string
		err   error
	)

	query = "INSERT INTO categories (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4)"

	err = r.db.WithTx(func(tx *database.Tx) error {
		err = tx.WithStmt(query, func(stmt *database.Stmt) error {
			_, err = stmt.Exec(category.Name, category.Description, "now()", "now()")
			if err != nil {
				return err
			}

			return nil
		})

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

func (r *CategoryRepository) UpdateCategory(id int64, category *entity.Category) error {
	var (
		query string
		err   error
	)

	query = "UPDATE categories SET name = $1, description = $2, updated_at = $3 WHERE id = $4"

	err = r.db.WithTx(func(tx *database.Tx) error {
		err = tx.WithStmt(query, func(stmt *database.Stmt) error {
			_, err = stmt.Exec(category.Name, category.Description, "now()", id)
			if err != nil {
				return err
			}

			return nil
		})

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

func (r *CategoryRepository) DeleteCategory(id int64) error {
	var (
		query string
		err   error
	)

	query = "DELETE FROM categories WHERE id = $1"

	err = r.db.WithTx(func(tx *database.Tx) error {
		err = tx.WithStmt(query, func(stmt *database.Stmt) error {
			_, err = stmt.Exec(id)
			if err != nil {
				return err
			}
			return nil
		})

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

func (r *CategoryRepository) GetCategoryByID(id int64) (*entity.ResponseCategory, error) {
	var (
		category     entity.Category
		respCategory entity.ResponseCategory
		err          error
		query        string
	)

	query = "SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1"

	err = r.db.WithStmt(query, func(stmt *database.Stmt) error {
		err = stmt.Query(func(rows *database.Rows) error {
			if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
				return err
			}

			return nil
		}, id)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if category.ID == 0 {
		return nil, errors.New(constants.ErrCategoryNotFound)
	}

	createdAt, _ := datetime.ParseTime(category.CreatedAt)
	updatedAt, _ := datetime.ParseTime(category.UpdatedAt)

	respCategory = entity.ResponseCategory{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return &respCategory, nil
}

func (r *CategoryRepository) GetAllCategories() ([]entity.ResponseCategory, error) {
	var (
		categories []entity.Category
		err        error
		query      string
	)

	query = "SELECT id, name, description, created_at, updated_at FROM categories"

	err = r.db.WithStmt(query, func(stmt *database.Stmt) error {
		err = stmt.Query(func(rows *database.Rows) error {
			var category entity.Category
			if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
				return err
			}

			categories = append(categories, category)
			return nil
		})

		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	var respCategories []entity.ResponseCategory
	for _, category := range categories {
		createdAt, _ := datetime.ParseTime(category.CreatedAt)
		updatedAt, _ := datetime.ParseTime(category.UpdatedAt)

		respCategory := entity.ResponseCategory{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		respCategories = append(respCategories, respCategory)
	}

	return respCategories, nil
}
