package repository

import (
	"github.com/pandusatrianura/kasir_api_service/internal/users/entity"
	"github.com/pandusatrianura/kasir_api_service/pkg/database"
)

type IUserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
}

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var (
		user  entity.User
		err   error
		roles []entity.Role
		role  entity.Role
	)

	queryUser := "SELECT id, name, email, password, photo, phone FROM users WHERE email = $1"
	err = u.db.WithStmt(queryUser, func(stmt *database.Stmt) error {
		scanFn := func(rows *database.Rows) error {
			return rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Photo, &user.Phone)
		}

		return stmt.Query(scanFn, email)
	})

	if err != nil {
		return nil, err
	}

	queryRole := "SELECT roles.id, roles.name FROM roles join user_role on user_role.role_id = roles.id WHERE user_role.user_id = $1"
	scanFn := func(stmt *database.Stmt) error {
		err = stmt.QueryRow(user.ID).Scan(&role.ID, &role.Name)
		if err != nil {
			return err
		}

		roles = append(roles, role)
		return nil
	}

	err = u.db.WithStmt(queryRole, scanFn)
	if err != nil {
		return nil, err
	}

	user.Roles = roles
	return &user, nil
}
