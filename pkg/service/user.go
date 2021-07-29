package service

import (
	"nnw_s/pkg/model"
	"nnw_s/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetUserById(userId string) (*model.User, error) {
	user, err := us.repo.GetUserByIdDb(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) CreateUser(user model.User) (*string, error) {
	id, err := us.repo.CreateUserDb(user)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (us *UserService) UpdateUser(user model.User) error {
	err := us.repo.UpdateUserDb(user)
	if err != nil {
		return err
	}

	return nil
}
