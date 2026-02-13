package service

import (
	"github.com/pandusatrianura/kasir_api_service/internal/users/entity"
	"github.com/pandusatrianura/kasir_api_service/internal/users/repository"
)

type IUserService interface {
	GetUserByEmail(email string) (*entity.User, error)
}

type UserService struct {
	userRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &UserService{userRepo: userRepo}
}

func (u UserService) GetUserByEmail(email string) (*entity.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
